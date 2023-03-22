import {ListItemButton, ListItemIcon, ListItemText} from "@mui/material";
import React from "react";
import {Link as RouterLink, LinkProps as RouterLinkProps} from "react-router-dom";

export interface ListItemLinkProps {
    icon: React.ReactElement;
    primary: string;
    to: string;
}

const ClientLink =
    React.forwardRef<HTMLAnchorElement, RouterLinkProps>(
        function Link(itemProps, ref) {
            return <RouterLink ref={ref} {...itemProps} role={undefined}/>;
        },
    );

export function ListItemLink({icon, primary, to}: ListItemLinkProps) {
    return (
        <ListItemButton component={ClientLink} to={to}>
            <ListItemIcon>{icon}</ListItemIcon>
            <ListItemText primary={primary}/>
        </ListItemButton>
    );
}
