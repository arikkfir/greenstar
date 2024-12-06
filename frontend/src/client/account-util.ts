import {Account} from "./account.ts";

export type AccountID = Account['id']

export type AccountsMap = {[p: AccountID]: Account}

export function MapAccountsByID(accounts: Account[]): AccountsMap {
    const accountsMap: AccountsMap = {}
    accounts.forEach(acc => accountsMap[acc.id] = acc)
    return accountsMap
}
