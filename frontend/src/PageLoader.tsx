import {Box, CircularProgress} from "@mui/material";

export function PageLoader() {
    // TODO: use backdrop for page loader: https://mui.com/material-ui/react-backdrop/
    return (
        <Box sx={{display: 'flex'}}>
            <CircularProgress/>
        </Box>
    );
}
