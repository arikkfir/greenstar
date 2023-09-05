import {useEffect} from 'react'
import {GrowthBook, GrowthBookProvider} from "@growthbook/growthbook-react";
import {useTheme} from "./hooks/theme.tsx";
import {Descope, useSession, useUser} from "@descope/react-sdk";
import {PageLoader} from "./PageLoader.tsx";
import {ThemeProvider} from "@mui/system";
import {CssBaseline} from "@mui/material";
import {SnackbarProvider} from "notistack";
import {ApolloWrapper} from "./ApolloWrapper.tsx";
import {AppLayout} from "./AppLayout.tsx";
import {AppRoutes} from "./AppRoutes.tsx";

interface AppProps {
    tenant: string
    growthBook: GrowthBook
}

function App({tenant, growthBook}: AppProps) {
    const theme = useTheme()

    const {isAuthenticated, isSessionLoading} = useSession();
    const {user, isUserLoading} = useUser()

    useEffect(() => {
        // noinspection JSIgnoredPromiseFromCall
        growthBook.loadFeatures({autoRefresh: true});
    }, [growthBook]);

    useEffect(() => {
        if (user && isAuthenticated) {
            growthBook.setAttributes({
                id: user.userId,
                tenant: tenant,
            });
        }
    }, [user, isAuthenticated, growthBook, tenant])

    if (isSessionLoading || isUserLoading) {
        return (
            <PageLoader/>
        )
    }

    if (!isAuthenticated) {
        return (
            <Descope flowId="authenticate" tenant={tenant} theme={theme.palette.mode}/>
        )
    }

    return (
        <GrowthBookProvider growthbook={growthBook}>
            <ThemeProvider theme={theme}>
                <CssBaseline/>
                <SnackbarProvider>
                    <ApolloWrapper>
                        <AppLayout tenant={tenant}>
                            <AppRoutes tenant={tenant}/>
                        </AppLayout>
                    </ApolloWrapper>
                </SnackbarProvider>
            </ThemeProvider>
        </GrowthBookProvider>
    )
}

export default App
