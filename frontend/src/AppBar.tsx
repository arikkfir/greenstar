import {AccountCircle} from "@mui/icons-material";
import MenuIcon from "@mui/icons-material/Menu";
import {AppBar, Box, IconButton, Menu, MenuItem, Toolbar, Typography} from "@mui/material";
import {useAtom} from "jotai/react";
import React from "react";
import state, {User} from './state';

export default function TopAppBar() {
    const menuId = 'primary-search-account-menu';
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
    const isMenuOpen = Boolean(anchorEl);
    const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => setAnchorEl(event.currentTarget);
    const handleMenuClose = () => setAnchorEl(null);
    const [user] = useAtom<User>(state.user)
    return (
        <AppBar position="static">
            <Toolbar>
                <IconButton size="large" edge="start" color="inherit" aria-label="open drawer" sx={{mr: 2}}>
                    <MenuIcon/>
                </IconButton>
                <Typography variant="h6" noWrap>Greenstar</Typography>
                <Box sx={{flexGrow: 1}}/>
                <IconButton size="large" edge="end" aria-label="account of current user" aria-controls={menuId}
                            aria-haspopup="true" onClick={handleProfileMenuOpen} color="inherit">
                    <AccountCircle/>
                </IconButton>
                <Menu
                    anchorEl={anchorEl}
                    anchorOrigin={{vertical: 'top', horizontal: 'right'}}
                    id={menuId}
                    keepMounted
                    transformOrigin={{vertical: 'top', horizontal: 'right'}}
                    open={isMenuOpen}
                    onClose={handleMenuClose}
                >
                    <MenuItem onClick={handleMenuClose}>Profile of {user.name}</MenuItem>
                    <MenuItem onClick={handleMenuClose}>My account</MenuItem>
                </Menu>
            </Toolbar>
        </AppBar>
    )
}
