import MenuIcon from "@mui/icons-material/Menu";
import {AppBar, Box, IconButton, Toolbar, Typography} from "@mui/material";
import React from "react";
import {ProfileMenu} from "./ProfileMenu";
import {User} from "./User";

export interface AppBarProps {
    user: User | undefined
}

export function TopAppBar({user}: AppBarProps) {
    return (
        <AppBar position="static" sx={{zIndex: (theme) => theme.zIndex.drawer + 1}}>
            <Toolbar>
                <IconButton size="large" edge="start" color="inherit" aria-label="open drawer" sx={{mr: 2}}>
                    <MenuIcon/>
                </IconButton>
                <Typography variant="h6" noWrap>Greenstar</Typography>
                <Box sx={{flexGrow: 1}}/>
                {user && [<ProfileMenu key="profile-menu" user={user}/>]}
            </Toolbar>
        </AppBar>
    )
}
