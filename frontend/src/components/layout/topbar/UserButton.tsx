import {useNavigate} from "react-router";
import {useDescope, useSession, useUser} from "@descope/react-sdk";
import {MouseEvent, useCallback, useState} from "react";
import {Avatar, Divider, IconButton, Menu, MenuItem} from "@mui/material";
import AccountCircleIcon from "@mui/icons-material/AccountCircle";
import SettingsIcon from "@mui/icons-material/Settings";
import LogoutIcon from "@mui/icons-material/Logout";

export function UserButton() {
    const navigate = useNavigate()
    const sdk = useDescope();
    const {isAuthenticated} = useSession()
    const {user} = useUser()
    const [userMenuAnchorEl, setUserMenuAnchorEl] = useState<null | HTMLElement>(null);
    const openUserMenu = useCallback((event: MouseEvent<HTMLElement>) => {
        setUserMenuAnchorEl(event.currentTarget);
    }, [setUserMenuAnchorEl]);
    const navigateTo = useCallback((path: string) => {
        setUserMenuAnchorEl(null);
        navigate(path)
    }, [setUserMenuAnchorEl, navigate]);

    if (!isAuthenticated) {
        return <></>
    }

    const givenName = user?.givenName ? user.givenName[0] : ""
    const familyName = user?.familyName ? user.familyName[0] : ""
    const initials = givenName ? givenName[0] : "" + familyName ? familyName[0] : ""

    return (
        <>
            <IconButton size="large" color="inherit" onClick={openUserMenu}>
                {initials ? <Avatar sx={{width: 32, height: 32}}>{initials}</Avatar> : <AccountCircleIcon/>}
            </IconButton>
            <Menu id="menu-appbar"
                  anchorEl={userMenuAnchorEl}
                  keepMounted
                  open={!!userMenuAnchorEl}
                  onClose={() => setUserMenuAnchorEl(null)}
            >
                <MenuItem onClick={() => navigateTo("/profile")}>
                    <AccountCircleIcon sx={{mr: 1}} fontSize="small"/>
                    My Profile
                </MenuItem>
                <MenuItem onClick={() => navigateTo("/settings")}>
                    <SettingsIcon sx={{mr: 1}} fontSize="small"/>
                    Settings
                </MenuItem>
                <Divider/>
                <MenuItem onClick={() => sdk.logout()}>
                    <LogoutIcon sx={{mr: 1}} fontSize="small"/>
                    Log out
                </MenuItem>
            </Menu>
        </>
    )
}