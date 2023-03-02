import {Backdrop, CircularProgress, PaletteMode, Stack, useMediaQuery} from "@mui/material";
import CssBaseline from "@mui/material/CssBaseline";
import {createTheme, ThemeProvider} from "@mui/material/styles";
import React, {useEffect, useState} from 'react';

import {createBrowserRouter, createRoutesFromElements, Route, RouterProvider,} from "react-router-dom";
import {RouteProps} from "react-router/dist/lib/components";
import {About} from "./About";

import {TopAppBar} from "./AppBar";
import {Content} from "./Content";
import {Dashboard} from "./Dashboard";
import {Forbidden} from "./Forbidden";
import {LoginError} from "./LoginError";
import {Transactions} from "./Transactions";
import {LoginStatus, User} from "./User";

export interface AppProps {
    environment: string
    version: string
    userInfoURL: string
    loginURL: string
    adminAPIURL: string
    operationsAPIURL: string
    publicAPIURL: string
}

export function App({
                        environment,
                        version,
                        userInfoURL,
                        loginURL,
                        adminAPIURL,
                        operationsAPIURL,
                        publicAPIURL
                    }: AppProps) {
    useEffect(() => {
        document.title = `Greenstar (v${version}/${environment})`;
    }, [version, environment]);
    console.info(`Loading application v${version} (${environment})`)

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

    //
    // Manage logged-in user information state by obtaining the user information from the server
    // React to the server response as follows:
    // - HTTP 2xx: obtain the user information and store it in state
    // - HTTP 401: means that the user is not logged in
    // - HTTP 403: means that the user is not allowed to log-in
    // - Otherwise: set the user state to "error"
    //
    const [loginStatus, setLoginStatus] = useState<LoginStatus>("pending")
    useEffect(() => {
        // TODO: try to prevent excessive (double) execution (might be due to devMode)
        console.info("Requesting user information from server")
        fetch(userInfoURL, {cache: 'no-store', credentials: 'include'})
            .then(response => {
                if (response.status >= 200 && response.status <= 299) {
                    response.json().then(data => new User(data)).then(setLoginStatus)
                } else if (response.status === 401) {
                    // TODO: Use cookies to prevent excessive redirects (maybe server can do that too?)
                    window.location.href = loginURL
                } else if (response.status === 403) {
                    setLoginStatus("forbidden")
                } else {
                    console.error("Server failed to provide user-info: \n" +
                        `\tRequest type:        ${response.url}\n`,
                        `\tResponse Status:     ${response.status} ${response.statusText}\n` +
                        `\tResponse headers:    ${response.headers}\n`,
                        `\tResponse redirected? ${response.redirected}\n`,
                        `\tResponse type:       ${response.type}\n`,
                        `\tResponse data:       ${response.text()}\n`)
                    setLoginStatus("error")
                }
            })
            .catch(reason => {
                console.error("Failed to send user-info request to server: ", reason)
                setLoginStatus("error");
            })
    }, [loginURL, userInfoURL])

    //
    // Build the routes based on login status
    //
    let routes: React.ReactElement<RouteProps, typeof Route>[]
    let user: User | undefined
    if (loginStatus === "forbidden") {
        routes = [<Route key="forbidden" path="/" element={<Forbidden/>}/>]
    } else if (loginStatus === "error") {
        routes = [<Route key="error" path="/" element={<LoginError/>}/>]
    } else if (loginStatus === "pending") {
        routes = []
    } else {
        user = loginStatus
        routes = [
            <Route key="dashboard" path="/" element={<Dashboard/>}/>,
            <Route key="transactions" path="/transactions" element={<Transactions/>}/>,
            <Route key="about" path="/about" element={<About/>}/>,
        ]
    }

    const router = createBrowserRouter(
        createRoutesFromElements(<Route path="/"
                                        element={<Content adminAPIURL={adminAPIURL}
                                                          operationsAPIURL={operationsAPIURL}
                                                          publicAPIURL={publicAPIURL}/>}>{routes}</Route>)
    )

    return (
        <ThemeProvider theme={theme}>
            <CssBaseline/>
            <Backdrop sx={{color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1}} open={loginStatus === "pending"}>
                <CircularProgress color="inherit"/>
            </Backdrop>
            <Stack direction="column"
                   spacing={0}
                   position="fixed"
                   sx={{width: '100vw', height: "100vh", overflow: 'hidden'}}>
                <TopAppBar user={user}/>
                <RouterProvider router={router}/>
            </Stack>
        </ThemeProvider>
    );
}
