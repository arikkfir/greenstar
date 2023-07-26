import {Box, CircularProgress} from "@mui/material";

export function PageLoader() {
    // TODO: improve page loader
    return (
        <Box sx={{display: 'flex'}}>
            <CircularProgress/>
        </Box>
    );
}
