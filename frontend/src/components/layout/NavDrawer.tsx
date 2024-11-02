import {Divider, List, ListItem, ListItemButton, ListItemIcon, ListItemText} from "@mui/material";
import AutoGraphIcon from '@mui/icons-material/AutoGraph';
import {useLocation, useNavigate} from "react-router-dom";
import RequestQuoteIcon from '@mui/icons-material/RequestQuote';
import ElectricalServicesIcon from '@mui/icons-material/ElectricalServices';
import {Theme as SystemTheme} from "@mui/system/createTheme/createTheme";
import {SxProps} from "@mui/system/styleFunctionSx";

export interface NavDrawerProps<Theme extends object = SystemTheme> {
    sx?: SxProps<Theme>
}

export function NavDrawer({sx}: NavDrawerProps) {
    const navigate = useNavigate()
    const location = useLocation();
    return (
        <List sx={sx}>
            <ListItem disablePadding>
                <ListItemButton selected={location.pathname == "/"} onClick={() => navigate("/")}>
                    <ListItemIcon>
                        <AutoGraphIcon/>
                    </ListItemIcon>
                    <ListItemText primary="Dashboard"/>
                </ListItemButton>
            </ListItem>
            <ListItem disablePadding>
                <ListItemButton selected={location.pathname == "/transactions"}
                                onClick={() => navigate("/transactions")}>
                    <ListItemIcon>
                        <RequestQuoteIcon/>
                    </ListItemIcon>
                    <ListItemText primary="Transactions"/>
                </ListItemButton>
            </ListItem>
            <Divider/>
            <ListItem disablePadding>
                <ListItemButton selected={location.pathname == "/api"} onClick={() => navigate("/api")}>
                    <ListItemIcon>
                        <ElectricalServicesIcon/>
                    </ListItemIcon>
                    <ListItemText primary="API Playground"/>
                </ListItemButton>
            </ListItem>
        </List>
    )
}
