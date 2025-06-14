import { test } from "@playwright/test"
import { Site } from "./site.ts"
import { requireEnvVar } from "../util/env.ts"
import { graphQLClient } from "../util/graphql-client.ts"
import { GetTransactionsSummary } from "./query.ts"
import { createInitializationTransaction, TransactionRow } from "./transaction_row.ts"

test("bank-yahav", async ({ page }) => {
    const tenantID  = requireEnvVar("TENANT_ID")
    const accountID = requireEnvVar("ACCOUNT_ID")
    const username  = requireEnvVar("BANK_YAHAV_USERNAME")
    const password  = requireEnvVar("BANK_YAHAV_PASSWORD")
    const pinno     = requireEnvVar("BANK_YAHAV_PINNO")

    const result = await graphQLClient.query(GetTransactionsSummary, { tenantID })
    if (result.error) {
        throw new Error(`Failed to get transactions summary for tenant ${tenantID}: ${result.error.message}`)
    }

    const totalTransactions: number = result.data!.tenant!.totalTransactions

    const site = new Site(page, tenantID, accountID, username, password, pinno)
    await site.open()
    await site.login()
    await site.openTransactionsPage()

    // Go back as much as possible in the start-date picker and select the first eligible day in that month
    const fromDatePicker = site.createFromDatePicker()
    await fromDatePicker.open()
    while (await fromDatePicker.hasPrevMonth()) {
        await fromDatePicker.navigateToPreviousMonth()
    }
    await fromDatePicker.selectFirstDayOfMonth()

    // Go back as much as possible in the to-date picker and select the last eligible day in that month
    const toDatePicker = site.createToDatePicker()
    await toDatePicker.open()
    while (await toDatePicker.hasPrevMonth()) {
        await toDatePicker.navigateToPreviousMonth()
    }
    await toDatePicker.selectLastDayOfMonth()

    let rowCount = 0
    do {
        // TODO: consider replacing this timeout with: await page.waitForLoadState("networkidle")
        await page.waitForTimeout(1000 * 3)

        // STEP 5: SCRAPE DATA
        const count = await site.getTransactionRowCount()
        let row: TransactionRow | null = null
        for (let rowIndex = 0; rowIndex < count; rowIndex++) {
            // STEP 8: SEND TRANSACTION DATA TO GREENSTAR API
            row = await site.getTransactionRow(++rowCount, rowIndex);
            if (row) {
                await row.create()
            } else {
                throw new Error(`Failed to obtain transaction row ${rowCount} / ${count}`)
            }
        }

        // If we just finished processing the first page of transactions, AND
        // This tenant has NO transactions (this is the first scrape ever), THEN
        // We need an initialization transaction to bring the balance to be up-to-date
        if (row && rowCount == count && totalTransactions === 0) {

            // Obtain the last transaction row in the page (which is essentially the first transaction we've scraped)
            await createInitializationTransaction(row)

        }

        // STEP 6: PROGRESS 1 MONTH IN "TO" DATE (before "from" date!)
        await toDatePicker.open()
        if (await toDatePicker.hasNextMonth()) {

            await toDatePicker.navigateToNextMonth()
            await toDatePicker.selectLastDayOfMonth()

        } else {
            await toDatePicker.close()
            break
        }

        // STEP 7: PROGRESS 1 MONTH IN "FROM" DATE
        await fromDatePicker.open()
        await fromDatePicker.navigateToNextMonth()
        await fromDatePicker.selectFirstDayOfMonth()

    } while (true)

    // Logout
    await site.logout()
})
