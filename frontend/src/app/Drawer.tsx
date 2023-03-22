import {Info, LineAxis, ReceiptLong} from "@mui/icons-material";
import {Divider, Drawer, List, ListItem, ListItemButton, ListItemIcon, ListItemText} from "@mui/material";
import React from "react";
import {ListItemLink} from "./common/ListItemLink";

export interface AppDrawerProps {
    adminAPIURL: string
    operationsAPIURL: string
    publicAPIURL: string
}

export function AppDrawer({adminAPIURL, operationsAPIURL, publicAPIURL}: AppDrawerProps) {
    return (
        <Drawer variant="permanent" PaperProps={{sx: {position: "static"}}}>
            <List>
                <ListItem disablePadding>
                    <ListItemLink icon={<LineAxis/>} primary="Dashboard" to="/"/>
                </ListItem>
                <ListItem disablePadding>
                    <ListItemLink icon={<ReceiptLong/>} primary="Transactions" to="/transactions"/>
                </ListItem>
                <Divider/>
                <ListItem disablePadding>
                    <ListItemButton component="a" href={adminAPIURL}>
                        <ListItemIcon>{<ReceiptLong/>}</ListItemIcon>
                        <ListItemText primary="Admin API"/>
                    </ListItemButton>
                </ListItem>
                <ListItem disablePadding>
                    <ListItemButton component="a" href={operationsAPIURL}>
                        <ListItemIcon>{<ReceiptLong/>}</ListItemIcon>
                        <ListItemText primary="Operations API"/>
                    </ListItemButton>
                </ListItem>
                <ListItem disablePadding>
                    <ListItemButton component="a" href={publicAPIURL}>
                        <ListItemIcon>{<ReceiptLong/>}</ListItemIcon>
                        <ListItemText primary="Public API"/>
                    </ListItemButton>
                </ListItem>
                <Divider/>
                <ListItem disablePadding>
                    <ListItemLink icon={<Info/>} primary="About" to="/about"/>
                </ListItem>
            </List>
        </Drawer>
    )
}
