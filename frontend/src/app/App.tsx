import {useAuth0} from "@auth0/auth0-react";
import {Backdrop, CircularProgress, PaletteMode, Stack, useMediaQuery} from "@mui/material";
import CssBaseline from "@mui/material/CssBaseline";
import {createTheme, ThemeProvider} from "@mui/material/styles";
import React, {useEffect, useState} from 'react';
import {createBrowserRouter, createRoutesFromElements, Route, RouterProvider,} from "react-router-dom";
import {RouteProps} from "react-router/dist/lib/components";
import {About} from "./About";
import {Content} from "./Content";
import {Dashboard} from "./Dashboard";
import {LoginError} from "./LoginError";
import {Organization} from "./Organization";
import {TopAppBar} from "./TopAppBar";
import {Transactions} from "./Transactions";

export interface AppProps {
    organization: Organization
    adminAPIURL: string
    operationsAPIURL: string
    publicAPIURL: string
}


export function App({organization, adminAPIURL, publicAPIURL, operationsAPIURL}: AppProps) {
    useEffect(() => {
        document.title = "Greenstar"
    }, []);

    //
    // Build the theme of the app, based on either the OS current theme, or a user-specified preference, if one exists
    // TODO: get/store palette mode user preference in local storage
    //
    const osPaletteMode = useMediaQuery('(prefers-color-scheme: dark)') ? 'dark' : 'light';
    const [preferredPaletteMode] = useState<PaletteMode | null>(null)
    const theme = React.useMemo(
        () => createTheme({
            palette: {mode: preferredPaletteMode || osPaletteMode},
        }),
        [osPaletteMode, preferredPaletteMode],
    );

    // Manage user state
    const {
        error: authError,
        isAuthenticated,
        isLoading,
        // user,
        // getAccessTokenSilently,
        // getAccessTokenWithPopup,
        // getIdTokenClaims,
    } = useAuth0();

    //
    // Build the routes based on login status
    //
    let routes: React.ReactElement<RouteProps, typeof Route>[]
    if (authError) {
        routes = [<Route key="error" path="/" element={<LoginError error={authError}/>}/>]
    } else if (isAuthenticated) {
        routes = [
            <Route key="dashboard" path="/" element={<Dashboard/>}/>,
            <Route key="transactions" path="/transactions" element={<Transactions/>}/>,
            <Route key="about" path="/about" element={<About/>}/>,
        ]
    } else {
        routes = []
    }

    const router = createBrowserRouter(
        createRoutesFromElements(
            <Route path="/" element={
                <Content adminAPIURL={adminAPIURL}
                         operationsAPIURL={operationsAPIURL}
                         publicAPIURL={publicAPIURL}/>}>
                {routes}
            </Route>)
    )

    return (
        <ThemeProvider theme={theme}>
            <CssBaseline/>
            <Backdrop sx={{color: '#fff', zIndex: theme => theme.zIndex.drawer + 1}} open={isLoading}>
                <CircularProgress color="inherit"/>
            </Backdrop>
            <Stack direction="column"
                   spacing={0}
                   position="fixed"
                   sx={{width: '100vw', height: "100vh", overflow: 'hidden'}}>
                <TopAppBar organization={organization}/>
                <RouterProvider router={router}/>
            </Stack>
        </ThemeProvider>
    )
}
