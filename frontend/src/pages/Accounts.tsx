import {TreeItem, TreeView} from "@mui/lab";
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import {Box, LinearProgress, Stack} from "@mui/material";
import {Account} from "../gql/graphql";
import {useApolloClient, useQuery} from "@apollo/client";
import {accountChildren, rootAccounts} from "../services/accounts";
import {useCallback, useMemo, useState} from "react";

type AccountNode = {
    id: Account['id'],
    displayName: Account['displayName'],
    childCount: Account['childCount'],
}

type ChildMappings = {
    [key: AccountNode['id']]: Array<AccountNode>
}

export function Accounts({tenantID}: { tenantID: string }) {
    const client = useApolloClient();

    const {
        data: roots,
        loading: rootsLoading,
        error: rootsLoadingError
    } = useQuery(rootAccounts, {variables: {tenantID}});

    const rootNodes = useMemo(
        () => (roots?.tenant?.accounts.map(({id, displayName, childCount}) => ({id, displayName, childCount})) || []),
        [roots?.tenant?.accounts],
    )

    const [childMappings, setChildMappings] = useState<ChildMappings>({})

    const renderNode = useCallback((node: AccountNode) => {
        let children = null
        if (node.childCount) {
            if (childMappings[node.id]) {
                children = childMappings[node.id]?.map((child) => renderNode(child))
            } else {
                children = "Loading..."
            }
        }
        return <TreeItem key={node.id} nodeId={node.id} label={node.displayName}>{children}</TreeItem>
    }, [childMappings]);

    if (rootsLoading) {
        return (
            <Box sx={{width: '100%'}}>
                <LinearProgress/>
            </Box>
        )
    }

    if (rootsLoadingError) {
        return (
            <Box sx={{width: '100%'}}>
                Error loading accounts!
            </Box>
        )
    }

    const fetchMissingChildren = async (nodeIds: string[]) => {
        let newChildMappings: ChildMappings = {}
        Promise.all(
            nodeIds.map(async (id) => {
                newChildMappings[id] = []
                const queryOptions = {
                    query: accountChildren,
                    variables: {tenantID, accountID: id},
                };
                const result = await client.query(queryOptions);
                const children = (result.data.tenant?.account?.children || []);
                newChildMappings[id] = children.map(({id, displayName, childCount}) => ({id, displayName, childCount}))
            })
        ).then(() => setChildMappings(newChildMappings))
    }

    return (
        <Stack direction="row" sx={{padding: 1}}>
            <TreeView defaultCollapseIcon={<ExpandMoreIcon/>}
                      defaultExpandIcon={<ChevronRightIcon/>}
                      onNodeToggle={(_, nodeIds) => fetchMissingChildren(nodeIds)}
                      sx={{height: 240, flexGrow: 0, maxWidth: 250, overflowY: 'auto'}}>
                {rootNodes.map(node => renderNode(node))}
            </TreeView>
            <Box sx={{flexGrow: 1, overflow: 'auto'}}>
                Text text text
            </Box>
        </Stack>
    )
}
