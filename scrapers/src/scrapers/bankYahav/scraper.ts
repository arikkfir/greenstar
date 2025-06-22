/**
 * Bank Yahav Scraper
 *
 * This module implements a Playwright test that scrapes transaction data from Bank Yahav's website.
 * It logs in to the bank's website, navigates to the transactions page, and iterates through
 * transaction history, downloading and processing transactions for storage in the system.
 *
 * The scraper handles both regular account transactions and debit card transactions,
 * and can initialize account balance if needed.
 */

import { test } from "@playwright/test"
import { Site } from "./site.ts"
import { client } from "../../graphql/client.ts"
import { Config } from "./configuration.ts"
import { AccountTransaction } from "./account-transaction.ts"
import { AccountTransactionsPage } from "./account-transactions-page.ts"
import { BalanceTransaction } from "./balance-transaction.ts"
import { ITransaction, saveTransaction } from "./transaction.ts"
import { gql } from "../../graphql"

/**
 * Query to get a summary of transactions for a tenant
 *
 * Retrieves the total number of transactions for the specified tenant.
 *
 * @param {ID!} tenantID - The ID of the tenant to query
 * @returns {Object} An object containing the total number of transactions
 */
export const GetTransactionsSummary = gql(`
    query GetTransactionsCount($tenantID: ID!) {
        tenant(id: $tenantID) {
            totalTransactions
        }
    }
`)

/**
 * Main test function for scraping Bank Yahav transactions
 *
 * This test performs the following operations:
 * 1. Retrieves configuration and authentication details
 * 2. Gets current transaction summary from the GraphQL API
 * 3. Logs into the Bank Yahav website
 * 4. Navigates to and scrapes transaction data, moving backward through history
 * 5. Creates initialization transaction if needed
 * 6. Processes both regular and debit card transactions
 * 7. Saves all transactions to the system
 *
 * @param {Object} param0 - Playwright test context
 * @param {Page} param0.page - Playwright page object
 */
test("bank-yahav", async ({ page }) => {
    const config    = new Config()
    const tenantID  = config.tenantID
    const accountID = config.accountID
    const username  = config.username
    const password  = config.password
    const pinno     = config.pinno

    const result = await client.query(GetTransactionsSummary, { tenantID })
    if (result.error) {
        throw new Error(`Failed to get transactions summary for tenant ${tenantID}: ${result.error.message}`)
    }
    const totalTransactions: number = result.data!.tenant!.totalTransactions

    // Open the site and login
    const site = new Site(page, tenantID, accountID, username, password, pinno)
    await site.open()
    await site.login()

    // Navigate to the transactions page
    const txPage: AccountTransactionsPage = await site.openTransactionsPage()

    // Move to requested date
    await txPage.setDateRange(config.startDate, config.endDate)

    // Export current account transactions page as an Excel (XLS) file
    if (config.downloadXLS) {
        await txPage.downloadTransactionsExcel()
    }

    const transactions: ITransaction[] = []

    // Scrape account-level transactions from the current date range
    const accountTransactions = await txPage.getTransactions()

    // If found at least one account-level transaction, and the tenant currently has no transactions, then this is
    // essentially the first transaction ever of this tenant; thus, given that account-level transactions also provide
    // the resulting account balance, we can infer what was the account balance before that first transaction; therefore
    // we will add a balance initialization transaction first.
    if (accountTransactions.length > 0 && totalTransactions == 0) {
        const firstTx: AccountTransaction = accountTransactions[0]
        const balanceBeforeTx: number     = await firstTx.getBalanceBeforeTransaction()
        const balanceTx                   = new BalanceTransaction(
            (await firstTx.getDate()).minus({ days: 7 }),
            "initialization",
            "Balance initialization",
            balanceBeforeTx < 0 ? Math.abs(balanceBeforeTx) : null,
            balanceBeforeTx > 0 ? Math.abs(balanceBeforeTx) : null,
        )
        transactions.push(balanceTx)
    }

    // Iterate over all account-level transactions and add each one:
    //  - If it's a debit-card aggregate transaction, navigate to the page showing its underlying transactions behind it
    //    and add each such underlying transaction to the final list of transactions
    //  - Otherwise, add the account-level transaction to the final list of transactions
    for (let accountTx of accountTransactions) {
        if (await accountTx.isDebitCardTransaction()) {
            console.info(`Processing debit card aggregate transaction: ${await accountTx.toDebugString()}`)
            const debitCardTxPage = await accountTx.navigateToDebitCardTransactions()
            if (config.downloadXLS) {
                await debitCardTxPage.downloadTransactionsExcel()
            }
            transactions.push(...await debitCardTxPage.getTransactions())
            await debitCardTxPage.close()
            await txPage.awaitUntilVisible()
        } else {
            console.info(`Processing account transaction: ${await accountTx.toDebugString()}`)
            transactions.push(accountTx)
        }
    }

    // Sort transactions based on date
    const transactionsWithDates = await Promise.all(
        transactions.map(async (tx) => ({ tx, date: await tx.getDate() })),
    )
    const sortedTransactions    = transactionsWithDates
        .sort((a, b) => a.date.valueOf() - b.date.valueOf())
        .map(item => item.tx)

    // Save the resulting full list of transactions that contains:
    // - Optional balance initialization transaction
    // - Account-level transactions
    // - Debit-card underlying transactions
    for (let i = 0; i < sortedTransactions.length; i++) {
        await saveTransaction(tenantID, accountID, sortedTransactions[i], i)
    }

    // Logout
    await site.logout()
})
