import React, { useMemo } from 'react';
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    Cell
} from "recharts";
import { Box, useTheme, colors, type Theme, type Palette, Typography } from "@mui/material";
import type { PaletteColor } from '@mui/material/styles';
import type { EEGBarPlot } from "../../types/vizTypes.ts";

type BarPlotVisualizerProps = EEGBarPlot & {
    width?: number | `${number}%` | undefined,
    height?: number | `${number}%` | undefined
};

// 🎨 Custom Tooltip for Categorical Data
const CustomTooltip = ({ active, payload, label, theme, xAxisName, yAxisName }: any) => {
    if (active && payload && payload.length) {
        return (
            <Box
                sx={{
                    p: 1,
                    bgcolor: theme.palette.grey[900],
                    color: '#fff',
                    border: '1px solid ' + theme.palette.grey[700],
                    borderRadius: 1,
                    opacity: 0.95,
                    pointerEvents: 'none'
                }}
            >
                <Typography variant="caption" sx={{ color: theme.palette.grey[400], display: 'block', mb: 0.5 }}>
                    {xAxisName ?? 'Category'}: {label}
                </Typography>
                {payload.map((item: any, index: number) => (
                    <Typography
                        key={index}
                        variant="body2"
                        sx={{ color: item.fill, fontWeight: 'bold' }}
                    >
                        {yAxisName ?? 'Value'}: {item.value.toFixed(4)}
                    </Typography>
                ))}
            </Box>
        );
    }
    return null;
};

// --- Color Resolution Helpers ---
const isColorObject = (paletteMember: Palette[keyof Palette]): paletteMember is PaletteColor => {
    return typeof paletteMember === 'object' && paletteMember !== null && 'main' in paletteMember;
};

const resolveColor = (theme: Theme, colorKey: string | undefined): string => {
    if (!colorKey) return theme.palette.primary.main;

    const paletteKey = colorKey as keyof typeof theme.palette;
    const paletteMember = theme.palette[paletteKey];

    if (isColorObject(paletteMember)) return paletteMember.main;

    const muiColorsKey = colorKey as keyof typeof colors;
    if (colors[muiColorsKey] && (colors[muiColorsKey] as any)['500']) {
        return (colors[muiColorsKey] as any)['500'];
    }

    return colorKey;
};

const BarPlotVisualizer = ({
                               labels,
                               values,
                               preferredColor,
                               xAxisName,
                               yAxisName,
                               yMin,
                               yMax,
                               width = "100%",
                               height = 256
                           }: BarPlotVisualizerProps) => {
    const theme = useTheme();

    // 1. Transform separate arrays into Recharts-friendly objects
    const data = useMemo(() => {
        return labels.map((label, index) => ({
            name: label,
            value: values[index] ?? 0
        }));
    }, [labels, values]);

    if (!labels.length || !values.length) {
        return <Box sx={{ p: 4, textAlign: 'center', color: 'text.secondary' }}>No data to visualize.</Box>;
    }

    const resolvedColor = resolveColor(theme, preferredColor);
    const yDomain: [number | 'auto', number | 'auto'] = [yMin ?? 'auto', yMax ?? 'auto'];

    return (
        <ResponsiveContainer width={width} height={height}>
            <BarChart
                data={data}
                margin={{ top: 10, right: 30, left: 20, bottom: 20 }}
            >
                <CartesianGrid strokeDasharray="3 3" opacity={0.3} vertical={false} />

                <XAxis
                    dataKey="name"
                    tick={{ fontSize: 12 }}
                    label={xAxisName ? { value: xAxisName, position: 'bottom', offset: 0, style: { fontSize: '14px' } } : undefined}
                />

                <YAxis
                    domain={yDomain}
                    tick={{ fontSize: 12 }}
                    label={yAxisName ? { value: yAxisName, angle: -90, position: 'left', style: { textAnchor: 'middle', fontSize: '14px' } } : undefined}
                />

                <Tooltip
                    content={<CustomTooltip theme={theme} xAxisName={xAxisName} yAxisName={yAxisName} />}
                />

                <Bar
                    dataKey="value"
                    fill={resolvedColor}
                    isAnimationActive={false}
                    radius={[4, 4, 0, 0]} // Rounded top corners
                />
            </BarChart>
        </ResponsiveContainer>
    );
};

export default BarPlotVisualizer;