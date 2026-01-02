import type {BrainZone, RhythmType} from "./eegTypes.ts";

const AnalysisMode = {
    GROUP: 'GROUP',
    SINGLE: 'SINGLE'
} as const;
type AnalysisMode = (typeof AnalysisMode)[keyof typeof AnalysisMode];

interface EEGFileConfig {
    id: string;
    filename: string;
    experimentName: string;
    timeColumn: string;
    amplitudeColumn: string;
    rawFile: File | null;
    serverId: string | null;
}

interface EEGSingleAnalysisFormData {
    analysisMode: typeof AnalysisMode.SINGLE;
    file: EEGFileConfig;
    brainZone: BrainZone;
    rhythms: RhythmType[];
}
interface EEGGroupAnalysisFormData {
    analysisMode: typeof AnalysisMode.GROUP;
    files: EEGFileConfig[];
    brainZone: BrainZone;
    rhythm: RhythmType;
}
type EEGAnalysisFormData = EEGSingleAnalysisFormData | EEGGroupAnalysisFormData;

export type { EEGFileConfig, EEGAnalysisFormData, AnalysisMode };