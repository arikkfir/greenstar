import {LineAxis, ReceiptLong} from "@mui/icons-material";
import {Divider, Drawer, List, ListItem, ListItemButton, ListItemIcon, ListItemText} from "@mui/material";
import {RequiresTenantPermission} from "./components/RequiresTenantPermissions.tsx";
import {RequiresFeatureFlag} from "./components/RequiresFeatureFlag.tsx";
import {RequiresGlobalTenant} from "./components/RequiresGlobalTenant.tsx";

interface AppDrawerProps {
    tenant: string
}

export function AppDrawer({tenant}: AppDrawerProps) {
    return (
        <Drawer variant="permanent" PaperProps={{sx: {position: "static"}}}>
            <List>
                <ListItem disablePadding>
                    <ListItemButton href="/">
                        <ListItemIcon><LineAxis/></ListItemIcon>
                        <ListItemText primary="Dashboard"/>
                    </ListItemButton>
                </ListItem>
                <RequiresTenantPermission tenant={tenant} anyOf={["accounts:read"]}>
                    <ListItem disablePadding>
                        <ListItemButton href="/accounts">
                            <ListItemIcon><LineAxis/></ListItemIcon>
                            <ListItemText primary="Accounts"/>
                        </ListItemButton>
                    </ListItem>
                </RequiresTenantPermission>
                <RequiresTenantPermission tenant={tenant} anyOf={["transactions:read"]}>
                    <ListItem disablePadding>
                        <ListItemButton href="/transactions">
                            <ListItemIcon><LineAxis/></ListItemIcon>
                            <ListItemText primary="Transactions"/>
                        </ListItemButton>
                    </ListItem>
                </RequiresTenantPermission>
                <RequiresGlobalTenant tenant={tenant}>
                    <Divider/>
                    <ListItem disablePadding>
                        <ListItemButton href="/tenants">
                            <ListItemIcon><LineAxis/></ListItemIcon>
                            <ListItemText primary="Tenants"/>
                        </ListItemButton>
                    </ListItem>
                </RequiresGlobalTenant>
                <RequiresFeatureFlag mustBeOn={["show-graphql-console"]}>
                    <Divider/>
                    <ListItem disablePadding>
                        <ListItemButton href="/api">
                            <ListItemIcon><ReceiptLong/></ListItemIcon>
                            <ListItemText primary="API Explorer"/>
                        </ListItemButton>
                    </ListItem>
                </RequiresFeatureFlag>
            </List>
        </Drawer>
    )
}
