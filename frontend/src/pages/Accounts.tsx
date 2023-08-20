import {Account} from "../gql/graphql.ts";
import {useApolloClient, useQuery} from "@apollo/client";
import {graphql} from "../gql";
import {useSnackbar} from "notistack";
import {useCallback, useMemo, useState} from "react";
import AccountBalanceWalletIcon from '@mui/icons-material/AccountBalanceWallet';
import {Box, CircularProgress, Icon, LinearProgress, Paper, Stack, Typography} from "@mui/material";
import {TreeItem, TreeView} from "@mui/lab";
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';

export const rootAccounts = graphql(/* GraphQL */ `
    query rootAccounts ($tenantID: ID!) {
        tenant(id: $tenantID) {
            id
            accounts {
                id
                displayName
                childCount
                icon
            }
        }
    }
`)

export const accountChildren = graphql(/* GraphQL */ `
    query accountChildren ($tenantID: ID!, $accountID: ID!) {
        tenant(id: $tenantID) {
            id
            account(id: $accountID) {
                id
                children {
                    id
                    displayName
                    childCount
                    icon
                }
            }
        }
    }
`)

type AccountNode = {
    id: Account['id'],
    displayName: Account['displayName'],
    childCount: Account['childCount'],
    icon: string
}

type ChildMappings = {
    [key: AccountNode['id']]: Array<AccountNode>
}

interface AccountsProps {
    tenantID: string
}

export function Accounts({tenantID}: AccountsProps) {
    const {enqueueSnackbar} = useSnackbar();
    const client = useApolloClient();
    const {data: roots, loading: rootsLoading, error: rootsLoadingError} = useQuery(rootAccounts, {variables: {tenantID}});
    const rootNodes: AccountNode[] = useMemo(
        () => (roots?.tenant?.accounts || []),
        [roots?.tenant?.accounts],
    )
    const [childMappings, setChildMappings] = useState<ChildMappings>({})
    const [loadingNodes, setLoadingNodes] = useState<Array<AccountNode['id']>>([])

    const renderNode = useCallback((node: AccountNode) => {
        let children = null
        if (node.childCount) {
            if (childMappings[node.id]) {
                children = childMappings[node.id]?.map((child) => renderNode(child))
            } else {
                children = [<Box key={-1 * Math.random()}/>]
            }
        }
        let icon = <AccountBalanceWalletIcon/>
        if (node.icon) {
            icon = <Icon>{node.icon}</Icon>
            console.warn(`TODO: use icon '${node.icon}' for account '${node.displayName}'`)
            // TODO: use 'icon' instead of AccountBalanceWalletIcon in the layout below
        }
        return (
            <TreeItem key={node.id}
                      nodeId={node.id}
                      label={
                          <Box sx={{alignItems: "center", display: 'flex'}}>
                              <Box sx={{mr: 0.5}}>
                                  {loadingNodes.includes(node.id) && <CircularProgress size={24}/>}
                                  {!loadingNodes.includes(node.id) && icon}
                              </Box>
                              <Typography variant="body2" sx={{fontWeight: 'inherit', flexGrow: 1, mr: 2}}>
                                  {node.displayName}
                              </Typography>
                              <Typography variant="caption" color="inherit">
                                  {"labelInfo"}
                              </Typography>
                          </Box>
                      }>
                {children}
            </TreeItem>
        )
    }, [childMappings, loadingNodes]);

    if (rootsLoading) {
        return (
            <Box sx={{width: '100%'}}>
                <LinearProgress/>
            </Box>
        )
    } else if (rootsLoadingError) {
        return (
            <Box sx={{width: '100%'}}>
                Error loading accounts!
            </Box>
        )
    }

    const fetchMissingChildren = async (nodeIds: string[]) => {
        let newChildMappings: ChildMappings = {}
        const loadPromises = nodeIds.map(async (id) => {
            if (!childMappings[id] && !loadingNodes.includes(id)) {
                setLoadingNodes([...loadingNodes, id])
            }
            const result = await client.query({
                query: accountChildren,
                variables: {tenantID, accountID: id},
            });
            newChildMappings[id] = result.data.tenant?.account?.children || []
        });
        loadPromises.push(new Promise<void>(res => setTimeout(res, 1000)));
        Promise.all(loadPromises)
            .then(() => setChildMappings(newChildMappings))
            .then(() => setLoadingNodes(loadingNodes.filter(id => !nodeIds.includes(id))))
            .catch(err => enqueueSnackbar(err.message, {variant: 'error'}))
    }

    return (
        <Stack direction="row" spacing={2} sx={{width: '100%', height: "100%"}}>
            <Paper elevation={3} sx={{padding: 2, minWidth: '25%', height: "100%"}}>
                <TreeView defaultCollapseIcon={<ExpandMoreIcon/>}
                          defaultExpandIcon={<ChevronRightIcon/>}
                          onNodeToggle={(_, nodeIds) => fetchMissingChildren(nodeIds)}
                          sx={{height: "100%"}}>
                    {rootNodes.map(node => renderNode(node))}
                </TreeView>
            </Paper>
            <Paper elevation={3} sx={{padding: 3, flexGrow: 1, overflow: 'auto'}}>
                Text text text
            </Paper>
        </Stack>
    )
}
