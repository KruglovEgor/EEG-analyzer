import ResultsBlock from "../result/ResultsBlock.tsx";
import {Button, SvgIcon, Typography} from "@mui/material";
import {CircleAlertIcon, HomeIcon} from "lucide-react";
import {t} from "i18next";

interface ErrorPageProps {
    message: string;
    statusCode?: number;
    onGoHome?: () => void;
}

const ErrorPage = ({ message, statusCode, onGoHome = () => {} }: ErrorPageProps) => {
    return (
        <>
            <ResultsBlock
                headerText={t('errorResult_header')}
                headerTextColor="error"
                headerIcon={<SvgIcon color="error"><CircleAlertIcon /></SvgIcon>}
            >
                {statusCode && <Typography variant="h1" color="error">{statusCode}</Typography>}
                <Typography color="error">{message}</Typography>
                <Button
                    variant="contained"
                    size="large"
                    fullWidth
                    onClick={onGoHome}
                    startIcon={<HomeIcon />}
                    sx={{mt: 2}}
                >{ t('misc_button_goHome') }</Button>
            </ResultsBlock>
        </>
    );
}

export default ErrorPage;