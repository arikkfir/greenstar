import {Outlet} from "react-router-dom";
import {Box, useTheme} from "@mui/material";
import {TopBar} from "./components/layout/topbar/TopBar.tsx";
import {NavDrawer} from "./components/layout/NavDrawer.tsx";
import {Footer} from "./components/layout/Footer.tsx";
import {useSession, useUser} from "@descope/react-sdk";
import {SpinnerBlock} from "./components/SpinnerBlock.tsx";

export function Layout() {
    const theme = useTheme();
    const {isSessionLoading} = useSession();
    const {isUserLoading} = useUser();

    if (isSessionLoading || isUserLoading) {
        return (
            <SpinnerBlock open={true}/>
        );
    }

    return (
        <Box sx={{
            p: 0,
            m: 0,
            width: '100vw',
            height: '100vh',
            display: 'flex',
            alignItems: 'stretch',
            flexDirection: 'column',
        }}>
            <Box component="header" sx={{flexGrow: 0, flexShrink: 0, display: 'flex'}}>
                <TopBar/>
            </Box>
            <Box sx={{flexGrow: 1, display: 'flex', flexDirection: 'row', alignItems: 'stretch'}}>
                <Box component="nav" sx={{
                    flexGrow: 0,
                    flexShrink: 0,
                    p: 1,
                    borderRight: 1,
                    borderRightColor: theme.palette.grey[300]
                }}>
                    <NavDrawer/>
                </Box>
                <Box component="main" sx={{display: 'flex', alignItems: 'stretch', flexGrow: 1, flexShrink: 1, p: 2}}>
                    <Outlet/>
                </Box>
            </Box>
            <Box component="footer" sx={{
                flexGrow: 0,
                flexShrink: 0,
                display: 'flex',
                flexDirection: 'row',
                borderTop: 1,
                borderTopColor: theme.palette.grey[300]
            }}>
                <Footer/>
            </Box>
        </Box>
    )
}
