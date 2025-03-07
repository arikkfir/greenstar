import { DateTime } from "luxon"

export interface TransactionAccountInfo {
    id: string
    icon: string
    displayName: string
}

export interface TransactionRow {
    id: string
    createdAt: DateTime
    updatedAt: DateTime
    currency: string
    amount: number
    date: DateTime
    description?: string
    referenceID: string
    sourceAccount: TransactionAccountInfo
    targetAccount: TransactionAccountInfo
}
