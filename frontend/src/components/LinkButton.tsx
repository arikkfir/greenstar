import {Button} from "@mui/material";
import {useLocation, useNavigate} from "react-router";

interface LinkButtonProperties {
    to: string,
    children: any
}

export function LinkButton({to, children}: LinkButtonProperties) {
    const navigate = useNavigate()
    const location = useLocation();
    return (
        <Button onClick={() => navigate(to)}
                size="medium"
                variant="contained"
                sx={{
                    boxShadow: 'none',
                    border: location.pathname === to ? 1 : undefined,
                    display: 'block',
                    fontWeight: '700',
                    pl: 2, pr: 2,
                    ml: 2, mr: 2,
                    '&:hover': {boxShadow: 'none'},
                    '&:active': {boxShadow: 'none'},
                    '&:focus': {boxShadow: 'none'},
                }}
        >
            {children}
        </Button>
    )
}
