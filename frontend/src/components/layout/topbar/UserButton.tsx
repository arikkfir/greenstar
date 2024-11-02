import {MouseEvent, useCallback, useState} from "react";
import {Avatar, Divider, IconButton, Menu, MenuItem} from "@mui/material";
import AccountCircle from "@mui/icons-material/AccountCircle";
import {useNavigate} from "react-router-dom";
import Settings from '@mui/icons-material/Settings';
import Logout from '@mui/icons-material/Logout';
import {useDescope, useSession, useUser} from "@descope/react-sdk";

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
                {initials ? <Avatar sx={{width: 32, height: 32}}>{initials}</Avatar> : <AccountCircle/>}
            </IconButton>
            <Menu id="menu-appbar"
                  anchorEl={userMenuAnchorEl}
                  keepMounted
                  open={!!userMenuAnchorEl}
                  onClose={() => setUserMenuAnchorEl(null)}
            >
                <MenuItem onClick={() => navigateTo("/user/profile")}>
                    <AccountCircle sx={{mr: 1}} fontSize="small"/>
                    My Profile
                </MenuItem>
                <MenuItem onClick={() => navigateTo("/settings")}>
                    <Settings sx={{mr: 1}} fontSize="small"/>
                    Settings
                </MenuItem>
                <Divider/>
                <MenuItem onClick={() => sdk.logout()}>
                    <Logout sx={{mr: 1}} fontSize="small"/>
                    Log out
                </MenuItem>
            </Menu>
        </>
    )
}