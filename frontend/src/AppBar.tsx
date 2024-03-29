import {
    AppBar as MuiAppBar,
    Avatar,
    Box,
    Divider,
    IconButton,
    ListItemIcon,
    Menu,
    MenuItem,
    Toolbar,
    Typography
} from "@mui/material";
import {useDescope, useUser} from "@descope/react-sdk";
import {LogoutRounded as LogoutRoundedIcon, PersonRounded as PersonRoundedIcon} from "@mui/icons-material";
import {MouseEvent, useCallback, useState} from "react";

function UserMenu() {
    const {user} = useUser()
    const {logout} = useDescope()
    const handleLogout = useCallback(() => logout(), [logout])

    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const openProfileMenu = (event: MouseEvent<HTMLElement>) => setAnchorEl(event.currentTarget);
    const closeMenu = useCallback(() => setAnchorEl(null), []);
    return (
        <>
            <IconButton size="large" edge="end" onClick={openProfileMenu}>
                <Avatar imgProps={{referrerPolicy: "no-referrer"}}
                        alt={user.name ? user.name[0] : "??"}
                        src={user.picture}></Avatar>
            </IconButton>
            <Menu id="primary-search-account-menu"
                  keepMounted
                  anchorEl={anchorEl}
                  anchorOrigin={{vertical: 'top', horizontal: 'right'}}
                  transformOrigin={{vertical: 'top', horizontal: 'right'}}
                  open={Boolean(anchorEl)}
                  onClose={closeMenu}>
                <MenuItem>
                    <ListItemIcon>
                        <PersonRoundedIcon fontSize="small"/>
                    </ListItemIcon>
                    Profile
                </MenuItem>
                <Divider/>
                <MenuItem onClick={handleLogout}>
                    <ListItemIcon>
                        <LogoutRoundedIcon fontSize="small"/>
                    </ListItemIcon>
                    Logout
                </MenuItem>
            </Menu>
        </>
    )
}

export function AppBar() {
    return (
        <MuiAppBar position="static" sx={{zIndex: (theme) => theme.zIndex.drawer + 1}}>
            <Toolbar>
                <Typography variant="h6" noWrap>
                    GreenSTAR
                </Typography>
                <Box sx={{flexGrow: 1}}/>
                <UserMenu/>
            </Toolbar>
        </MuiAppBar>
    )
}
