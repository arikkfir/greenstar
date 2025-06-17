import { Locator } from "@playwright/test"
import { graphQLClient } from "../util/graphql-client.ts"
import { CreateTransaction } from "./query.ts"
import { DateTime } from "luxon"

export class TransactionRow {
    private readonly rowDataLocator: Locator
    private readonly rowColumnsLocator: Locator

    constructor(
        public readonly tenantID: string,
        public readonly accountID: string,
        public readonly rowCount: number,
        public readonly rowIndex: number,
        private readonly rowLocator: Locator,
    ) {
        this.rowDataLocator    = this.rowLocator.locator("> div.row")
        this.rowColumnsLocator = this.rowDataLocator.locator("> div.col")
    }

    async getDate() {
        const dateStr              = await this.rowColumnsLocator.nth(0).locator("> span").innerText()
        const [ day, month, year ] = dateStr.split("/")
        const i                    = this.rowCount - 1
        const date                 = DateTime.fromObject({
            year: Number(year),
            month: Number(month),
            day: Number(day),
            hour: Math.floor(i / 3600),
            minute: Math.floor(i / 60) % 60,
            second: i % 60,
        }, { zone: "Asia/Jerusalem" })
        return date.toJSDate()
    }

    async getReferenceID() {
        return await this.rowColumnsLocator.nth(1).locator("> span").innerText() || ""
    }

    async getDescription() {
        return await this.rowColumnsLocator.nth(2).locator("> span:not(.icon)").innerText() || ""
    }

    async getDebit() {
        const value = await this.rowColumnsLocator.nth(3).innerText()
        if (value) {
            return parseFloat(value.replace(/,/g, "").trim())
        } else {
            return 0
        }
    }

    async getCredit() {
        const value = await this.rowColumnsLocator.nth(4).innerText()
        if (value) {
            return parseFloat(value.replace(/,/g, "").trim())
        } else {
            return 0
        }
    }

    async getBalance() {
        const value = await this.rowColumnsLocator.nth(5).locator("> div > span:not(.icon)").innerText()
        if (value) {
            return parseFloat(value.replace(/,/g, "").trim())
        } else {
            return 0
        }
    }

    async getAmount() {
        const debit  = await this.getDebit()
        const credit = await this.getCredit()
        if (debit) {
            return debit
        } else if (credit) {
            return credit
        } else {
            throw new Error(
                `Transaction #${this.rowIndex + 1}: missing both debit and credit\n` +
                `- Debit: ${debit}\n` +
                `- Credit: ${credit}`,
            )
        }
    }

    async getSourceAccountID() {
        const debit = await this.getDebit()
        return debit ? this.accountID || "" : null
    }

    async getTargetAccountID() {
        const credit = await this.getCredit()
        return credit ? this.accountID || "" : null
    }

    async create() {
        const date                           = await this.getDate()
        const referenceID                    = await this.getReferenceID()
        const description                    = await this.getDescription()
        const balance                        = await this.getBalance()
        const amount                         = await this.getAmount()
        const sourceAccountID: string | null = await this.getSourceAccountID()
        const targetAccountID: string | null = await this.getTargetAccountID()
        // TODO: const expandedData   = row.locator("> div.expanded-data")

        console.info(`Creating transaction ${this.rowIndex}: `,
            `  (tenant=${this.tenantID})\n`,
            `  (date=${date})\n`,
            `  (referenceID=${referenceID})\n`,
            `  (amount=${amount})\n`,
            `  (description=${description})\n`,
            `  (sourceAccountID=${sourceAccountID})\n`,
            `  (targetAccountID=${targetAccountID})\n`,
            `  (balance=${balance})`)
        const result = await graphQLClient.mutation(CreateTransaction, {
            tx: {
                tenantID: this.tenantID,
                date: DateTime.fromJSDate(date),
                referenceID,
                amount,
                currency: "ILS",
                description,
                sourceAccountID,
                targetAccountID,
            },
        })
        if (result.error) {
            throw new Error(
                `Failed to create Transaction #${this.rowIndex + 1}\n` +
                `- Message: ${result.error.message}\n` +
                `- Network Error: ${JSON.stringify(result.error.networkError, null, 2)}\n` +
                `- GraphQL Errors: ${JSON.stringify(result.error.graphQLErrors, null, 2)}`,
            )
        }
    }
}

export async function createInitializationTransaction(row: TransactionRow) {
    const tenantID    = row.tenantID
    const date        = DateTime.fromJSDate(await row.getDate()).minus({ days: 1 })
    const referenceID = "initialization"
    const debit       = await row.getDebit()
    const credit      = await row.getCredit()

    let balance = await row.getBalance()
    if (debit) {
        balance = balance + debit
    } else if (credit) {
        balance = balance - credit
    } else {
        throw new Error(
            `Transaction #${row.rowIndex}: missing both debit and credit:\n` +
            `- Debit: ${debit}\n` +
            `- Credit: ${credit}`,
        )
    }
    const amount = Math.abs(balance)

    const currency        = "ILS"
    const description     = "Balance initialization"
    const sourceAccountID = balance < 0 ? row.accountID : null
    const targetAccountID = balance > 0 ? row.accountID : null

    console.info(`Creating BALANCE-INIT transaction ${row.rowCount}/${row.rowIndex} initialization transaction: `,
        `  (tenant=${tenantID})\n`,
        `  (date=${date})\n`,
        `  (referenceID=${referenceID})\n`,
        `  (credit=${credit})\n`,
        `  (debit=${debit})\n`,
        `  (balance=${balance})\n`,
        `  (amount=${amount})\n`,
        `  (description=${description})`)

    const result = await graphQLClient.mutation(CreateTransaction, {
        tx: { tenantID, date, referenceID, amount, currency, description, sourceAccountID, targetAccountID },
    })
    if (result.error) {
        throw new Error(
            `Failed to create initialization transaction:\n` +
            `- Message: ${result.error.message}\n` +
            `- Network Error: ${JSON.stringify(result.error.networkError, null, 2)}\n` +
            `- GraphQL Errors: ${JSON.stringify(result.error.graphQLErrors, null, 2)}`,
        )
    }
}
