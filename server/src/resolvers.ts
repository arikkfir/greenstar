import { Context } from "./context.js"
import { Account, Currency, CurrencyRate, Resolvers, Scraper, ScraperType, Tenant } from "./schema/graphql.js"
import { DateTimeScalar } from "./schema/scalars-luxon-datetime.js"
import { VoidScalar } from "./schema/scalars-void.js"

import { TransactionsSummaryResult } from "./data/transactions.js"
import pkg from "../package.json" with { type: "json" }
import { required } from "./util.js"

export interface AccountRow extends Account {
    tenantID: Tenant["id"],
    parentID: Account["id"],
}

export interface ScraperRow extends Scraper {
    tenantID: Tenant["id"],
    scraperTypeID: ScraperType["id"],
}

export const GraphResolvers: Resolvers<Context> = {
    Mutation: {
        noOp: async (_: any, _args, _ctx) => null,
        createAccount: async (_: any, args, ctx) => ctx.data.createAccount(args),
        createCurrencyRate: async (_: any, args, ctx) => ctx.data.createCurrencyRate(args),
        createScraper: async (_: any, args, ctx) => ctx.data.createScraper(args),
        createScraperRun: async (_: any, args, ctx) => ctx.data.createScraperRun(args),
        createTransaction: async (_: any, args, ctx) => ctx.data.createTransaction(args),
        createTenant: async (_: any, args, ctx) => ctx.data.createTenant(args),
        deleteAccount: async (_: any, args, ctx) => ctx.data.deleteAccount(args),
        deleteScraper: async (_: any, args, ctx) => ctx.data.deleteScraper(args),
        deleteTenant: async (_: any, args, ctx) => ctx.data.deleteTenant(args),
        deleteTransaction: async (_: any, args, ctx) => ctx.data.deleteTransaction(args),
        moveAccount: async (_: any, args, ctx) => ctx.data.moveAccount(args),
    },
    Query: {
        version: async (_: any, _args, _ctx) => pkg.version,
        tenants: async (_: any, args, ctx) => ctx.data.fetchTenants(args),
        tenant: async (_: any, args, ctx) => ctx.data.fetchTenant(args.id),
        currencies: async (_parent, _args: any, ctx): Promise<Currency[]> => ctx.data.fetchCurrencies(),
        currency: async (_parent, args, ctx) => ctx.data.fetchCurrency(args.code),
        currencyRate: async (_parent, args, ctx) => ctx.data.fetchCurrencyRate(
            args.date,
            args.sourceCurrencyCode,
            args.targetCurrencyCode,
        ),
        currencyRates: async (_parent, args, ctx): Promise<CurrencyRate[]> =>
            ctx.data.fetchCurrencyRates(
                args.startDate || undefined,
                args.endDate || undefined,
                args.sourceCurrencyCode || undefined,
                args.targetCurrencyCode || undefined,
            ),
        scraperTypes: async (_parent, _args: any, ctx) => ctx.data.fetchScraperTypes(),
    },
    Scraper: {
        runs: async (scraper, _args: any, ctx) =>
            ctx.data.fetchScraperRuns((scraper as ScraperRow).tenantID, scraper.id),
    },
    Tenant: {
        rootAccounts: async (tenant, _: any, ctx): Promise<Account[]> => ctx.data.fetchRootAccounts(tenant.id),
        account: async (tenant, args, ctx) => ctx.data.fetchAccount(tenant.id, args.id),
        accounts: async (tenant, args, ctx) => ctx.data.fetchAccounts(tenant.id, args.filter),
        accountsBalanceOverTime: async (tenant, args, ctx) =>
            ctx.data.fetchAccountsBalanceOverTime(tenant.id, args),
        transactions: async (tenant, args, ctx) =>
            ctx.data.fetchTransactions(tenant.id, "all", args),
        firstTransactionDate: async (tenant, _args, ctx) => {
            const summary: TransactionsSummaryResult = await ctx.data.fetchTransactionsSummary(tenant.id)
            return summary.firstTransactionDate || null
        },
        lastTransactionDate: async (tenant, _args, ctx) => {
            const summary = await ctx.data.fetchTransactionsSummary(tenant.id)
            return summary.lastTransactionDate || null
        },
        totalTransactions: async (tenant, _args, ctx) => {
            const summary = await ctx.data.fetchTransactionsSummary(tenant.id)
            return summary.totalCount
        },
        scrapers: async (tenant, args, ctx) => ctx.data.fetchScrapers(tenant.id, args.scraperTypeID || undefined),
        scraper: async (tenant, args, ctx) => ctx.data.fetchScraper(tenant.id, args.id),
    },
    Account: {
        parent: async (account, _args: any, ctx) => ctx.data.fetchAccount(
            (account as AccountRow).tenantID,
            (account as AccountRow).parentID,
        ),
        children: async (account, _: any, ctx) =>
            ctx.data.fetchAccountChildren((account as AccountRow).tenantID, account.id),
        tenant: async (account, _: any, ctx) => required(
            await ctx.data.fetchTenant((account as AccountRow).tenantID),
            "Tenant not found",
        ),
        incomingTransactions: async (account, args, ctx) =>
            ctx.data.fetchTransactions(
                (account as AccountRow).tenantID,
                "incoming",
                { ...args, involvingAccountID: account.id },
            ),
        outgoingTransactions: async (account, args, ctx) =>
            ctx.data.fetchTransactions(
                (account as AccountRow).tenantID,
                "outgoing",
                { ...args, involvingAccountID: account.id },
            ),
        transactions: async (account, args, ctx) =>
            ctx.data.fetchTransactions(
                (account as AccountRow).tenantID,
                "all",
                { ...args, involvingAccountID: account.id },
            ),
        balance: async (account, args, ctx) =>
            ctx.data.fetchAccountBalance((account as AccountRow).tenantID, account.id, args),
        balanceOverTime: async (account, args, ctx) =>
            ctx.data.fetchAccountBalanceOverTime(
                (account as AccountRow).tenantID,
                account.id,
                args,
            ),
    },
    DateTime: DateTimeScalar,
    Void: VoidScalar,
}
