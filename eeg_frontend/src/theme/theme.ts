import {createTheme} from "@mui/material/styles";

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
        primary: { main: '#00e5ff' },
        secondary: { main: '#f50057' },
        warning: { main: '#ffc400' },
        background: { default: '#0a0e17', paper: '#111827' },
        text: { primary: '#e2e8f0', secondary: '#94a3b8' },
    },
    typography: {
        fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
        h6: { fontWeight: 600 },
        h5: { fontWeight: 700 },
    },
    components: {
        MuiPaper: {
            styleOverrides: {
                root: { borderRadius: 12, border: '1px solid rgba(255, 255, 255, 0.08)' },
            },
        },
    },
});

export { darkTheme };