import MenuIcon from "@mui/icons-material/Menu";
import {AppBar, Box, IconButton, Toolbar, Typography} from "@mui/material";
import React from "react";
import {Organization} from "./Organization";
import {Profile} from "./Profile";

export interface TopAppBarProps {
    organization: Organization
}

export function TopAppBar({organization}: TopAppBarProps) {
    return (
        <AppBar position="static" sx={{zIndex: (theme) => theme.zIndex.drawer + 1}}>
            <Toolbar>
                <IconButton size="large" edge="start" color="inherit" aria-label="open drawer" sx={{mr: 2}}>
                    <MenuIcon/>
                </IconButton>
                <Typography variant="h6" noWrap>
                    Greenstar
                </Typography>
                <Box sx={{flexGrow: 1}}/>
                <Profile organization={organization}/>
            </Toolbar>
        </AppBar>
    )
}
