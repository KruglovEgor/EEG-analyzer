import React, { useState } from "react";
import type { AnalysisMode, EEGFileConfig, EEGAnalysisFormData } from "../../types/configTypes.ts";
import { useTranslation } from "react-i18next";
import { Box, Button, Chip, FormControl, InputLabel, MenuItem, Select, Stack, FormHelperText } from "@mui/material"; // Added FormHelperText

import ModeConfigurationBlock from "./ModeConfigurationBlock.tsx";
import ConfigurationBlock from "./ConfigurationBlock.tsx";
import { Activity, Play, Upload } from "lucide-react";
import FileDropzone from "./FileDropzone.tsx";
import FileItem from "./FileItem.tsx";
import { ALL_BRAIN_ZONES, ALL_RHYTHM_TYPES, RHYTHM_TYPES_BY_BRAIN_ZONE } from '../../types/eegTypes.ts';
import type { BrainZone, RhythmType } from '../../types/eegTypes.ts';

interface ConfigurationPageProps {
    isAnalyzing?: boolean;
    onSetAnalyzing?: (analyzing: boolean) => void;
    onRunAnalysis?: (data: EEGAnalysisFormData) => void;
}

const ConfigurationPage = ({
                               isAnalyzing = false,
                               onSetAnalyzing = (_) => {},
                               onRunAnalysis = (_) => {}
                           }: ConfigurationPageProps) => {
    const [mode, setMode] = useState<AnalysisMode>("GROUP");
    const { t } = useTranslation();
    const [files, setFiles] = useState<EEGFileConfig[]>([]);

    // Default to Frontal
    const [brainZone, setBrainZone] = useState<BrainZone>("FRONTAL");

    // Initialize with the defaults for Frontal
    const [selectedRhythms, setSelectedRhythms] = useState<RhythmType[]>(
        [RHYTHM_TYPES_BY_BRAIN_ZONE["FRONTAL"][0]]
    );

    const handleFileUpload = (fileList: FileList) => {
        if (fileList.length === 0) return;

        const newFiles = Array.from(fileList).map(f => ({
            id: Math.random().toString(36).substr(2, 9),
            filename: f.name,
            experimentName: f.name.replace(/\.(csv|txt)$/i, ''), // Regex matches both extensions case-insensitive
            timeColumn: 'Time',
            amplitudeColumn: 'A0',
            rawFile: f,
            serverId: null,
        }));

        if (mode === 'SINGLE') {
            // In Single mode, we only keep the NEWEST file if multiple were dropped, or just the one.
            // (Previous logic was a bit ambiguous if fileList had >1 items)
            setFiles([newFiles[0]]);
        } else {
            setFiles(prev => [...prev, ...newFiles]);
        }
    };

    const updateFile = (id: string, field: keyof EEGFileConfig, value: any) => {
        setFiles(files.map(f => f.id === id ? { ...f, [field]: value } : f));
    };

    const removeFile = (id: string) => {
        setFiles(files.filter(f => f.id !== id));
    };

    const handleBrainZoneChange = (newBrainZone: BrainZone) => {
        setBrainZone(newBrainZone);
        // Reset rhythms to the default for this zone when zone changes
        // This is good UX so users don't see 'Occipital' with 'Frontal' rhythms selected
        setSelectedRhythms(RHYTHM_TYPES_BY_BRAIN_ZONE[newBrainZone]);
    };

    const runAnalysis = () => {
        if (files.length === 0) return;

        // Validation logic
        if (mode === 'GROUP') {
            // Strict check: Group mode must have exactly 1 rhythm
            const targetRhythm = selectedRhythms[0];
            if (!targetRhythm) return;

            const formData: EEGAnalysisFormData = {
                analysisMode: 'GROUP',
                files: files,
                brainZone: brainZone,
                rhythm: targetRhythm
            };
            onSetAnalyzing(true);
            onRunAnalysis(formData);

        } else {
            // Single mode validation
            if (selectedRhythms.length === 0) return;

            const formData: EEGAnalysisFormData = {
                analysisMode: 'SINGLE',
                file: files[0],
                brainZone: brainZone,
                rhythms: selectedRhythms
            };
            onSetAnalyzing(true);
            onRunAnalysis(formData);
        }
    }

    return (
        <>
            <ModeConfigurationBlock mode={mode} onModeChange={setMode}/>

            <ConfigurationBlock
                headerText={t('config_file_upload')}
                headerIcon={<Upload size={18} />}
            >
                <FileDropzone
                    isMultiple={mode === "GROUP"}
                    onFileUpload={handleFileUpload}
                    accept={".txt,.csv"}
                />

                <Box sx={{ mt: 2, maxHeight: 300, overflowY: 'auto' }}>
                    {files.map(f => (
                        <FileItem
                            key={f.id}
                            file={f}
                            onUpdate={updateFile}
                            onRemove={removeFile}
                        />
                    ))}
                    {files.length === 0 && (
                        <Box p={2} textAlign="center" color="text.secondary">
                            {t('config_file_noFilesSelected')}
                        </Box>
                    )}
                </Box>
            </ConfigurationBlock>

            <ConfigurationBlock
                headerText={t('config_rhythm')}
                headerIcon={<Activity size={18} />}
            >
                <Stack spacing={2}>
                    <FormControl fullWidth size="small">
                        <InputLabel>{t('config_rhythm_brainZone')}</InputLabel>
                        <Select
                            value={brainZone}
                            label={t('config_rhythm_brainZone')}
                            onChange={(e) => handleBrainZoneChange(e.target.value as BrainZone)}
                        >
                            {ALL_BRAIN_ZONES.map(z =>
                                <MenuItem key={z} value={z}>{t(`misc_brainZone_${z}`)}</MenuItem>
                            )}
                        </Select>
                    </FormControl>

                    <FormControl fullWidth size="small">
                        <InputLabel id="rhythm-select-label">
                            {(mode === "GROUP") ? t('config_rhythm_targetRhythm') : t('config_rhythm_targetRhythms')}
                        </InputLabel>
                        <Select
                            labelId="rhythm-select-label"
                            multiple={(mode === "SINGLE")}
                            // Logic: In GROUP mode, Value is a String. In SINGLE mode, Value is an Array.
                            value={(mode === "GROUP") ? (selectedRhythms[0] || '') : selectedRhythms}
                            label={(mode === "GROUP") ? t('config_rhythm_targetRhythm') : t('config_rhythm_targetRhythms')}

                            onChange={(e) => {
                                const val = e.target.value;
                                // Force internal state to always be an array for consistency
                                if (typeof val === 'string') {
                                    setSelectedRhythms([val as RhythmType]);
                                } else {
                                    setSelectedRhythms(val as RhythmType[]);
                                }
                            }}
                            renderValue={(selected) => (
                                <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                                    {/* Handle rendering for both string (Group) and array (Single) values */}
                                    {(Array.isArray(selected) ? selected : [selected]).map((value) => {
                                        if (!value) return null; // Guard against empty string
                                        return (
                                            <Chip
                                                key={value}
                                                label={t(`misc_rhythm_${value}`)}
                                                size="small"
                                            />
                                        )
                                    })}
                                </Box>
                            )}
                        >
                            {ALL_RHYTHM_TYPES.map(r =>
                                <MenuItem key={r} value={r}>
                                    {t(`misc_rhythm_${r}`)}
                                </MenuItem>
                            )}
                        </Select>
                    </FormControl>

                    <Button
                        variant="contained"
                        size="large"
                        fullWidth
                        // STRICTER DISABLE LOGIC
                        disabled={
                            files.length === 0 ||
                            // isAnalyzing ||
                            selectedRhythms.length === 0 ||
                            (mode === 'GROUP' && selectedRhythms.length !== 1) ||
                            (mode === 'SINGLE' && files.length !== 1)
                        }
                        onClick={runAnalysis}
                        startIcon={<Play />}
                        sx={{ mt: 2 }}
                    >
                        {isAnalyzing ? '...' : t('config_submit')}
                    </Button>
                </Stack>
            </ConfigurationBlock>
        </>
    );
}

export default ConfigurationPage;