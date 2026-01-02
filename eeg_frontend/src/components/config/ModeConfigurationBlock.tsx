// src/components/ModeConfigurationBlock.tsx

import * as React from 'react';
import {Paper, Tabs, Tab, Box, Typography, Tooltip} from "@mui/material";
import { useTranslation } from 'react-i18next';
import type { AnalysisMode } from '../../types/configTypes'; // 👈 Import the constant and type

// Define the component props
interface ModeConfigurationProps {
    mode?: AnalysisMode,
    // Callback function to send the new mode (GROUP or SINGLE) back to the parent
    onModeChange: (mode: AnalysisMode) => void;
}

const ModeConfigurationBlock: React.FC<ModeConfigurationProps> = (
    { mode = 'GROUP', onModeChange }
) => {
    const { t } = useTranslation();

    // Local state for the selected tab index (0 for Group, 1 for Single File)
    const [tabValue, setTabValue] = React.useState(0);

    const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
        setTabValue(newValue);

        // Map the tab index to the AnalysisMode constant and call the parent callback
        const newMode: AnalysisMode = newValue === 0 ? 'GROUP' : 'SINGLE';
        onModeChange(newMode);
    };

    // Determine the descriptive text based on the current tab value
    // We use dynamic translation keys here: 'mode.group' or 'mode.single'
    const descriptionKey = tabValue === 0 ? 'config_mode_description_group' : 'config_mode_description_single';

    return (
        <Paper sx={{ overflow: 'hidden' }}>
            <Tabs
                value={tabValue}
                onChange={handleTabChange}
                variant="fullWidth"
                indicatorColor="primary"
                textColor="primary"
            >
                {/* 1. Use translation for the Tab labels */}
                <Tooltip title={t('config_mode_description_group')}>
                    <Tab label={t('config_mode_label_group')} />
                </Tooltip>
                <Tooltip title={t('config_mode_description_single')}>
                    <Tab label={t('config_mode_label_single')} />
                </Tooltip>
            </Tabs>
            <Box sx={{ p: 2 }}>
                <Typography variant="caption" color="text.secondary">
                    {/* 2. Use translation for the description */}
                    {t(descriptionKey)}
                </Typography>
            </Box>
        </Paper>
    );
};

export default ModeConfigurationBlock;