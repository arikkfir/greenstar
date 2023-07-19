import {Box, CircularProgress} from "@mui/material";

export function PageLoader() {
    return (
        <Box sx={{display: 'flex'}}>
            <CircularProgress/>
        </Box>
    );
}
