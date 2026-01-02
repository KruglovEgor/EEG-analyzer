import ResultsBlock from "./ResultsBlock.tsx";
import {BarChart2, BarChart2Icon, BarChartIcon, BrainCircuit, HomeIcon, LineChartIcon, Microscope} from "lucide-react";
import LinePlotVisualizer from "./LinePlotVisualizer.tsx";
import LinePlotPairBlock from "./LinePlotPairBlock.tsx";
import type {EEGAnalysisResponse} from "../../types/communicationTypes.ts";
import {t} from "i18next";
import {Button, Typography} from "@mui/material";
import BarPlotVisualizer from "./BarPlotVisualizer.tsx";
import {ALL_RHYTHM_TYPES, isRhythmType, type RhythmType} from "../../types/eegTypes.ts";
import BarPlotPairBlock from "./BarPlotPairBlock.tsx";
import type {EEGPlotPair} from "../../types/vizTypes.ts";

interface ResultsPageProps {
    analysis: Partial<EEGAnalysisResponse>,
    onGoHome?: () => void
}

const localizeRhythm = (rhythm: RhythmType) => t(`misc_rhythm_${rhythm}`);

const elemMaxMinPower = <T extends (string | RhythmType)>(
    powers: [T, number][],
    getName: (item: [T, number]) => string,
    isPercentage: boolean = false
) => {
    if (!powers || powers.length === 0) return null;

    // Find extremas
    const maxElem = powers.reduce((prev, curr) => (curr[1] > prev[1] ? curr : prev));
    const minElem = powers.reduce((prev, curr) => (curr[1] < prev[1] ? curr : prev));

    // Formatting helper for scientific values
    const format = (val: number) =>
        Math.abs(val) >= 1e-4 ? val.toFixed(4) : val.toExponential(2);

    return (
        <>
            <Typography variant="body2">
                {t('results_reportPower_maxPower', {
                    name: getName(maxElem),
                    value: isPercentage ? (format(maxElem[1] * 100) + '%') : format(maxElem[1])
                })}
            </Typography>
            <Typography variant="body2">
                {t('results_reportPower_minPower', {
                    name: getName(minElem),
                    value: isPercentage ? (format(minElem[1] * 100) + '%') : format(minElem[1])
                })}
            </Typography>
        </>
    );
};

