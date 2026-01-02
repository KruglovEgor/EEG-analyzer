// src/i18n.js
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import HttpBackend from 'i18next-http-backend';

const PUBLIC_BASE_PATH = import.meta.env.BASE_URL || '/';

i18n
    // Load translation files via http (needed for production)
    .use(HttpBackend)
    // Detect user language (e.g., from browser settings)
    .use(LanguageDetector)
    // Pass i18n instance to react-i18next
    .use(initReactI18next)
    .init({
        // Fallback language if translation is missing
        fallbackLng: 'ru',
        // Debug output to console during development
        debug: true,

        // Configuration for the HTTP Backend to load translation files
        backend: {
            loadPath: `${PUBLIC_BASE_PATH}locales/{{lng}}/{{ns}}.json`,
        },

        interpolation: {
            // React already protects against XSS
            escapeValue: false,
        },

        // Configuration for Language Detector
        detection: {
            // Order of language detection (e.g., querystring, cookie, header)
            order: ['querystring', 'cookie', 'localStorage', 'sessionStorage', 'navigator', 'htmlTag'],
            // Cache the detected language
            caches: ['localStorage'],
        }
    });

export default i18n;