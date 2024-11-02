import {AppBar, Toolbar, Typography} from "@mui/material";
import {UserButton} from "./UserButton.tsx";
import {ThemeButton} from "./Theme.tsx";

export function TopBar() {
    return (
        <AppBar position="relative">
            <Toolbar>
                <Typography variant="h6" component="div" sx={{flexGrow: 1}}>
                    GreenSTAR
                </Typography>
                <UserButton/>
                <ThemeButton/>
            </Toolbar>
        </AppBar>
    )
}
