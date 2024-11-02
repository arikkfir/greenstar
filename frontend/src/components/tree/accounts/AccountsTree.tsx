import {Account} from "../../../client/account.ts";
import {SyntheticEvent, useCallback, useEffect, useState} from "react";
import {TreeItemEx} from "./TreeItemEx.tsx";
import {AccountNode} from "./AccountNode.ts";
import {RichTreeViewPro} from "@mui/x-tree-view-pro";

export interface AccountsTreeProps {
    accounts: Account[]
    loading?: boolean
    onAccountSelected?: (account: Account | undefined) => void
}

export function AccountsTree({accounts, onAccountSelected}: AccountsTreeProps) {
    const [rootNodes, setRootNodes] = useState<AccountNode[]>([])

    useEffect(() => {
        let nodesByID: { [key: string]: AccountNode } = {}
        let newNodes: AccountNode[] = accounts.map(acc => {
            const node = Object.assign(acc, {
                children: [],
            })
            nodesByID[acc.id] = node
            return node
        })

        const rootNodes: AccountNode[] = []
        newNodes.forEach(n => {
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
        setRootNodes(rootNodes)
    }, [accounts, setRootNodes]);

    const getItemLabel = useCallback((account: AccountNode) => account.displayName, [])

    const handleItemSelectionToggle = (_: SyntheticEvent, id: string, isSelected: boolean) => {
        if (onAccountSelected) {
            if (isSelected) {
                onAccountSelected(accounts.find(n => n.id === id))
            } else {
                onAccountSelected(undefined)
            }
        }
    };

    return (
        <RichTreeViewPro<AccountNode, false> getItemLabel={getItemLabel}
                                             items={rootNodes}
                                             onItemSelectionToggle={handleItemSelectionToggle}
                                             slots={{item: TreeItemEx}}/>
    )
}
