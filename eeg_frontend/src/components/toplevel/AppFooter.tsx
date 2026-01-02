import * as React from 'react';
import {Box} from "@mui/system";
import { Typography } from '@mui/material';

interface AppFooterProps {
    footerText: string;
}

const AppFooter = ({ footerText }: AppFooterProps) => (
    <Box component="footer" sx={{ py: 2, borderTop: '1px solid rgba(255, 255, 255, 0.08)', textAlign: 'center', bgcolor: 'background.paper', position: 'fixed', bottom: 0, width: '100%', zIndex: 1200 }}>
        <Typography variant="caption" color="text.secondary">
            { footerText }
        </Typography>
    </Box>
);

export default AppFooter;