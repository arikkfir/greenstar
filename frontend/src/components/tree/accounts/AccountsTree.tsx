import {forwardRef, HTMLAttributes, ReactNode, Ref, useCallback, useEffect, useState} from "react";
import {
    RichTreeViewPro,
    RichTreeViewProProps,
    TreeItem2Checkbox,
    TreeItem2Content,
    TreeItem2GroupTransition,
    TreeItem2Icon,
    TreeItem2IconContainer,
    TreeItem2Label,
    TreeItem2Provider,
    TreeItem2Root,
    useTreeItem2,
    UseTreeItem2Parameters
} from "@mui/x-tree-view-pro";
import {useCurrencyFormatter} from "../../../hooks/locale.tsx";
import {Account} from "../../../client/account.ts";
import {DynamicIcon} from "../../DynamicIcon.tsx";
import {Box} from "@mui/material";

export interface AccountNode extends Account {
    children: Account[]
}

function buildAccountsTree(accounts: Account[]): AccountNode[] {
    let nodesByID: { [key: string]: AccountNode } = {}
    let nodes: AccountNode[] = accounts.map(acc => {
        const node = Object.assign(acc, {children: []})
        nodesByID[acc.id] = node
        return node
    })

    const rootNodes: AccountNode[] = []
    nodes.forEach(n => {
        if (n.parentID) {
            const parent = nodesByID[n.parentID]
            if (!parent) {
                throw new Error(`parent account '${n.parentID}' of account '${n.id}' could not be found`)
            }
            parent.children?.push(n)
        } else {
            rootNodes.push(n)
        }
    })
    return rootNodes
}

export interface AccountTreeItemProps
    extends Omit<UseTreeItem2Parameters, 'rootRef'>, Omit<HTMLAttributes<HTMLLIElement>, 'onFocus'> {
}

export const AccountTreeItem = forwardRef(function (props: AccountTreeItemProps, ref: Ref<HTMLLIElement>): ReactNode {
    const {id, itemId, label, disabled, children, ...other} = props;
    const {
        getRootProps,
        getContentProps,
        getIconContainerProps,
        getCheckboxProps,
        getLabelProps,
        getGroupTransitionProps,
        publicAPI,
        status,
    } = useTreeItem2({id, itemId, children, label, disabled, rootRef: ref});
    const item: AccountNode = publicAPI.getItem(itemId) as AccountNode
    const currencyFormatter = useCurrencyFormatter()
    return (
        <TreeItem2Provider itemId={itemId}>
            <TreeItem2Root {...getRootProps(other)} >
                <TreeItem2Content {...getContentProps()} data-testid={`accountNode:${item.id}`}>
                    <TreeItem2IconContainer {...getIconContainerProps()} data-testid={`expander:${item.id}`}>
                        <TreeItem2Icon status={status}/>
                    </TreeItem2IconContainer>
                    <TreeItem2Checkbox {...getCheckboxProps()} data-testid={`checkbox:${item.id}`} />
                    <Box sx={{flexGrow: 1, display: 'flex', gap: 1}}>
                        <DynamicIcon icon={item.icon || 'Work'} data-testid={`icon:${item.id}`}/>
                        <TreeItem2Label sx={{flexGrow: 1}} {...getLabelProps()} data-testid={`label:${item.id}`} />
                        <Box data-testid={`balance:${item.id}`} sx={{
                            flexGrow: 0,
                            flexShrink: 0
                        }}>{item.balance ? currencyFormatter(item.balance) : "-"}</Box>
                    </Box>
                </TreeItem2Content>
                {children && <TreeItem2GroupTransition {...getGroupTransitionProps()} />}
            </TreeItem2Root>
        </TreeItem2Provider>
    );
})

export interface AccountsTreeProps extends Omit<RichTreeViewProProps<Account, false>, 'getItemLabel' | 'items' | 'selectedItems' | 'slots'> {
    accounts: Account[]
    selectedAccount?: Account
}

export function AccountsTree({selectedAccount, accounts}: AccountsTreeProps) {
    const [rootNodes, setRootNodes] = useState<AccountNode[]>([])
    const [selectedItemID, setSelectedItemID] = useState<string | undefined>(selectedAccount ? selectedAccount.id : undefined);

    useEffect(() => setSelectedItemID(selectedAccount ? selectedAccount.id : undefined), [selectedAccount])
    useEffect(() => setRootNodes(buildAccountsTree(accounts)), [accounts, buildAccountsTree, setRootNodes]);

    const getItemLabel = useCallback((account: AccountNode) => account.displayName, [])

    return (
        <RichTreeViewPro<AccountNode, false> getItemLabel={getItemLabel}
                                             items={rootNodes}
                                             selectedItems={selectedItemID}
                                             slots={{item: AccountTreeItem}}/>
    )
}
