import type {EEGAnalysisResponse} from "./communicationTypes.ts";

export type AppState =
    | { status: 'IDLE' }
    | { status: 'LOADING' }
    | { status: 'SUCCESS'; response: EEGAnalysisResponse }
    | { status: 'ERROR'; message: string };