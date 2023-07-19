import {LineAxis, ReceiptLong} from "@mui/icons-material";
import {Divider, Drawer, List, ListItem, ListItemButton, ListItemIcon, ListItemText} from "@mui/material";
import {globalTenant, gqlPlaygroundURL} from "../../config";
import {getJwtPermissions, getSessionToken} from "@descope/react-sdk";
import {useFeatureIsOn} from "@growthbook/growthbook-react";

interface AppDrawerProps {
    tenant: string
}

export function AppDrawer({tenant}: AppDrawerProps) {
    const globalPermissions = getJwtPermissions(getSessionToken(), "global");
    const tenantPermissions = getJwtPermissions(getSessionToken(), tenant);
    const showGraphQLQueryLink = useFeatureIsOn("show-graphql-console");

    return (
        <Drawer variant="permanent" PaperProps={{sx: {position: "static"}}}>
            <List>
                <ListItem disablePadding>
                    <ListItemButton href="/">
                        <ListItemIcon><LineAxis/></ListItemIcon>
                        <ListItemText primary="Dashboard"/>
                    </ListItemButton>
                </ListItem>
                {(tenantPermissions.includes("Read accounts") || tenantPermissions.includes("Manage accounts")) && (
                    <>
                        <ListItem disablePadding>
                            <ListItemButton href="/accounts">
                                <ListItemIcon><LineAxis/></ListItemIcon>
                                <ListItemText primary="Accounts"/>
                            </ListItemButton>
                        </ListItem>
                    </>
                )}
                {(tenantPermissions.includes("Read transactions") || tenantPermissions.includes("Manage transactions")) && (
                    <>
                        <ListItem disablePadding>
                            <ListItemButton href="/transactions">
                                <ListItemIcon><LineAxis/></ListItemIcon>
                                <ListItemText primary="Transactions"/>
                            </ListItemButton>
                        </ListItem>
                    </>
                )}
                <Divider/>
                {tenant === globalTenant && globalPermissions.includes("Manage tenants") && (
                    <>
                        <ListItem disablePadding>
                            <ListItemButton href="/tenants">
                                <ListItemIcon><LineAxis/></ListItemIcon>
                                <ListItemText primary="Tenants"/>
                            </ListItemButton>
                        </ListItem>
                    </>
                )}
                {showGraphQLQueryLink && (
                    <>
                        <Divider/>
                        <ListItem disablePadding>
                            <ListItemButton href={gqlPlaygroundURL}>
                                <ListItemIcon><ReceiptLong/></ListItemIcon>
                                <ListItemText primary="GraphQL Playground"/>
                            </ListItemButton>
                        </ListItem>
                    </>
                )}
            </List>
        </Drawer>
    )
}
