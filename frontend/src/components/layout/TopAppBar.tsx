import MenuIcon from "@mui/icons-material/Menu";
import {AppBar, Box, IconButton, Toolbar, Typography} from "@mui/material";
import {UserMenu} from "./UserMenu";

export function TopAppBar() {
    return (
        <AppBar position="static" sx={{zIndex: (theme) => theme.zIndex.drawer + 1}}>
            <Toolbar>
                <IconButton size="large" edge="start" color="inherit" aria-label="open drawer" sx={{mr: 2}}>
                    <MenuIcon/>
                </IconButton>
                <Typography variant="h6" noWrap>
                    GreenSTAR
                </Typography>
                <Box sx={{flexGrow: 1}}/>
                <UserMenu/>
            </Toolbar>
        </AppBar>
    )
}
