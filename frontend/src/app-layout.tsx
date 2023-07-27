import React, {PropsWithChildren} from "react";
import {Box, Stack} from "@mui/material";
import {TopAppBar} from "./components/layout/TopAppBar";
import {AppDrawer} from "./components/layout/AppDrawer";
import {Copyright} from "./components/layout/Copyright";

interface AppLayoutProps {
    tenant: string
}

export function AppLayout({children, tenant}: PropsWithChildren<AppLayoutProps>) {
    return (
        <Stack direction="column"
               spacing={0}
               position="fixed"
               sx={{width: '100vw', height: "100vh", overflow: 'hidden'}}>
            <TopAppBar/>
            <Stack direction="row" spacing={0} sx={{flexGrow: 1, overflow: 'auto'}}>
                <AppDrawer tenant={tenant}/>
                <Stack direction="column" sx={{flexGrow: 1, overflow: 'auto'}}>
                    <Box sx={{flexGrow: 1, p: 3}}>
                        {children}
                    </Box>
                    <Copyright/>
                </Stack>
            </Stack>
        </Stack>
    )
}
