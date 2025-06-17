import { gql } from "../graphql"

export const GetTransactionsSummary = gql(`
    query GetTransactionsCount($tenantID: ID!) {
        tenant(id: $tenantID) {
            totalTransactions
        }
    }
`)

export const CreateTransaction = gql(`
    mutation CreateTransaction($tx: CreateTransaction!) {
        createTransaction(tx: $tx) {
            id
        }
    }
`)
