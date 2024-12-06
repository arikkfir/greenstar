import {Divider, IconButton, Menu, MenuItem, useColorScheme} from "@mui/material";
import {MouseEvent, useCallback, useState} from "react";
import ContrastIcon from "@mui/icons-material/Contrast";
import DarkModeIcon from "@mui/icons-material/DarkMode";
import LightModeIcon from "@mui/icons-material/LightMode";
import SettingsBrightnessIcon from "@mui/icons-material/SettingsBrightness";

export type Mode = 'light' | 'dark' | 'system';

export function ThemeButton() {
    const {mode, setMode} = useColorScheme()
    const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
    const openMenu = useCallback((event: MouseEvent<HTMLElement>) => setMenuAnchorEl(event.currentTarget), [setMenuAnchorEl]);
    const changeMode = useCallback((mode: Mode) => {
        setMenuAnchorEl(null);
        setMode(mode)
    }, [setMenuAnchorEl, setMode]);

    if (!mode) {
        return <></>
    }

    return (
        <>
            <IconButton size="large" color="inherit" onClick={openMenu}>
                <ContrastIcon/>
            </IconButton>
            <Menu id="menu-appbar"
                  anchorEl={menuAnchorEl}
                  keepMounted
                  open={!!menuAnchorEl}
                  onClose={() => setMenuAnchorEl(null)}
            >
                <MenuItem onClick={() => changeMode("dark")}>
                    <DarkModeIcon sx={{mr: 1}} fontSize="small"/>
                    Dark
                </MenuItem>
                <MenuItem onClick={() => changeMode("light")}>
                    <LightModeIcon sx={{mr: 1}} fontSize="small"/>
                    Light
                </MenuItem>
                <Divider/>
                <MenuItem onClick={() => changeMode("system")}>
                    <SettingsBrightnessIcon sx={{mr: 1}} fontSize="small"/>
                    System
                </MenuItem>
            </Menu>
        </>
    )
}
