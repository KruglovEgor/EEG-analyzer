// --- RHYTHM TYPES ---
const RhythmType = {
    ALPHA: 'ALPHA', // 8–13 Hz (Idling)
    BETA: 'BETA',   // 13–30 Hz (Active, concentration)
    GAMMA: 'GAMMA', // 30–100 Hz (Higher cognitive processing)
    DELTA: 'DELTA', // 0.5–4 Hz (Deep sleep, infancy)
    THETA: 'THETA', // 4–8 Hz (Drowsiness, meditation, memory retrieval)
    KAPPA: 'KAPPA', // 8–13 Hz (A variant of alpha, often over temporal)
    LAMBDA: 'LAMBDA', // 4–8 Hz (Associated with visual exploration)
    MU: 'MU'        // 8–13 Hz (Associated with motor inhibition)
} as const;
const ALL_RHYTHM_TYPES = Object.values(RhythmType) as RhythmType[];
const isRhythmType = (value: any): value is RhythmType => {
    return ALL_RHYTHM_TYPES.includes(value);
};

// --- BRAIN ZONES ---
const BrainZone = {
    FRONTAL: 'FRONTAL',     // Cognitive control, executive function
    PARIETAL: 'PARIETAL',   // Sensory integration, spatial awareness
    CENTRAL: 'CENTRAL',     // Motor and somatosensory cortex
    TEMPORAL: 'TEMPORAL',   // Auditory processing, language, memory
    OCCIPITAL: 'OCCIPITAL'  // Visual processing
} as const;
const ALL_BRAIN_ZONES = Object.values(BrainZone) as BrainZone[];

const DEFAULT_RHYTHM_BANDS = {
    [RhythmType.DELTA]: [0.5, 4],
    [RhythmType.THETA]: [4, 8],
    [RhythmType.ALPHA]: [8, 13],
    [RhythmType.BETA]: [13, 30],
    [RhythmType.GAMMA]: [30, 100],
    [RhythmType.MU]: [8, 13],
    [RhythmType.LAMBDA]: [4, 8],
    [RhythmType.KAPPA]: [8, 13],
} as const;

// 🚨 FIX 1: Update the type to include 'readonly' tuple!
type RhythmBand = readonly [number, number];

const DEFAULT_RHYTHM_BANDS_CONFIG: Record<RhythmType, RhythmBand> = DEFAULT_RHYTHM_BANDS;

// --- 2. RHYTHM TYPES BY BRAIN ZONE (Common Groupings) ---
// This map suggests which rhythms are typically most dominant or studied 
// within each major brain zone.
// Note: This is simplified, as all rhythms exist everywhere, but analysis often
// focuses on these key regional associations.
const RHYTHM_TYPES_BY_BRAIN_ZONE: Record<BrainZone, RhythmType[]> = {
    [BrainZone.FRONTAL]: [
        RhythmType.BETA,
        RhythmType.THETA,
        RhythmType.GAMMA,
    ],
    [BrainZone.PARIETAL]: [
        RhythmType.ALPHA,
        RhythmType.BETA,
    ],
    [BrainZone.CENTRAL]: [
        RhythmType.MU,
        RhythmType.BETA,
    ],
    [BrainZone.TEMPORAL]: [
        RhythmType.THETA,
        RhythmType.GAMMA,
    ],
    [BrainZone.OCCIPITAL]: [
        RhythmType.ALPHA,
        RhythmType.LAMBDA,
    ]
} as const;

export type RhythmType = (typeof RhythmType)[keyof typeof RhythmType];
export type BrainZone = (typeof BrainZone)[keyof typeof BrainZone];
export type { RhythmBand };
export {
    ALL_BRAIN_ZONES,
    ALL_RHYTHM_TYPES,
    isRhythmType,
    DEFAULT_RHYTHM_BANDS_CONFIG as DEFAULT_RHYTHM_BANDS,
    RHYTHM_TYPES_BY_BRAIN_ZONE
};
