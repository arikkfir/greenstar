import { Box, useTheme } from "@mui/material"
import { Descope, useSession, useUser } from "@descope/react-sdk"
import { useTenantID } from "./hooks/tenant.ts"
import { SpinnerBlock } from "./components/SpinnerBlock.tsx"
import { LocaleContext } from "./providers/LocaleProvider.tsx"
import { Outlet, Route, Routes } from "react-router"
import { DashboardPage } from "./pages/DashboardPage.tsx"
import { TransactionsPage } from "./pages/TransactionsPage.tsx"
import { APIPlaygroundPage } from "./pages/APIPlayground.tsx"
import { SettingsPage } from "./pages/SettingsPage.tsx"
import { UserProfilePage } from "./pages/UserProfilePage.tsx"
import { ThemeButton } from "./components/layout/topbar/ThemeButton.tsx"
import { TopBar } from "./components/layout/topbar/TopBar.tsx"
import { Footer } from "./components/layout/Footer.tsx"
import { useContext } from "react"

export function App() {
    const locale = useContext(LocaleContext)
    const { isAuthenticated, isSessionLoading } = useSession()
    const { isUserLoading } = useUser()
    const tenantID = useTenantID()
    const theme = useTheme()

    if (!locale.resolved || isSessionLoading || isUserLoading) {
        return <SpinnerBlock open={true} />
    }

    if (!isAuthenticated) {
        let descopeOpts: any = {
            flowId: "custom-sign-in",
            tenant: tenantID,
            theme: theme.palette.mode,
        }
        return (
            <Box
                sx={{
                    p: 0,
                    m: 0,
                    width: "100vw",
                    height: "100vh",
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                    justifyContent: "center",
                }}
            >
                <ThemeButton />
                <Descope {...descopeOpts} />
            </Box>
        )
    }

    return (
        <Routes>
            <Route element={<RootLayout />}>
                <Route index element={<DashboardPage />} />
                <Route path="transactions" element={<TransactionsPage />} />
                <Route path="api" element={<APIPlaygroundPage />} />
                <Route path="settings" element={<SettingsPage />} />
                <Route path="profile" element={<UserProfilePage />} />
            </Route>
        </Routes>
    )
}

export function RootLayout() {
    const theme = useTheme()
    return (
        <Box
            sx={{
                backgroundColor: theme.palette.background.default,
                p: 0,
                m: 0,
                width: "100vw",
                height: "100vh",
                overflow: "hidden",
                display: "flex",
                flexDirection: "column",
                alignItems: "stretch",
                alignContent: "stretch",
                gap: 0,
            }}
        >
            <TopBar sx={{ flexGrow: 0, flexShrink: 0 }} />
            <Box
                component="main"
                sx={{
                    flexGrow: 1,
                    flexShrink: 1,
                    display: "flex",
                    overflow: "hidden",
                    flexDirection: "column",
                    justifyContent: "center",
                    alignItems: "stretch",
                    alignContent: "stretch",
                }}
            >
                <Outlet />
            </Box>
            <Footer sx={{ flxGrow: 0, flexShrink: 0 }} />
        </Box>
    )
}
