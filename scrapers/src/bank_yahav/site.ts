import { expect, Locator, Page } from "@playwright/test"
import { DatePicker } from "./date_picker.ts"
import { TransactionRow } from "./transaction_row.ts"

export class Site {
    private transactionRowsLocator: Locator

    constructor(
        private readonly page: Page,
        private readonly tenantID: string,
        private readonly accountID: string,
        private readonly username: string,
        private readonly password: string,
        private readonly pinno: string,
    ) {
        this.transactionRowsLocator = page.locator(`div.transactions > div.transactionList > div.list-item-holder > div[role=row]`)
    }

    async open() {
        await this.page.goto("/")
        await this.page.waitForLoadState("networkidle")
    }

    async login() {
        const loginRevealButtonLocator = this.page.locator(`button#mainLoginFormBtn`)
        await expect(loginRevealButtonLocator).toBeVisible()
        await loginRevealButtonLocator.click()

        const iframeLocator = this.page.locator("iframe.login-iframe").contentFrame()

        const usernameLocator = iframeLocator.locator(`input#username`)
        await expect(usernameLocator).toBeVisible()
        await usernameLocator.fill(this.username)

        const pinnoLocator = iframeLocator.locator(`input#pinno`)
        await expect(pinnoLocator).toBeVisible()
        await pinnoLocator.fill(this.pinno)

        const passwordLocator = iframeLocator.locator(`input#password`)
        await expect(passwordLocator).toBeVisible()
        await passwordLocator.fill(this.password)

        const submitButtonLocator = iframeLocator.locator(`button[type=submit].btn-primary`).getByText("כניסה")
        await expect(submitButtonLocator).toBeVisible()
        console.info("Submit button count: ", await submitButtonLocator.count())
        await submitButtonLocator.click()

        while (this.page.url().startsWith("https://www.bank-yahav.co.il/")) {
            await this.page.waitForLoadState("networkidle")
            await this.page.waitForTimeout(1000)
            console.info("Submit button count: ", await submitButtonLocator.count())
            await submitButtonLocator.click()
        }

        while (this.page.url().startsWith("https://login.yahav.co.il/login/")) {
            await this.page.waitForLoadState("networkidle")
            await this.page.waitForTimeout(5000)

            await expect(usernameLocator).toBeVisible()
            await usernameLocator.fill(this.username)

            await expect(pinnoLocator).toBeVisible()
            await pinnoLocator.fill(this.pinno)

            await expect(passwordLocator).toBeVisible()
            await passwordLocator.fill(this.password)

            const submitButtonLocator = this.page.locator("button[type=submit].btn-primary").getByText("כניסה")
            await expect(submitButtonLocator).toBeVisible()
            await submitButtonLocator.click()
            await this.page.waitForLoadState("networkidle")
        }

        if (!this.page.url().startsWith("https://digital.yahav.co.il/")) {
            throw new Error(`Unexpected URL: ${this.page.url()}`)
        }

        // WAIT FOR SPINNER TO GO AWAY
        await this.page.waitForTimeout(1000 * 10)
    }

    async logout() {
        await this.page.getByRole("button").getByText("יציאה").click()
        await this.page.waitForLoadState("networkidle")
        await this.page.waitForTimeout(1000 * 3)
    }

    async openTransactionsPage() {
        const checkingAccountLinkLocator = this.page.locator(`a[href="#/main/accounts/current/"]`)
        await expect(checkingAccountLinkLocator).toBeVisible()
        await checkingAccountLinkLocator.click()
        await this.page.waitForLoadState("networkidle")
    }

    createFromDatePicker() {
        return new DatePicker(this.page, "datePickerPositionFrom")
    }

    createToDatePicker() {
        return new DatePicker(this.page, "datePickerPositionTo")
    }

    async getTransactionRowCount() {
        return await this.transactionRowsLocator.count()
    }

    async getTransactionRow(rowCount: number, rowIndex: number) {
        const rowLocator = this.transactionRowsLocator.nth(rowIndex)
        if (await rowLocator.count() > 0) {
            return new TransactionRow(
                this.tenantID,
                this.accountID,
                rowCount,
                rowIndex,
                rowLocator,
            )
        } else {
            return null
        }
    }
}
