import {PersonRounded as PersonRoundedIcon} from "@mui/icons-material";
import {Avatar, IconButton, ListItemIcon, Menu, MenuItem} from "@mui/material";
import React from "react";
import {User} from "./User";

export interface ProfileMenuProps {
    user: User
}

const menuId = 'primary-search-account-menu';

export function ProfileMenu({user}: ProfileMenuProps) {
    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
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
            </Menu>
        </>
    )
}
