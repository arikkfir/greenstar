import {Box, Stack, useTheme} from "@mui/material";
import React from "react";
import {Outlet,} from "react-router-dom";
import {Copyright} from "./Copyright";
import {AppDrawer} from "./Drawer";

export interface ContentProps {
    adminAPIURL: string
    operationsAPIURL: string
    publicAPIURL: string
}

export function Content({adminAPIURL, operationsAPIURL, publicAPIURL}: ContentProps) {
    const theme = useTheme()
    return (
        <Stack direction="row" spacing={0} sx={{flexGrow: 1, overflow: 'auto'}}>
            <AppDrawer adminAPIURL={adminAPIURL} operationsAPIURL={operationsAPIURL} publicAPIURL={publicAPIURL}/>
            <Stack direction="column" sx={{
                flexGrow: 1,
                overflow: 'auto',
                "&::-webkit-scrollbar": {
                    width: 20
                },
                "&::-webkit-scrollbar-track": {
                    // backgroundColor: "orange"
                },
                "&::-webkit-scrollbar-thumb": {
                    backgroundColor: theme.palette.primary.main,
                    borderRadius: 0
                }
            }}>
                <Box sx={{p: 3}}>
                    <Outlet/>
                </Box>
                <Copyright/>
            </Stack>
        </Stack>
    )
}
