import {StrictMode, useMemo} from 'react'
import {createRoot} from 'react-dom/client'
import './index.css'
import {App} from './App.tsx'
import {AuthProvider} from "@descope/react-sdk";
import {BrowserRouter} from "react-router";
import {createTheme, CssBaseline, ThemeProvider} from "@mui/material";

function WithTheme({children}: any) {
    const theme = useMemo(() => createTheme({
        colorSchemes: {dark: true, light: true},
        cssVariables: {
            colorSchemeSelector: 'class'
        },
        // components: {
        //     MuiDataGrid: {
        //         styleOverrides: {
        //             root: {
        //                 backgroundColor: 'red',
        //             },
        //         },
        //     },
        // }
    }), []);

    return (
        <ThemeProvider theme={theme}>
            <CssBaseline/>
            {children}
        </ThemeProvider>
    )
}

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <AuthProvider sessionTokenViaCookie={false} projectId={import.meta.env.VITE_DESCOPE_PROJECT_ID}>
            <BrowserRouter>
                <WithTheme>
                    <App/>
                </WithTheme>
            </BrowserRouter>
        </AuthProvider>
    </StrictMode>,
)
