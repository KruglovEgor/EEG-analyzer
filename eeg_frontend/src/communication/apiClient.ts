// src/communication/apiClient.ts
//
// Thin client for sending EEGAnalysisRequest payloads to a backend and receiving
// EEGAnalysisResponse payloads.
//
// Goals:
// - Keep runtime-safe parsing at the boundary (JSON -> TS types).
// - Support partial dictionaries for rhythms (dataByRhythm) and experiment maps.
// - Tolerate common backend encodings:
//     * dataByRhythm / dataByExperiment as object maps OR [key, value][] entries.
//     * absolutePowers / relativePowers as [key, number][] OR Record<string, number>.
//
// IMPORTANT: AnalysisMode is a runtime constant (not just a type), so we must
// import it as a value.

// src/communication/apiClient.ts
//
// Thin client for sending EEGAnalysisRequest payloads to a backend and receiving
// EEGAnalysisResponse payloads.
//
// Goals:
// - Keep runtime-safe parsing at the boundary (JSON -> TS types).
// - Support partial dictionaries for rhythms (dataByRhythm) and experiment maps.
// - Tolerate common backend encodings:
//     * dataByRhythm / dataByExperiment as object maps OR [key, value][] entries.
//     * absolutePowers / relativePowers as [key, number][] OR Record<string, number>.
import type { EEGAnalysisRequest, EEGAnalysisResponse } from "../types/communicationTypes.ts";
import type { RhythmType } from "../types/eegTypes.ts";
import { isRhythmType } from "../types/eegTypes.ts";
import type { EEGPlotPair } from "../types/vizTypes.ts";
import type {EEGAnalysisFormData} from "../types/configTypes.ts";
import {generateUUID} from "../util/uuidUtils.ts";

export interface ApiClientOptions {
    /** Base URL like "http://localhost:8000" */
    baseUrl: string;
    /** Endpoint path like "/analysis" (default) */
    analysisEndpoint?: string;
    /** Inject custom fetch (tests) */
    fetchFn?: typeof fetch;
    /** Abort support */
    signal?: AbortSignal;
    /** If true, unknown rhythm keys from JSON are dropped instead of cast. */
    dropUnknownRhythms?: boolean;
}

export class ApiError extends Error {
    public readonly status: number;
    public readonly statusText: string;
    public readonly bodyText: string;

    constructor(status: number, statusText: string, bodyText: string) {
        super(`Request failed: ${status} ${statusText}${bodyText ? ` - ${bodyText}` : ""}`);
        this.status = status;
        this.statusText = statusText;
        this.bodyText = bodyText;
    }
}

export function eegFormDataToRequest(formData: EEGAnalysisFormData): EEGAnalysisRequest {
    const analysisId = generateUUID();

    if (formData.analysisMode === 'SINGLE') {
        return {
            analysisId: analysisId,
            analysisMode: 'SINGLE',
            brainZone: formData.brainZone,
            file: formData.file,
            rhythms: formData.rhythms
        };
    }
    if (formData.analysisMode === 'GROUP') {
        return {
            analysisId: analysisId,
            analysisMode: 'GROUP',
            brainZone: formData.brainZone,
            files: formData.files,
            rhythm: formData.rhythm
        };
    }
    throw new Error("Cannot create a request out of current form data: unknown analysis mode");
}

function isPlainObject(value: unknown): value is Record<string, unknown> {
    return typeof value === "object" && value !== null && !Array.isArray(value);
}

function requireString(value: unknown, field: string): string {
    if (typeof value !== "string") {
        throw new Error(`Expected ${field} to be a string`);
    }
    return value;
}

function requireNumber(value: unknown, field: string): number {
    if (typeof value !== "number" || Number.isNaN(value)) {
        throw new Error(`Expected ${field} to be a number`);
    }
    return value;
}

function requireArray(value: unknown, field: string): unknown[] {
    if (!Array.isArray(value)) {
        throw new Error(`Expected ${field} to be an array`);
    }
    return value;
}

function parsePlotPair(value: unknown, field: string): EEGPlotPair {
    // The plot structure is nested; the app should be robust even if the backend
    // slightly changes metadata fields. Require only that it's an object.
    if (!isPlainObject(value)) {
        throw new Error(`Expected ${field} to be an object`);
    }
    return value as unknown as EEGPlotPair;
}

function normalizePairs(raw: unknown, field: string): Array<[string, number]> {
    // Accept:
    // 1) [ [key, number], ... ]
    // 2) { key: number, ... }
    if (Array.isArray(raw)) {
        return raw.map((entry, i) => {
            if (!Array.isArray(entry) || entry.length !== 2) {
                throw new Error(`Expected ${field}[${i}] to be [string, number]`);
            }
            const k = requireString(entry[0], `${field}[${i}][0]`);
            const v = requireNumber(entry[1], `${field}[${i}][1]`);
            return [k, v];
        });
    }

    if (isPlainObject(raw)) {
        const out: Array<[string, number]> = [];
        for (const [k, v] of Object.entries(raw)) {
            out.push([k, requireNumber(v, `${field}.${k}`)]);
        }
        return out;
    }

    throw new Error(`Expected ${field} to be an entries array or object map`);
}

