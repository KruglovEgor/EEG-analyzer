import React from "react";
import { Box, Typography, Divider } from "@mui/material";
import ResultsBlock from "./ResultsBlock.tsx";
import BarPlotVisualizer from "./BarPlotVisualizer.tsx";
import type { EEGBarPlot, PreferredColor } from "../../types/vizTypes.ts";

interface BarPlotItem {
    data: EEGBarPlot;
    title: string;
    preferredColor?: string;
    subtext?: React.ReactNode;
}

interface BarPlotPairBlockProps {
    headerText: string;
    headerIcon: React.ReactNode;
    leftPlot: BarPlotItem;
    rightPlot: BarPlotItem;
    height?: number;
}

const BarPlotPairBlock = ({
                              headerText,
                              headerIcon,
                              leftPlot,
                              rightPlot,
                              height = 256
                          }: BarPlotPairBlockProps) => {

    const renderPlotItem = (item: BarPlotItem) => (
        <Box sx={{
            flexBasis: { xs: '100%', md: 'calc(50% - 12px)' },
            minWidth: 0,
            display: 'flex',
            flexDirection: 'column'
        }}>
            {/* 1. Individual Plot Title */}
            <Typography variant="subtitle2" color="text.secondary" sx={{ mb: 1, fontWeight: 600 }}>
                {item.title}
            </Typography>

            {/* 2. The Visualizer */}
            <BarPlotVisualizer
                {...item.data}
                preferredColor={item.preferredColor ?? item.data.preferredColor}
                height={height}
            />

            {/* 3. Subtext/Children Slot */}
            {item.subtext && (
                <Box sx={{ mt: 1.5 }}>
                    <Divider sx={{ mb: 1, opacity: 0.5 }} />
                    <Box sx={{ typography: 'body2', color: 'text.secondary' }}>
                        {item.subtext}
                    </Box>
                </Box>
            )}
        </Box>
    );

    return (
        <ResultsBlock
            headerText={headerText}
            headerIcon={headerIcon}
        >
            <Box sx={{
                display: 'flex',
                gap: 3,
                flexWrap: 'wrap',
                mt: 1
            }}>
                {renderPlotItem(leftPlot)}
                {renderPlotItem(rightPlot)}
            </Box>
        </ResultsBlock>
    );
}

export default BarPlotPairBlock;