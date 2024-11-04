import {createBrowserRouter, RouterProvider} from "react-router-dom";
import {createTheme, CssBaseline, ThemeProvider, useColorScheme} from "@mui/material";
import {useMemo} from "react";
import {useTenantID} from "./client/common.ts";
import {Layout} from "./Layout.tsx";
import {DashboardPage} from "./pages/DashboardPage.tsx";
import {TransactionsPage} from "./pages/TransactionsPage.tsx";
import {APIPlayground} from "./pages/APIPlayground.tsx";
import {LocaleProvider} from "./providers/LocaleProvider.tsx";
import {UserProfilePage} from "./pages/UserProfilePage.tsx";
import {Descope, useSession, useUser} from "@descope/react-sdk";
import {SpinnerBlock} from "./components/SpinnerBlock.tsx";

function WithTheme({children}: any) {
    const theme = useMemo(() => createTheme({
        colorSchemes: {dark: true, light: true},
        cssVariables: {
            colorSchemeSelector: 'class'
        }
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
            {children}
        </ThemeProvider>
    )
}

export function App() {
    const {isAuthenticated, isSessionLoading} = useSession();
    const {isUserLoading} = useUser();
    const tenantID = useTenantID();
    const {mode} = useColorScheme()

    const router = useMemo(() => createBrowserRouter([
        {
            id: "root",
            path: "/",
            element: <Layout/>,
            children: [
                {
                    index: true,
                    element: <DashboardPage/>,
                },
                {
                    id: "transactions",
                    path: "/transactions",
                    element: <TransactionsPage/>,
                },
                {
                    id: "userProfile",
                    path: "/user/profile",
                    element: <UserProfilePage/>,
                },
                {
                    id: "settings",
                    path: "/settings",
                    element: <h1>Settings</h1>,
                },
                {
                    id: "api",
                    path: "/api",
                    element: <APIPlayground/>,
                },
            ],
        },
    ]), [tenantID])


    if (isSessionLoading || isUserLoading) {
        return (
            <SpinnerBlock open={true}/>
        );
    }

    if (!isAuthenticated) {
        let descopeOpts: any = {
            flowId: "custom-sign-in",
            tenant: tenantID,
        }
        if (mode === "dark" || mode === "light") {
            descopeOpts.theme = mode
        }
        return (
            <Descope {...descopeOpts}/>
        );
    }

    return (
        <WithTheme>
            <CssBaseline/>
            <LocaleProvider>
                <LocaleProvider>
                    <RouterProvider router={router}/>
                </LocaleProvider>
            </LocaleProvider>
        </WithTheme>
    )
}
