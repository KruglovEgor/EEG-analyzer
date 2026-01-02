import {useTranslation} from "react-i18next";
import type {EEGFileConfig} from "../../types/configTypes.ts";
import React, {useState} from "react";
import {Box, Grid, IconButton, Paper, TextField, Typography, Tooltip} from "@mui/material";
import {FileText, Settings, Trash2} from "lucide-react";


const FileItem = ({
                      file,
                      onUpdate,
                      onRemove
                  }: {
    file: EEGFileConfig,
    onUpdate: (id: string, field: keyof EEGFileConfig, val: any) => void,
    onRemove: (id: string) => void
}) => {
    const { t } = useTranslation();
    const [expanded, setExpanded] = useState(false);

    const formatFileSize = (bytes: number | undefined): string => {
        if (bytes === undefined || bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
    };

    return (
        <Paper sx={{ p: 2, mb: 1, backgroundColor: 'rgba(255,255,255,0.02)' }}>
            <Grid container spacing={2} alignItems="center">
                <Grid component="div">
                    <FileText size={20} color="#00e5ff" />
                </Grid>
                <Grid component="div">
                    <Tooltip title={`ID: #${file.id}`}>
                        <Typography variant="body2" noWrap sx={{ fontWeight: 500 }}>
                            {file.filename}
                        </Typography>
                    </Tooltip>
                    <Typography variant="caption" color="text.secondary">
                        {formatFileSize(file.rawFile?.size)}
                    </Typography>
                </Grid>
                <Grid component="div">
                    <TextField
                        fullWidth
                        size="small"
                        label={t('config_file_itemExperimentName')}
                        value={file.experimentName}
                        onChange={(e) => onUpdate(file.id, 'experimentName', e.target.value)}
                        variant="outlined"
                    />
                </Grid>
                <Grid component="div" sx={{ display: 'flex', justifyContent: 'flex-end' }}>
                    <IconButton size="small" onClick={() => setExpanded(!expanded)}>
                        <Settings size={18} />
                    </IconButton>
                    <IconButton size="small" color="error" onClick={() => onRemove(file.id)}>
                        <Trash2 size={18} />
                    </IconButton>
                </Grid>
            </Grid>

            {expanded && (
                <Box sx={{ mt: 2, pt: 2, borderTop: '1px solid rgba(255,255,255,0.1)' }}>
                    <Grid container spacing={2}>
                        <Grid component="div">
                            <TextField
                                fullWidth size="small"
                                label={t('config_file_itemTimeCol')}
                                value={file.timeColumn}
                                onChange={(e) => onUpdate(file.id, 'timeColumn', e.target.value)}
                            />
                        </Grid>
                        <Grid component="div">
                            <TextField
                                fullWidth size="small"
                                label={t('config_file_itemAmplitudeCol')}
                                value={file.amplitudeColumn}
                                onChange={(e) => onUpdate(file.id, 'amplitudeColumn', e.target.value)}
                            />
                        </Grid>
                    </Grid>
                </Box>
            )}
        </Paper>
    );
};

export default FileItem;