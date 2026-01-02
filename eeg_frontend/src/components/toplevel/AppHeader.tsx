import * as React from 'react';
import {AppBar, Button, MenuItem, Select, type SelectChangeEvent, Toolbar, Typography} from "@mui/material";
import {Brain, RefreshCw, Globe} from "lucide-react";
import {useTranslation} from "react-i18next";

interface AppHeaderProps {
    headerText: string;
    resetText: string;
    onClick?: () => void;
}

const AppHeader = ({ headerText, resetText, onClick }: AppHeaderProps) => {
    const { i18n } = useTranslation();
    const currentLanguage = i18n.language || "ru";
    const handleLanguageChange = (event: SelectChangeEvent<string>) => {
        const newLanguage = event.target.value;
        // i18n.changeLanguage is the i18next function to switch locales
        i18n.changeLanguage(newLanguage);
    };

    return (
        <AppBar position="fixed" color="transparent" elevation={0} sx={{ borderBottom: '1px solid rgba(255,255,255,0.08)', bgcolor: 'background.default', zIndex: 1200 }}>
            <Toolbar>
                <Brain
                    color="#00e5ff"
                    style={{ marginRight: 12 }}
                    onClick={onClick}
                />
                <Typography
                    variant="h5"
                    color="text.primary"
                    sx={{ flexGrow: 1 }}
                    onClick={onClick}
                >
                    {headerText}
                </Typography>

                <Select
                    value={currentLanguage} // Current language from i18n
                    onChange={handleLanguageChange}
                    variant="outlined"
                    size="small"
                    startAdornment={<Globe size={24} style={{ marginRight: 4, color: 'text.secondary' }} />}
                    sx={{
                        marginRight: 2,
                        color: 'text.primary',
                        '.MuiOutlinedInput-notchedOutline': { borderColor: 'rgba(255,255,255,0.2)' },
                        '&:hover .MuiOutlinedInput-notchedOutline': { borderColor: 'primary.main' },
                        '.MuiSvgIcon-root': { color: 'text.primary' } // Arrow icon color
                    }}
                >
                    <MenuItem value="en">English</MenuItem>
                    <MenuItem value="ru">Русский</MenuItem>
                    {/* Add more languages here as needed */}
                </Select>

                <Button startIcon={<RefreshCw size={16}/>} onClick={() => window.location.reload()}>
                    {resetText}
                </Button>
            </Toolbar>
        </AppBar>
    );
}

export default AppHeader;