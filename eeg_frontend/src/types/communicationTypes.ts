import type {BrainZone, RhythmType} from "./eegTypes.ts";
import type {AnalysisMode, EEGFileConfig} from "./configTypes.ts";
import type {EEGPlotPair} from "./vizTypes.ts";

interface EEGBaseAnalysisRequest {
    analysisId: string; // e.g. generated UUID
    brainZone?: BrainZone;
}
interface EEGSingleAnalysisRequest extends EEGBaseAnalysisRequest {
    analysisMode: typeof AnalysisMode.SINGLE; // or AnalysisMode values
    file: EEGFileConfig;
    rhythms: RhythmType[];
}
interface EEGGroupAnalysisRequest extends EEGBaseAnalysisRequest {
    analysisMode: typeof AnalysisMode.GROUP;
    files: EEGFileConfig[];
    rhythm: RhythmType;
}
type EEGAnalysisRequest = EEGSingleAnalysisRequest | EEGGroupAnalysisRequest;

interface EEGBaseAnalysisResponse {
    analysisId: string;
}
interface EEGGroupAnalysisResponse extends EEGBaseAnalysisResponse {
    analysisMode: typeof AnalysisMode.GROUP;
    experimentNames: string[],
    rhythm: RhythmType;
    absolutePowers: [string, number][];
    relativePowers: [string, number][];
    dataByExperiment: Record<string, EEGPlotPair>;
}
interface EEGSingleAnalysisResponse extends EEGBaseAnalysisResponse {
    analysisMode: typeof AnalysisMode.SINGLE;
    experimentName: string,
    rhythms: RhythmType[];
    absolutePowers: [RhythmType, number][];
    relativePowers: [RhythmType, number][];
    dataByRhythm: Partial<Record<RhythmType, EEGPlotPair>>;
}
type EEGAnalysisResponse = EEGGroupAnalysisResponse | EEGSingleAnalysisResponse;

export type {
    EEGAnalysisRequest,
    EEGAnalysisResponse
};