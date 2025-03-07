import pkg from "../../package.json"
import { AppBar, Typography } from "@mui/material"

export function Footer() {
    return (
        <AppBar position="relative" component="footer">
            <Typography>
                Copyright © 2017 - {new Date().getFullYear()},
                &nbsp;
                <a href="https://github.com/arikkfir">Arik Kfir</a>,
                &nbsp;
                All Rights Reserved ({pkg.version})
            </Typography>
        </AppBar>
    )
}
