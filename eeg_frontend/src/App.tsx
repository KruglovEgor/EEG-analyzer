import React, {useCallback, useEffect, useMemo, useRef, useState} from "react";
import {darkTheme} from "./theme/theme.ts";
import {ThemeProvider, CssBaseline} from "@mui/material";
import RootContainer from "./components/toplevel/RootContainer.tsx";
import AppHeader from "./components/toplevel/AppHeader.tsx";
import AppFooter from "./components/toplevel/AppFooter.tsx";
import AppMainContainer from "./components/toplevel/AppMainContainer.tsx";
import {useTranslation} from "react-i18next";
import ConfigurationPage from "./components/config/ConfigurationPage.tsx";
import ResultsPage from "./components/result/ResultsPage.tsx";
import type {AppState} from "./types/appState.ts";
import ErrorPage from "./components/errorstate/ErrorPage.tsx";
import {createApiClient, eegFormDataToRequest} from "./communication/apiClient.ts";
import type {EEGAnalysisFormData} from "./types/configTypes.ts";
import LoadingPage from "./components/loading/LoadingPage.tsx";

export default function App(){
    const { t, i18n } = useTranslation();
    const [appState, setAppState]
        = useState<AppState>({ status: 'IDLE' });

    // 1. Store the controller in a ref
    const abortControllerRef = useRef<AbortController | null>(null);

    const apiClient = useMemo(() => createApiClient({
        baseUrl: import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:3000',
        analysisEndpoint: import.meta.env.VITE_API_ANALYSIS_ENDPOINT ?? '/analyze',
    }), []);

    useEffect(() => {
        const statusTitles = {
            'IDLE': t('titles_idle'),
            'LOADING': t('titles_loading'),
            'SUCCESS': t('titles_success'),
            'ERROR': t('titles_error')
        };

        if (document) document.title = statusTitles[appState.status];
    }, [appState]);

    const runAnalysis = useCallback(async (formData: EEGAnalysisFormData) => {
        const request = eegFormDataToRequest(formData);

        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
            abortControllerRef.current = null;
        }
        const abController = new AbortController();
        abortControllerRef.current = abController;

        setAppState({ status: 'LOADING' });

        try {
            const response = await apiClient.analyze(request, abController.signal);
            // if (response.analysisId !== request.analysisId) {
            //     throw new Error('request ID and response ID do not match');
            // }
            setAppState({ status: 'SUCCESS', response: response });
        } catch (e) {
            if (e instanceof DOMException && e.name === 'AbortError') {
                setAppState({ status: 'IDLE' })
                return;
            }
            setAppState({ status: 'ERROR', message: String(e) });
        } finally {
            if (abortControllerRef.current === abController) {
                abortControllerRef.current = null;
            }
        }
    }, [apiClient]);

    const onGoHome = useCallback(() => {
        // If there's an analysis running, cancel it and go home
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
            abortControllerRef.current = null;
        }

        setAppState({ status: 'IDLE' });
    }, []);

    return (
        <ThemeProvider theme={darkTheme}>
            <CssBaseline/>
            <RootContainer>
                <AppHeader
                    headerText={t('header_appName')}
                    resetText={t('header_resetBtn')}
                    onClick={onGoHome}
                />
                <AppMainContainer>
                    {appState.status === 'IDLE' &&
                        <ConfigurationPage onRunAnalysis={runAnalysis} />
                    }
                    {appState.status === 'LOADING' &&
                        <LoadingPage message={t('loading_analysis_mainText')}/>
                    }
                    {appState.status === 'SUCCESS' &&
                        <ResultsPage
                            analysis={appState.response}
                            onGoHome={onGoHome}
                        />
                    }
                    {appState.status === 'ERROR' &&
                        <ErrorPage
                            message={appState.message}
                            onGoHome={onGoHome}
                        />
                    }
                </AppMainContainer>
                <AppFooter footerText={t('footer_copyrightText')} />
            </RootContainer>
        </ThemeProvider>
    );
}