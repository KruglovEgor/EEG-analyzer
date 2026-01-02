import * as React from 'react';
import { Card, CardContent, Typography } from "@mui/material";

interface ConfigurationBlockProps {
    headerText: string;
    // 💡 This type is perfect for accepting any React element, including icons
    headerIcon: React.ReactNode;
    children?: React.ReactNode;
}

const ConfigurationBlock = ({ headerText, headerIcon, children }: ConfigurationBlockProps) => {

    return (
        <Card>
            <CardContent>
                <Typography variant="h6" gutterBottom display="flex" alignItems="center">
                    {/* 🚨 FIX: RENDER THE PASSED headerIcon PROP 🚨 */}
                    {/* We wrap it in a Span/Box to control spacing consistently */}
                    <span style={{ marginRight: 8, display: 'flex' }}>
                        {headerIcon}
                    </span>
                    {headerText}
                </Typography>
                {children}
            </CardContent>
        </Card>
    );
}

export default ConfigurationBlock;