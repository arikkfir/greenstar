import {Box, Link, Typography} from "@mui/material";
import React from "react";

export function Copyright() {
    return (
        <Box display="flex" justifyContent="center" alignItems="center">
            <Typography variant="body2" color="text.secondary" align="center">
                {'Copyright © '}
                <Link color="inherit" href="https://github.com/arikkfir">arikkfir</Link>
                {' '}
                {new Date().getFullYear()}.
            </Typography>
        </Box>
    );
}
