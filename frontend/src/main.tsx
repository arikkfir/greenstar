import { StrictMode, useMemo } from "react"
import { createRoot } from "react-dom/client"
import "./index.css"
import { App } from "./App.tsx"
import { BrowserRouter } from "react-router"
import { createTheme, CssBaseline, ThemeProvider } from "@mui/material"
import { LocaleProvider } from "./providers/LocaleProvider.tsx"

function WithTheme({ children }: any) {
    const theme = useMemo(
        () =>
            createTheme({
                colorSchemes: { dark: true, light: true },
                cssVariables: {
                    colorSchemeSelector: "class",
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
            }),
        [],
    )

    return (
        <ThemeProvider theme={theme}>
            <CssBaseline />
            {children}
        </ThemeProvider>
    )
}

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <LocaleProvider>
            <BrowserRouter>
                <WithTheme>
                    <App />
                </WithTheme>
            </BrowserRouter>
        </LocaleProvider>
    </StrictMode>,
)
