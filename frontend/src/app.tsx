import CssBaseline from "@mui/material/CssBaseline";
import {ThemeProvider} from "@mui/material/styles";
import {Route, Routes} from "react-router-dom";
import {PageLoader} from "./components/layout/PageLoader";
import {useTheme} from "./hooks/theme";
import {HomePage} from "./pages/HomePage";
import {NotFoundPage} from "./pages/NotFoundPage";
import {ApolloWrapper} from "./apollo-wrapper";
import {AppLayout} from "./app-layout";
import {Descope, useSession, useUser} from "@descope/react-sdk";
import {Tenants} from "./pages/Tenants";
import {GrowthBook, GrowthBookProvider} from "@growthbook/growthbook-react";
import {useEffect} from "react";

interface AppProps {
    tenant: string
    growthBook: GrowthBook
}

export function App({tenant, growthBook}: AppProps) {
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
            <Descope flowId="welcome" tenant={tenant} theme={theme.palette.mode}/>
        )
    }

    return (
        <GrowthBookProvider growthbook={growthBook}>
            <ThemeProvider theme={theme}>
                <CssBaseline/>
                <ApolloWrapper tenant={tenant}>
                    <AppLayout tenant={tenant}>
                        <Routes>
                            <Route path="/" element={<HomePage/>}/>
                            {tenant === "global" && (<Route path="/tenants" element={<Tenants/>}/>)}
                            <Route path="*" element={<NotFoundPage/>}/>
                        </Routes>
                    </AppLayout>
                </ApolloWrapper>
            </ThemeProvider>
        </GrowthBookProvider>
    )
}