function normalizeMap<V>(
    raw: unknown,
    field: string,
    valueParser: (value: unknown, key: string) => V
): Record<string, V> {
    // Accept:
    // 1) { key: value }
    // 2) [ [key, value], ... ]
    if (isPlainObject(raw)) {
        const out: Record<string, V> = {};
        for (const [k, v] of Object.entries(raw)) {
            out[k] = valueParser(v, k);
        }
        return out;
    }

    if (Array.isArray(raw)) {
        const out: Record<string, V> = {};
        for (let i = 0; i < raw.length; i++) {
            const entry = raw[i];
            if (!Array.isArray(entry) || entry.length !== 2) {
                throw new Error(`Expected ${field}[${i}] to be [key, value]`);
            }
            const k = requireString(entry[0], `${field}[${i}][0]`);
            out[k] = valueParser(entry[1], k);
        }
        return out;
    }

    throw new Error(`Expected ${field} to be an object map or entries array`);
}

function parseSingleResponse(
    raw: Record<string, unknown>,
    dropUnknownRhythms: boolean
): EEGAnalysisResponse {
    const analysisId = requireString(raw.analysisId, "analysisId");
    const experimentName = requireString(raw.experimentName, "experimentName");

    const rhythmsRaw = requireArray(raw.rhythms, "rhythms");
    const rhythms = rhythmsRaw.map((v, i) => requireString(v, `rhythms[${i}]`)) as RhythmType[];

    const absolutePairs = normalizePairs(raw.absolutePowers, "absolutePowers");
    const relativePairs = normalizePairs(raw.relativePowers, "relativePowers");

    // Convert keys to RhythmType where possible.
    const absolutePowers: Array<[RhythmType, number]> = [];
    for (const [k, v] of absolutePairs) {
        if (!dropUnknownRhythms || isRhythmType(k)) {
            absolutePowers.push([k as RhythmType, v]);
        }
    }

    const relativePowers: Array<[RhythmType, number]> = [];
    for (const [k, v] of relativePairs) {
        if (!dropUnknownRhythms || isRhythmType(k)) {
            relativePowers.push([k as RhythmType, v]);
        }
    }

    const dataRaw = raw.dataByRhythm ?? {};
    const dataObj = normalizeMap<EEGPlotPair>(dataRaw, "dataByRhythm", (v, key) =>
        parsePlotPair(v, `dataByRhythm.${key}`)
    );

    const dataByRhythm: Partial<Record<RhythmType, EEGPlotPair>> = {};
    for (const [k, v] of Object.entries(dataObj)) {
        if (dropUnknownRhythms) {
            if (isRhythmType(k)) {
                dataByRhythm[k] = v;
            }
        } else {
            dataByRhythm[k as RhythmType] = v;
        }
    }

    return {
        analysisId,
        analysisMode: 'SINGLE',
        experimentName,
        rhythms,
        absolutePowers,
        relativePowers,
        dataByRhythm,
    } as const;
}

function parseGroupResponse(raw: Record<string, unknown>): EEGAnalysisResponse {
    const analysisId = requireString(raw.analysisId, "analysisId");

    const experimentNamesRaw = requireArray(raw.experimentNames, "experimentNames");
    const experimentNames = experimentNamesRaw.map((v, i) =>
        requireString(v, `experimentNames[${i}]`)
    );

    const rhythm = requireString(raw.rhythm, "rhythm") as RhythmType;

    const absolutePowers = normalizePairs(raw.absolutePowers, "absolutePowers");
    const relativePowers = normalizePairs(raw.relativePowers, "relativePowers");

    const dataByExperiment = normalizeMap<EEGPlotPair>(
        raw.dataByExperiment,
        "dataByExperiment",
        (v, key) => parsePlotPair(v, `dataByExperiment.${key}`)
    );

    return {
        analysisId,
        analysisMode: 'GROUP',
        experimentNames,
        rhythm,
        absolutePowers,
        relativePowers,
        dataByExperiment,
    } as const;
}

function parseAnalysisResponse(rawJson: unknown, dropUnknownRhythms: boolean): EEGAnalysisResponse {
    if (!isPlainObject(rawJson)) {
        throw new Error("Response JSON must be an object");
    }

    const mode = requireString(rawJson.analysisMode, "analysisMode");

    if (mode === 'SINGLE') {
        return parseSingleResponse(rawJson, dropUnknownRhythms);
    }

    if (mode === 'GROUP') {
        return parseGroupResponse(rawJson);
    }

    throw new Error(`Unknown analysisMode: ${mode}`);
}

export function createApiClient(options: ApiClientOptions) {
    const fetchFn = options.fetchFn ?? fetch;
    const endpoint = options.analysisEndpoint ?? "/analysis";
    const baseUrl = options.baseUrl.replace(/\/+$/, "");
    const dropUnknownRhythms = options.dropUnknownRhythms ?? true;

    async function analyze(
        request: EEGAnalysisRequest,
        signal?: AbortSignal
    ): Promise<EEGAnalysisResponse> {
        const url = `${baseUrl}${endpoint}`;

        const res = await fetchFn(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Accept": "application/json",
            },
            body: JSON.stringify(request),
            signal: signal ?? options.signal,
        });

        if (!res.ok) {
            const text = await res.text().catch(() => "");
            throw new ApiError(res.status, res.statusText, text);
        }

        const json = (await res.json()) as unknown;
        return parseAnalysisResponse(json, dropUnknownRhythms);
    }

    return {
        analyze,
    } as const;
}

// Convenience: one-shot call without constructing the client.
export async function analyzeEEG(
    baseUrl: string,
    request: EEGAnalysisRequest,
    opts?: Omit<ApiClientOptions, "baseUrl">
): Promise<EEGAnalysisResponse> {
    const client = createApiClient({ baseUrl, ...opts });
    return client.analyze(request);
}