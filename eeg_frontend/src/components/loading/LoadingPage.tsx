import React from 'react';
import { Box, Typography, CircularProgress, keyframes } from "@mui/material";
import { Loader2Icon } from "lucide-react";
import { t } from "i18next";
import ResultsBlock from "../result/ResultsBlock.tsx";

interface LoadingPageProps {
    /** Optional custom message, defaults to 'misc_loading' translation */
    message?: string;
    /** Optional subtext for more context (e.g., "This may take a minute") */
    subtext?: string;
}

// Define a smooth rotation for the header icon
const spin = keyframes`
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
`;

const LoadingPage = ({ message, subtext }: LoadingPageProps) => {
    return (
        <ResultsBlock
            headerText={t('loading_header')}
            // Animated Lucide icon for the header
            headerIcon={
                <Box sx={{
                    display: 'flex',
                    animation: `${spin} 2s linear infinite`,
                    color: 'primary.main'
                }}>
                    <Loader2Icon size={20} />
                </Box>
            }
        >
            <Box sx={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                py: 4, // Add vertical padding for breathing room
                gap: 2
            }}>
                {/* Main indeterminate spinner */}
                <CircularProgress size={64} thickness={4} />

                <Box sx={{ textAlign: 'center' }}>
                    <Typography variant="h6" color="text.primary">
                        {message ?? t('loading_analysis_mainText')}
                    </Typography>

                    {subtext && (
                        <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
                            {subtext}
                        </Typography>
                    )}
                </Box>
            </Box>
        </ResultsBlock>
    );
}

export default LoadingPage;