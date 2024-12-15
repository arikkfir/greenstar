import { Account } from "./account.ts"

export type AccountID = Account["id"]

export type AccountsMap = { [p: AccountID]: Account }

export interface AccountNode extends Account {
    children: AccountNode[]
}

export function MapAccountsByID(accounts: Account[]): AccountsMap {
    const accountsMap: AccountsMap = {}
    accounts.forEach((acc) => (accountsMap[acc.id] = acc))
    return accountsMap
}

export function BuildAccountsTree(accounts: Account[]): AccountNode[] {
    const accountsMap = MapAccountsByID(accounts || [])
    const accountNodesMap: { [p: AccountID]: AccountNode } = {}
    for (let id in accountsMap) {
        accountNodesMap[id] = Object.assign({ children: [] }, accountsMap[id])
    }

    const rootNodes: AccountNode[] = []
    for (let id in accountNodesMap) {
        const account = accountNodesMap[id]
        if (account.parentID) {
            const parent = accountNodesMap[account.parentID]
            if (!parent) {
                throw new Error(`parent account '${account.parentID}' of account '${account.id}' could not be found`)
            }
            parent.children.push(account)
        } else {
            rootNodes.push(account)
        }
    }
    return rootNodes
}

export function FindAccountInTree(tree: AccountNode[], id: AccountID): AccountNode | undefined {
    for (let account of tree) {
        if (account.id === id) {
            return account
        }

        const child = FindAccountInTree(account.children, id)
        if (child) {
            return child
        }
    }
    return undefined
}
