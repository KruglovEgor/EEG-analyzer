import * as React from 'react';
import {Container, Box, Stack} from "@mui/material";

interface AppMainContainerProps {
    children?: React.ReactNode;
}

const AppMainContainer = ({children}: AppMainContainerProps) => (
    <Box component="main" sx={{ flexGrow: 1, pt: '80px', pb: '80px', px: 3, display: 'flex', justifyContent: 'center' }}>
        <Container maxWidth="lg">
            <Stack spacing={4}>
                {children}
            </Stack>
        </Container>
    </Box>
);

export default AppMainContainer;