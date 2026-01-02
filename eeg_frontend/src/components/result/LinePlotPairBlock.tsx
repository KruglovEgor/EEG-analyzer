import React from "react";
import {Box, Grid} from "@mui/material";
import { Activity, Zap } from "lucide-react"; // Example icons
import ResultsBlock from "./ResultsBlock.tsx"; // Assuming this path
import LinePlotVisualizer from "./LinePlotVisualizer.tsx"; // Assuming this path
import type { EEGPlotPair } from "../../types/vizTypes.ts";

interface LinePlotPairBlockProps {
    /** The title for the entire results block. */
    headerText: string;
    /** The icon for the entire results block. */
    headerIcon: React.ReactNode;
    /** The data containing the two plots (PSD and Signal). */
    plotPair: EEGPlotPair;
    height?: number
}

const LinePlotPairBlock = ({
    headerText,
    headerIcon,
    plotPair,
    height = 256
}: LinePlotPairBlockProps) => {

    const { psdPlot, signalPlot } = plotPair;

    return (
        <ResultsBlock
            headerText={headerText}
            headerIcon={headerIcon}
        >
            <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
                {/* 1. PSD Plot */}
                <Box sx={{
                    flexBasis: { xs: '100%', md: 'calc(50% - 8px)' }, // 50% width minus half the gap
                    minWidth: 0, // Allows content to shrink
                }}>
                    <LinePlotVisualizer {...psdPlot} height={height} />
                </Box>

                {/* 2. Signal/Time Domain Plot */}
                <Box sx={{
                    flexBasis: { xs: '100%', md: 'calc(50% - 8px)' },
                    minWidth: 0,
                }}>
                    <LinePlotVisualizer {...signalPlot} height={height} />
                </Box>
            </Box>
        </ResultsBlock>
    );
}

export default LinePlotPairBlock;