const ResultsPage = ({ analysis, onGoHome }: ResultsPageProps) => {
    return (
        <>
            {/* Header block */}
            <ResultsBlock headerText={
                analysis.analysisMode === "SINGLE" ? t('results_reportHeader_single') :
                    analysis.analysisMode === "GROUP" ? t('results_reportHeader_group') :
                    t('results_reportHeader_unknown')
            } headerIcon={<Microscope size={18} />}>
                {analysis.analysisMode && <Typography>
                    {
                        analysis.analysisMode === "SINGLE" ?
                            t('results_reportHeader_single_rhythmList', { rhythmList:
                                analysis.rhythms
                                    ?.map(localizeRhythm)
                                    ?.join(', ')
                                    || t('misc_list_empty')
                            }) :
                            analysis.analysisMode === "GROUP" ?
                                t('results_reportHeader_group_recordCount', { recordCount:
                                    analysis.experimentNames?.length || t('misc_list_empty')
                                }) :
                                t('misc_list_empty')
                    }
                </Typography>}
                {onGoHome && <Button
                    variant="contained"
                    size="large"
                    fullWidth
                    onClick={onGoHome}
                    startIcon={<HomeIcon />}
                    sx={{mt: 2}}
                >{ t('misc_button_goHome') }</Button>}
            </ResultsBlock>

            {/* Absolute vs relative powers */}
            {analysis.analysisMode && <BarPlotPairBlock
                    headerText={
                        analysis.analysisMode === "SINGLE" ?
                            t('results_reportPower_single', { experimentName: analysis.experimentName }) :
                            analysis.analysisMode === "GROUP" ?
                                t('results_reportPower_group', { rhythm: analysis.rhythm }) :
                                t('results_reportPower_unknown')
                    }
                    headerIcon={<BarChart2Icon size={18}/>}
                    leftPlot={{
                        title: t('results_reportPower_absolutePower'),
                        data: {
                            labels: (
                                analysis.analysisMode === "SINGLE" ?
                                    analysis.absolutePowers?.map(r => localizeRhythm(r[0])) :
                                    analysis.absolutePowers?.map(r => r[0])
                            ) ?? [],
                            values: analysis.absolutePowers
                                ?.map(([_, value]) => value)
                                ?? [],
                            preferredColor: 'secondary',
                            yMin: 0,
                            xAxisName: (
                                analysis.analysisMode === "SINGLE" ?
                                    t('results_reportPower_single_xAxis') :
                                    t('results_reportPower_group_xAxis')
                            ),
                            yAxisName: t('results_reportPower_absolutePower')
                        },
                        subtext: analysis.absolutePowers && elemMaxMinPower(
                            analysis.absolutePowers,
                            (([r, _]) => isRhythmType(r) ? localizeRhythm(r as RhythmType) : r)
                        )
                    }}
                    rightPlot={{
                        title: t('results_reportPower_relativePower'),
                        data: {
                            labels: (
                                analysis.analysisMode == "SINGLE" ?
                                    analysis.relativePowers?.map(([r, _]) => t(`misc_rhythm_${r}`)) :
                                    analysis.relativePowers?.map(([r, _]) => r)
                            ) ?? [],
                            values: analysis.relativePowers
                                ?.map(([_, value]) => value)
                                ?? [],
                            yMin: 0,
                            preferredColor: 'warning',
                            xAxisName: (
                                analysis.analysisMode == "SINGLE" ?
                                    t('results_reportPower_single_xAxis') :
                                    t('results_reportPower_group_xAxis')
                            ),
                            yAxisName: t('results_reportPower_relativePower')
                        },
                        subtext: analysis.relativePowers && elemMaxMinPower(
                            analysis.relativePowers,
                            (([r, _]) => isRhythmType(r) ? localizeRhythm(r as RhythmType) : r)
                        )
                    }}
            />}

            {analysis.analysisMode == "GROUP" && analysis.dataByExperiment &&
                (Object.entries(analysis.dataByExperiment) as Array<[string, EEGPlotPair]>).map(
                    ([experiment, plotPair], idx) =>
                        <LinePlotPairBlock
                            height={320}
                            headerText={t('results_linePlotPair_group_header', { experimentName: experiment })}
                            headerIcon={<LineChartIcon size={18}/>}
                            plotPair={{
                                psdPlot: {
                                    xAxisName: t('results_linePlotPair_axis_frequency'),
                                    yAxisName: t('results_linePlotPair_axis_psd'),
                                    yLogarithmic: true,
                                    seriesMetadata: [
                                        {
                                            dataKey: 'psd',
                                            legend: t('results_linePlotPair_legend_psd'),
                                            preferredColor: 'secondary'
                                        }
                                    ],
                                    showLegend: true,
                                    ...plotPair.psdPlot
                                },
                                signalPlot: {
                                    xAxisName: t('results_linePlotPair_axis_time'),
                                    yAxisName: t('results_linePlotPair_axis_amplitude'),
                                    seriesMetadata: [
                                        {dataKey: 'raw', legend: t('results_linePlotPair_legend_raw')},
                                        {
                                            dataKey: 'filtered',
                                            legend: t('results_linePlotPair_legend_filtered'),
                                            preferredColor: 'primary'
                                        }
                                    ],
                                    showLegend: true,
                                    ...plotPair.signalPlot
                                }
                            }}
                            key={`${idx}_${experiment}`}
                        />
                )
            }

            {analysis.analysisMode == "SINGLE" && analysis.dataByRhythm &&
                (Object.entries(analysis.dataByRhythm) as Array<[RhythmType, EEGPlotPair]>).map(
                    ([rhythm, plotPair], idx) =>
                    <LinePlotPairBlock
                        height={320}
                        headerText={t('results_linePlotPair_single_header', { rhythm: localizeRhythm(rhythm) })}
                        headerIcon={<LineChartIcon size={18}/>}
                        plotPair={{
                            psdPlot: {
                                xAxisName: t('results_linePlotPair_axis_frequency'),
                                yAxisName: t('results_linePlotPair_axis_psd'),
                                yLogarithmic: true,
                                seriesMetadata: [
                                    {
                                        dataKey: 'psd',
                                        legend: t('results_linePlotPair_legend_psd'),
                                        preferredColor: 'secondary'
                                    }
                                ],
                                showLegend: true,
                                ...plotPair.psdPlot
                            },
                            signalPlot: {
                                xAxisName: t('results_linePlotPair_axis_time'),
                                yAxisName: t('results_linePlotPair_axis_amplitude'),
                                seriesMetadata: [
                                    {dataKey: 'raw', legend: t('results_linePlotPair_legend_raw')},
                                    {
                                        dataKey: 'filtered',
                                        legend: t('results_linePlotPair_legend_filtered'),
                                        preferredColor: 'primary'
                                    }
                                ],
                                showLegend: true,
                                ...plotPair.signalPlot
                            }
                        }}
                        key={`${idx}_${rhythm}`}
                    />
                )
            }

            {/*<LinePlotPairBlock headerText={"vnm"} headerIcon={<BrainCircuit size={18} />} plotPair={{*/}
            {/*    psdPlot: {*/}
            {/*        xAxisName: 'Time',*/}
            {/*        yAxisName: 'PSD',*/}
            {/*        area: true,*/}
            {/*        yLogarithmic: true,*/}
            {/*        seriesMetadata: [*/}
            {/*                { dataKey: "psd", legend: "PSD", preferredColor: 'secondary' },*/}
            {/*        ],*/}
            {/*        data: [*/}
            {/*            { x: 0, psd: 1e-1 },*/}
            {/*            { x: 1, psd: 3e-2 },*/}
            {/*            { x: 2, psd: 6e-3 },*/}
            {/*            { x: 3, psd: 1e-3 },*/}
            {/*            { x: 4, psd: 1.5e-4 },*/}
            {/*        ]*/}
            {/*    },*/}
            {/*    signalPlot: {*/}
            {/*        xAxisName: 'Time',*/}
            {/*        yAxisName: 'Signal',*/}
            {/*        area: false,*/}
            {/*        seriesMetadata: [*/}
            {/*            { dataKey: "raw", legend: "Raw", preferredColor: 'secondary' },*/}
            {/*            { dataKey: "filtered", legend: "Filtered", preferredColor: 'warning' }*/}
            {/*        ],*/}
            {/*        data: [*/}
            {/*            { x: 0, raw: 2.5, filtered: 0 },*/}
            {/*            { x: 0.25, raw: 2.4, filtered: -0.1 },*/}
            {/*            { x: 0.5, raw: 2.6, filtered: 0.1 },*/}
            {/*            { x: 0.75, raw: 2.5, filtered: 0 },*/}
            {/*            { x: 1, raw: 2.3, filtered: -0.2 },*/}
            {/*            { x: 1.25, raw: 2.7, filtered: 0.2 },*/}
            {/*            { x: 1.5, raw: 2.0, filtered: -0.5 },*/}
            {/*            { x: 1.75, raw: 3.0, filtered: 0.5 },*/}
            {/*            { x: 2, raw: 2.5, filtered: 0 },*/}
            {/*            { x: 4.75, raw: 2, filtered: -0.5 },*/}
            {/*            { x: 5, raw: 2.5, filtered: -0.1 },*/}
            {/*        ]*/}
            {/*    }*/}
            {/*}} />*/}
        </>
    );
}

export default ResultsPage;