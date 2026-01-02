import * as React from 'react';
import { Box, Typography } from "@mui/material";
import { Upload } from "lucide-react";
import { useTranslation } from "react-i18next";

interface FileDropzoneProps {
    // Current analysis mode (e.g., true for GROUP, false for SINGLE)
    isMultiple: boolean;
    // Function to call when files are selected or dropped
    onFileUpload: (files: FileList) => void;
    // Acceptable file types
    accept: string;
}

const FileDropzone: React.FC<FileDropzoneProps> = ({ isMultiple, onFileUpload, accept }) => {
    const { t } = useTranslation();
    const [isDragging, setIsDragging] = React.useState(false);
    const dropzoneRef = React.useRef<HTMLDivElement>(null);

    // --- 1. Handlers for Input Change ---
    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        if (event.target.files) {
            onFileUpload(event.target.files);
            // Clear the input value so the same file can be selected again
            event.target.value = '';
        }
    };

    // --- 2. Handlers for Drag Events ---

    // Prevent default behavior (which is usually to open the file)
    const handleDragOver = (event: React.DragEvent) => {
        event.preventDefault();
        event.stopPropagation();
    };

    // Set dragging state to true for visual feedback
    const handleDragEnter = (event: React.DragEvent) => {
        event.preventDefault();
        event.stopPropagation();
        setIsDragging(true);
    };

    // Set dragging state to false when leaving the zone
    const handleDragLeave = (event: React.DragEvent) => {
        event.preventDefault();
        event.stopPropagation();
        // Check if the event is truly leaving the container, not just entering a child element
        if (!event.currentTarget.contains(event.relatedTarget as Node)) {
            setIsDragging(false);
        }
    };

    // Process dropped files
    const handleDrop = (event: React.DragEvent) => {
        event.preventDefault();
        event.stopPropagation();
        setIsDragging(false); // Reset visual state

        if (event.dataTransfer.files && event.dataTransfer.files.length > 0) {
            onFileUpload(event.dataTransfer.files);
            event.dataTransfer.clearData();
        }
    };

    return (
        <Box
            ref={dropzoneRef}
            // --- Attach Drag/Drop Handlers ---
            onDragOver={handleDragOver}
            onDragEnter={handleDragEnter}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}

            // Allow the box to act as a label for file input click
            component="label"

            // --- Dynamic Styling based on Drag State ---
            sx={{
                border: isDragging ? '2px solid' : '2px dashed', // Solid border when dragging
                borderColor: isDragging ? 'primary.main' : 'rgba(255,255,255,0.2)',
                bgcolor: isDragging ? 'rgba(0,229,255,0.08)' : 'transparent', // Highlight color

                // Existing styling for layout and hover
                borderRadius: 2,
                p: 4,
                textAlign: 'center',
                cursor: 'pointer',
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                transition: 'all 0.2s ease-in-out',
                '&:hover': {
                    borderColor: 'primary.main',
                    bgcolor: 'rgba(0,229,255,0.04)',
                    transform: 'translateY(-2px)'
                }
            }}
        >
            {/* Hidden File Input */}
            <input
                type="file"
                hidden
                multiple={isMultiple}
                accept={accept}
                onChange={handleInputChange}
            />

            <Upload size={40} style={{ opacity: isDragging ? 0.9 : 0.5, marginBottom: 16 }} />

            {/* Dynamic Text */}
            <Typography color="text.secondary">
                {isDragging ? t('config_file_dropzoneDropText') : t('config_file_dropzoneClickText')}
            </Typography>
        </Box>
    );
}

export default FileDropzone;