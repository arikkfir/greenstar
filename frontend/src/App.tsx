import {Box, Container, Typography, useMediaQuery} from "@mui/material";

import Button from '@mui/material/Button';
import CssBaseline from "@mui/material/CssBaseline";
import {createTheme, ThemeProvider} from "@mui/material/styles";
import {Provider as JotaiProvider} from "jotai";
import React from 'react';

import TopAppBar from "./AppBar";
import {Copyright} from "./Copyright";

export default function App() {
    const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');
    const theme = React.useMemo(
        () => createTheme({
            palette: {mode: prefersDarkMode ? 'dark' : 'light'},
        }),
        [prefersDarkMode],
    );

    return (
        <ThemeProvider theme={theme}>
            <CssBaseline/>
            <JotaiProvider>
                <Box>
                    <TopAppBar/>
                    <Container>
                        <Box sx={{my: 2}}/>
                        <Box display="flex" justifyContent="center" alignItems="center">
                            <Typography variant="h4" component="h1" gutterBottom>
                                Create React App example with TypeScript
                            </Typography>
                        </Box>
                        <Box display="flex" justifyContent="center" alignItems="center">
                            <Button variant="contained">Hello World</Button>
                        </Box>
                    </Container>
                    <Copyright/>
                </Box>
            </JotaiProvider>
        </ThemeProvider>
    );
};
