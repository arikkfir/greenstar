import {useAuth0} from "@auth0/auth0-react";
import {
    HourglassFullRounded as HourglassFullRoundedIcon,
    LogoutRounded as LogoutRoundedIcon,
    PersonRounded as PersonRoundedIcon
} from "@mui/icons-material";
import {Avatar, Divider, IconButton, ListItemIcon, Menu, MenuItem} from "@mui/material";
import React from "react";
import {LoginButton} from "./LoginButton";
import {Organization} from "./Organization";

const menuId = 'primary-search-account-menu';

export interface ProfileProps {
    organization: Organization
}

export function Profile({organization}: ProfileProps) {
    const {user, isAuthenticated, isLoading, logout} = useAuth0();
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

    if (isLoading) {
        return (
            <IconButton size="large" edge="end">
                <HourglassFullRoundedIcon/>
            </IconButton>
        )
    }

    if (!isAuthenticated || !user) {
        return (
            <LoginButton organizationId={organization.id}/>
        )
    }

    const openProfileMenu = (event: React.MouseEvent<HTMLElement>) => setAnchorEl(event.currentTarget);
    const closeMenu = () => setAnchorEl(null);
    return (
        <>
            <IconButton size="large" edge="end" onClick={openProfileMenu}>
                <Avatar imgProps={{referrerPolicy: "no-referrer"}}
                        alt={(user?.given_name || '?')[0] + (user?.family_name || '?')[0]}
                        src={user.picture}></Avatar>
            </IconButton>
            <Menu id={menuId}
                  keepMounted
                  anchorEl={anchorEl}
                  anchorOrigin={{vertical: 'top', horizontal: 'right'}}
                  transformOrigin={{vertical: 'top', horizontal: 'right'}}
                  open={Boolean(anchorEl)}
                  onClose={closeMenu}>
                <MenuItem onClick={closeMenu}>
                    <ListItemIcon>
                        <PersonRoundedIcon fontSize="small"/>
                    </ListItemIcon>
                    Account
                </MenuItem>
                <Divider/>
                <MenuItem onClick={() => logout()}>
                    <ListItemIcon>
                        <LogoutRoundedIcon fontSize="small"/>
                    </ListItemIcon>
                    Logout
                </MenuItem>
            </Menu>
        </>
    )
}
