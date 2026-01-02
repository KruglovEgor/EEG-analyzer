import {Card, CardContent, Typography} from "@mui/material";
import * as React from "react";

interface ResultsBlockProps {
    headerText: string;
    headerTextColor?: string;
    // 💡 This type is perfect for accepting any React element, including icons
    headerIcon: React.ReactNode;
    children?: React.ReactNode;
}

const ResultsBlock = ({ headerText, headerTextColor, headerIcon, children }: ResultsBlockProps) => {
    return (
        <Card>
            <CardContent>
                <Typography variant="h6" gutterBottom display="flex" alignItems="center" color={headerTextColor}>
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

export default ResultsBlock;