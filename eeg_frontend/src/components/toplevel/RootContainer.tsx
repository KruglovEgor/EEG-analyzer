// RootContainer.tsx
import { Box } from "@mui/system";
import * as React from "react"; // Import React for the type definition

// Define the component's props interface
interface RootContainerProps {
    // 1. Rename the prop from 'content' to 'children'
    // 2. Use the standard React type for elements/nodes/primitives passed as children
    children: React.ReactNode;
}

// Destructure 'children' from props
const RootContainer: React.FC<RootContainerProps> = ({ children }) => (
    <Box
        // These styles are great for setting up a full-page container
        sx={{
            minHeight: '100vh',
            bgcolor: 'background.default', // Assumes you have a theme with a 'background.default' color
            width: '100vw',
            overflowX: 'hidden',
            display: 'flex',
            flexDirection: 'column'
        }}
    >
        {/* Render the children here */}
        {children}
    </Box>
);

export default RootContainer;