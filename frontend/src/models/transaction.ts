// This file is generated by greenstar scripts. DO NOT EDIT.
// noinspection DuplicatedCode

export interface Transaction {
    id: string
    createdAt: Date
    updatedAt: Date
    amount: number
    convertedAmount?: number
    currency: string
    date: Date
    description?: string
    referenceId: string
    sourceAccountId: string
    targetAccountId: string
}

export const DateProperties = ["createdAt", "updatedAt", "date"]