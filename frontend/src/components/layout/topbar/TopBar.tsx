import { AppBar, Box, Toolbar, Typography } from "@mui/material"
import { SxProps } from "@mui/system/styleFunctionSx"
import { Theme as SystemTheme } from "@mui/system/createTheme/createTheme"
import { ThemeButton } from "./ThemeButton.tsx"
import { LinkButton } from "../../LinkButton.tsx"

export interface TopBarProperties<Theme extends object = SystemTheme> {
    sx: SxProps<Theme>
}

export function TopBar({ sx }: TopBarProperties) {
    return (
        <AppBar component="header" position="relative" sx={sx}>
            <Toolbar>
                <Typography variant="h5" component="div" sx={{ fontWeight: "700", mr: 5 }}>
                    GreenSTAR
                </Typography>
                <Box sx={{ display: "flex", flexGrow: 1 }}>
                    <LinkButton to="/">Dashboard</LinkButton>
                    <LinkButton to="/transactions">Transactions</LinkButton>
                    <LinkButton to="/api">API</LinkButton>
                </Box>
                <ThemeButton />
            </Toolbar>
        </AppBar>
    )
}
