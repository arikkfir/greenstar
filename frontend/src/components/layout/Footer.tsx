import {AppBar, Typography} from "@mui/material";
import pkg from "../../../package.json";
import {Link} from "react-router";
import {Theme as SystemTheme} from "@mui/system/createTheme/createTheme";
import {SxProps} from "@mui/system/styleFunctionSx";

export interface FooterProperties<Theme extends object = SystemTheme> {
    sx: SxProps<Theme>
}

export function Footer({sx}: FooterProperties) {
    return (
        <AppBar component="footer"
                position="relative"
                sx={{
                    boxShadow: '0 -2px 4px rgba(0,0,0,0.4)',
                    p: 0.5,
                    display: "flex",
                    flexDirection: "row",
                    justifyContent: "center",
                    alignContent: "center",
                    ...sx
                }}>
            <Typography variant="body1">
                Copyright © 2017 - {new Date().getFullYear()}&nbsp;&nbsp;
            </Typography>
            <Link to="https://github.com/arikkfir">Arik Kfir</Link>&nbsp;&nbsp;
            <Typography variant="body1">
                All Rights Reserved ({pkg.version})
            </Typography>
        </AppBar>
    )
}
