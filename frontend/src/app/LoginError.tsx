import {Stack, Typography} from "@mui/material";
import React from "react";

export interface LoginErrorProps {
    error: Error
}

export function LoginError({error}: LoginErrorProps) {
    return (
        <Stack>
            <Typography>Login failed: {error.message}</Typography>
        </Stack>
    );
}
