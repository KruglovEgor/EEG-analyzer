import React from 'react';
import type { EEGLinePlot, EEGSeriesMetadata, EEGCombinedDataPoint } from "../../types/vizTypes.ts";
// Import all necessary Recharts components
import {
    LineChart,
    AreaChart,
    Line,
    Area,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
    ResponsiveContainer
} from "recharts";
import {Box, useTheme, colors, type Theme, type Palette, Typography} from "@mui/material";
// Import Color object type from MUI to use in the type guard
import type { PaletteColor } from '@mui/material/styles';

type LinePlotVisualizerProps = EEGLinePlot & {
    width?: number | `${number}%` | undefined,
    height?: number | `${number}%` | undefined
};

// 🎨 NEW: Custom Tooltip Content Component
const CustomTooltip = (
    { active, payload, label, theme, xAxisName, yAxisName }: any
) => {
    if (active && payload && payload.length) {
        // Use MUI's dark paper background for the tooltip
        return (
            <Box
                sx={{
                    p: 1,
                    bgcolor: theme.palette.grey[900], // Dark background
                    color: '#fff', // White text
                    border: '1px solid ' + theme.palette.grey[700],
                    borderRadius: 1,
                    opacity: 0.95,
                    pointerEvents: 'none' // Ensures it doesn't block chart interaction
                }}
            >
                {/* Display X-Axis Label (Time or Frequency) */}
                <Typography variant="caption" sx={{ color: theme.palette.grey[400] }}>
                    {xAxisName ?? 'X-Value'}: {label.toFixed(4)}
                </Typography>
                {/* Display each curve's data point */}
                {payload.map((item: any, index: number) => {
                    const finalName = item.name === item.dataKey && yAxisName ? yAxisName : item.name;
                    const value = (item.value as number);
                    const formattedValue = Math.abs(value) >= 1e-4 ?
                        (value.toFixed(4)) :
                        (Math.abs(value) > 1e-18 ? value.toString() : 0);

                    return (
                        <Typography
                            key={index}
                            variant="body2"
                            sx={{color: item.stroke || item.fill}} // Color dot/text matches curve color
                        >
                            {finalName}: {formattedValue}
                        </Typography>
                    );
                })}
            </Box>
        );
    }
    return null;
};

// 1. Define the type of object that has a 'main' property (the 'Color' type)
// This is a type predicate function to safely check if the palette member is a Color object.
const isColorObject = (paletteMember: Palette[keyof Palette]): paletteMember is PaletteColor => {
    // A Color object has 'light', 'main', and 'dark' properties defined.
    return typeof paletteMember === 'object' && paletteMember !== null && 'main' in paletteMember;
};


// 2. The main resolution function
const resolveColor = (theme: Theme, colorKey: string | undefined): string => {
    if (!colorKey) {
        // Fallback to a safe grey if no color is provided
        return theme.palette.grey[500];
    }

    // --- A. Check for Theme Palette Colors (primary, secondary, etc.) ---
    const paletteKey = colorKey as keyof typeof theme.palette;
    const paletteMember = theme.palette[paletteKey];

    if (isColorObject(paletteMember)) {
        // TypeScript now knows this member has a '.main' property
        return paletteMember.main;
    }

    // --- B. Check for Utility Colors (red, blue, green, etc.) ---
    const muiColorsKey = colorKey as keyof typeof colors;

    // We must check if the color exists AND if it has the '500' shade,
    // using a type assertion to satisfy the index access requirement.
    if (colors[muiColorsKey] && (colors[muiColorsKey] as any)['500']) {
        return (colors[muiColorsKey] as any)['500'];
    }

    // --- C. Default fallback (Hex code or invalid string) ---
    return colorKey;
};

const LinePlotVisualizer = ({
                                xMin, xMax, yMin, yMax,
                                xAxisName, yAxisName,
                                yLogarithmic = false,
                                showLegend = false,
                                area = false,

                                // NEW Props from the refactored EEGLinePlot type
                                seriesMetadata,
                                data,

                                width = "100%",
                                height = 256
                            }: LinePlotVisualizerProps) => {

    const theme = useTheme();

    // Check if both metadata and data are present
    if (!seriesMetadata || seriesMetadata.length === 0 || !data || data.length === 0) {
        return <Box sx={{ p: 4, textAlign: 'center', color: 'text.secondary' }}>No data to visualize.</Box>;
    }

    // data is already the combinedData
    const combinedData = data;
    const ChartComponent = area ? AreaChart : LineChart;
    const lineThickness: number = 2;

    // Map the seriesMetadata to create the Recharts components
    const curveElements = seriesMetadata.map((series: EEGSeriesMetadata) => {
        const resolvedColor = resolveColor(theme, series.preferredColor);

        const Component = area ? Area : Line;

        return (
            <Component
                key={series.dataKey}
                dataKey={series.dataKey} // Direct dataKey from metadata
                name={series.legend || series.dataKey}
                stroke={resolvedColor}
                fill={area ? (resolvedColor) : undefined}
                dot={false}
                isAnimationActive={false}
                strokeWidth={lineThickness}
            />
        );
    });

    const yScale = yLogarithmic ? 'log' : 'linear';

    const yDomain: [number | 'auto', number | 'auto'] = [yMin ?? 'auto', yMax ?? 'auto'];
    const xDomain: [number | 'auto', number | 'auto'] = [xMin ?? 'auto', xMax ?? 'auto'];

    return (
        <ResponsiveContainer width={width} height={height}>
            <ChartComponent
                data={combinedData} // Use the combined data directly
                margin={{ top: 10, right: 30, left: 20, bottom: 20 }}
            >
                <CartesianGrid strokeDasharray="3 3" opacity={0.5} />

                <XAxis
                    dataKey="x" // X-Axis key is always 'x'
                    type="number"
                    domain={xDomain}
                    label={xAxisName ? { value: xAxisName, position: 'bottom', offset: 0, style: { fontSize: '14px' } } : undefined}
                />

                <YAxis
                    type="number"
                    scale={yScale}
                    domain={yDomain}
                    label={yAxisName ? { value: yAxisName, angle: -90, position: 'left', style: { textAnchor: 'middle', fontSize: '14px' } } : undefined}
                />

                <Tooltip
                    content={<CustomTooltip theme={theme} xAxisName={xAxisName} yAxisName={yAxisName}/>}
                />

                {showLegend && <Legend />}

                {curveElements}

            </ChartComponent>
        </ResponsiveContainer>
    );
}

export default LinePlotVisualizer;