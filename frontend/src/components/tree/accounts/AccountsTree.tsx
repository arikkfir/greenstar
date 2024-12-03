import {Account} from "../../../client/account.ts";
import {SyntheticEvent, useCallback, useEffect, useState} from "react";
import {TreeItemEx} from "./TreeItemEx.tsx";
import {AccountNode} from "./AccountNode.ts";
import {RichTreeViewPro} from "@mui/x-tree-view-pro";

export interface AccountsTreeProps {
    accounts: Account[]
    loading?: boolean
    selectedAccount?: Account
    onAccountSelected?: (account: Account | undefined) => void
}

export function AccountsTree({selectedAccount, accounts, onAccountSelected}: AccountsTreeProps) {
    const [rootNodes, setRootNodes] = useState<AccountNode[]>([])
    const [selectedItemID, setSelectedItemID] = useState<string | undefined>(selectedAccount ? selectedAccount.id : undefined);

    useEffect(() => {
        if (selectedAccount) {
            setSelectedItemID(selectedAccount.id)
        } else {
            setSelectedItemID(undefined)
        }
    }, [selectedAccount])

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

    const handleItemSelectionChange = useCallback((_: SyntheticEvent, id: string | null) => {
        if (onAccountSelected) {
            if (id) {
                onAccountSelected(accounts.find(n => n.id === id))
            } else {
                onAccountSelected(undefined)
            }
        }
    }, [onAccountSelected, accounts]);

    return (
        <RichTreeViewPro<AccountNode, false> getItemLabel={getItemLabel}
                                             items={rootNodes}
                                             selectedItems={selectedItemID}
                                             onSelectedItemsChange={handleItemSelectionChange}
                                             slots={{item: TreeItemEx}}/>
    )
}
