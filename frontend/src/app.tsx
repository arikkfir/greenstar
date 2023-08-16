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
import {GrowthBook, GrowthBookProvider, useFeatureIsOn} from "@growthbook/growthbook-react";
import {useEffect} from "react";
import {Accounts} from "./pages/Accounts";
import {APIExplorer} from "./pages/APIExplorer";
import {SnackbarProvider} from "notistack";

interface AppProps {
    tenant: string
    growthBook: GrowthBook
}

interface AppRoutesProps {
    tenant: string
}

function AppRoutes({tenant}: AppRoutesProps) {
    const showGraphQLQueryLink = useFeatureIsOn("show-graphql-console");
    return (
        <Routes>
            <Route path="/" element={<HomePage/>}/>
            <Route path="/accounts" element={<Accounts tenantID={tenant}/>}/>
            {tenant === "global" && (<Route path="/tenants" element={<Tenants/>}/>)}
            {showGraphQLQueryLink && (<Route path="/api" element={<APIExplorer/>}/>)}
            <Route path="*" element={<NotFoundPage/>}/>
        </Routes>
    )
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
