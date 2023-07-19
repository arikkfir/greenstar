import {LogoutRounded as LogoutRoundedIcon, PersonRounded as PersonRoundedIcon} from "@mui/icons-material";
import {Avatar, Divider, IconButton, ListItemIcon, Menu, MenuItem} from "@mui/material";
import {useState, MouseEvent, useCallback} from "react";
import {useDescope, useUser} from "@descope/react-sdk";

const menuId = 'primary-search-account-menu';

export function UserMenu() {
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
            <Menu id={menuId}
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
