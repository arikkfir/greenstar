import { expect, test } from "@playwright/test"
import { Site } from "./site.ts"
import { requireEnvVar } from "../util/env.ts"
import { graphQLClient } from "../util/graphql-client.ts"
import { GetTransactionsSummary } from "./query.ts"
import { createInitializationTransaction, TransactionRow } from "./transaction_row.ts"

function buildTransactionsFileName(fromDate: string, toDate: string) {
    const fromTokens = fromDate.split("/")
    const toTokens   = toDate.split("/")
    return `${fromTokens[2]}${fromTokens[1]}${fromTokens[0]}-${toTokens[2]}${toTokens[1]}${toTokens[0]}`
}

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

    // Setup file download handlers
    const xlsDownloadIconLocator = page.locator(`div.export-options-lft-content > div.export-options-cell > a > i.icon.xls`)
    await expect(xlsDownloadIconLocator).toBeVisible()
    const pdfDownloadIconLocator = page.locator(`div.export-options-lft-content > div.export-options-cell > a > i.icon.pdf`)
    await expect(pdfDownloadIconLocator).toBeVisible()

    // Iterate over months, downloading & scraping each one
    let rowCount = 0
    do {
        // TODO: consider replacing this timeout with: await page.waitForLoadState("networkidle")
        await page.waitForTimeout(1000 * 3)

        const exportFilePrefix = `transactions-${buildTransactionsFileName(await fromDatePicker.getDate(),
            await toDatePicker.getDate())}`

        // STEP 5.1: DOWNLOAD XLS DATA
        const xlsDownloadPromise = page.waitForEvent("download")
        await xlsDownloadIconLocator.click()
        const xlsDownload = await xlsDownloadPromise
        await xlsDownload.saveAs(`./${exportFilePrefix}.xls`)

        // STEP 5.1: DOWNLOAD PDF DATA
        const pdfDownloadPromise = page.waitForEvent("download")
        await pdfDownloadIconLocator.click()
        const pdfDownload = await pdfDownloadPromise
        await pdfDownload.saveAs(`./${exportFilePrefix}.pdf`)

        // STEP 5: SCRAPE DATA
        const count                    = await site.getTransactionRowCount()
        let row: TransactionRow | null = null
        for (let rowIndex = count - 1; rowIndex >= 0; rowIndex--) {
            // STEP 8: SEND TRANSACTION DATA TO GREENSTAR API
            row = await site.getTransactionRow(++rowCount, rowIndex)
            if (row) {
                await row.create()
            } else {
                throw new Error(`Failed to obtain transaction row ${rowCount} / ${count}`)
            }
        }

        // If we just finished processing the first page of transactions, AND
        // This tenant has NO transactions prior to processing this page (this is the first scrape ever), THEN
        // We need an initialization transaction to bring the balance to be up-to-date
        if (row && rowCount == count && totalTransactions === 0) {

            // Obtain the last transaction row in the page (which is essentially the first transaction we've scraped)
            const firstRow = await site.getTransactionRow(1, count - 1)
            if (firstRow) {
                await createInitializationTransaction(firstRow)
            } else {
                throw new Error(`Failed to obtain transaction row 1 / ${count}`)
            }

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
