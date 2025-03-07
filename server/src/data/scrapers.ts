import { PoolClient } from "pg"
import {
    MutationCreateScraperArgs,
    MutationDeleteScraperArgs,
    Scraper,
    ScraperParameter,
    ScraperParameterType,
    ScraperType,
    ScraperTypeParameter,
    Tenant,
} from "../schema/graphql.js"
import { convertObjectKeysToCamelCase } from "./pg_client.js"

export class ScrapersDataAccessLayer {
    constructor(private readonly pg: PoolClient) {
    }

    async createScraper(args: MutationCreateScraperArgs): Promise<Scraper> {
        const { tenantID, id, displayName, scraperTypeID, parameters } = args

        const res = await this.pg.query(`
            INSERT INTO scrapers (tenant_id, scraper_type_id, id, display_name)
            VALUES ($1, $2, $3, $4)
            RETURNING id
        `, [ tenantID, scraperTypeID, id, displayName ])
        if (res.rowCount != 1 || res.rows.length != 1) {
            throw new Error(`Incorrect number of affected or rows returned: ${JSON.stringify(res)}`)
        }
        const scraperID = res.rows[0].id as number

        if (parameters.length) {
            let parameterArgs = []
            for (let p of parameters) {
                parameterArgs.push(tenantID, scraperTypeID, scraperID, p.scraperTypeParameterID, p.value)
            }
            const values = parameters.reduce((prev, curr, i, arr) => {
                const v = prev + `$${i * 5 + 1}, $${i * 5 + 2}, $${i * 5 + 3}, $${i * 5 + 4}, $${i * 5 + 5}`
                return v + (i < arr.length - 1 ? `), (` : ``)
            }, "")
            const sql    = `
                INSERT
                INTO scraper_parameters (tenant_id, scraper_type_id, scraper_id, scraper_type_parameter_id, value)
                VALUES (${values})
            `

            const pres = await this.pg.query(sql, parameterArgs)
            if (pres.rowCount != parameters.length) {
                throw new Error(`Incorrect number of rows affected: ${JSON.stringify(pres)}`)
            }
        }

        return (await this.fetchScraper(tenantID, scraperTypeID, res.rows[0].id))!
    }

    async deleteScraper(args: MutationDeleteScraperArgs): Promise<void> {
        const res = await this.pg.query(`
            DELETE
            FROM scrapers
            WHERE tenant_id = $1 AND scraper_type_id = $2 AND id = $3
        `, [ args.tenantID, args.scraperTypeID, args.id ])
        if (res.rowCount != 1) {
            throw new Error(`Incorrect number of rows affected: ${JSON.stringify(res)}`)
        }
    }

    async fetchScraperParameterType(id: ScraperParameterType["id"]): Promise<ScraperParameterType | null> {
        const res = await this.pg.query(`
                    SELECT id,
                           created_at   AS createdat,
                           updated_at   AS updatedat,
                           display_name AS displayname
                    FROM scraper_parameter_types
                    WHERE id = $1
            `,
            [ id ])
        if (res.rows.length == 0) {
            return null
        } else if (res.rows.length > 1) {
            throw new Error(`Too many rows returned: ${JSON.stringify(res)}`)
        } else {
            return res.rows.map(convertObjectKeysToCamelCase)[0] as ScraperParameterType
        }
    }

    async fetchScraperParameterTypes(): Promise<ScraperParameterType[]> {
        const res = await this.pg.query(`
            SELECT id, created_at, updated_at, display_name
            FROM scraper_parameter_types
            ORDER BY display_name, id
        `)
        return res.rows.map(convertObjectKeysToCamelCase) as ScraperParameterType[]
    }

    async fetchScraperType(id: ScraperType["id"]): Promise<ScraperType | null> {
        const res = await this.pg.query(`
                    SELECT id,
                           created_at   AS createdat,
                           updated_at   AS updatedat,
                           display_name AS displayname
                    FROM scraper_types
                    WHERE id = $1
            `,
            [ id ])
        if (res.rows.length == 0) {
            return null
        } else if (res.rows.length > 1) {
            throw new Error(`Too many rows returned: ${JSON.stringify(res)}`)
        } else {
            return res.rows.map(convertObjectKeysToCamelCase)[0] as ScraperType
        }
    }

    async fetchScraperTypes(): Promise<ScraperType[]> {
        const res = await this.pg.query(`
                    SELECT id, created_at, updated_at, display_name
                    FROM scraper_types
                    ORDER BY display_name, id
            `,
            [])
        return res.rows.map(convertObjectKeysToCamelCase) as ScraperType[]
    }

    async fetchScraperTypeParameter(
        scraperTypeID: ScraperType["id"],
        id: ScraperTypeParameter["id"],
    ): Promise<ScraperTypeParameter | null> {
        const res = await this.pg.query(`
                    SELECT scraper_type_id,
                           id,
                           created_at,
                           updated_at,
                           display_name,
                           scraper_parameter_type_id
                    FROM scraper_type_parameters
                    WHERE scraper_type_id = $1 AND id = $2
            `,
            [ scraperTypeID, id ])
        if (res.rows.length == 0) {
            return null
        } else if (res.rows.length > 1) {
            throw new Error(`Too many rows returned: ${JSON.stringify(res)}`)
        } else {
            return res.rows.map(convertObjectKeysToCamelCase)[0] as ScraperTypeParameter
        }
    }

    async fetchScraperTypeParameters(scraperTypeID: ScraperType["id"]): Promise<ScraperTypeParameter[]> {
        const res = await this.pg.query(`
                    SELECT scraper_type_id,
                           id,
                           created_at,
                           updated_at,
                           display_name,
                           scraper_parameter_type_id
                    FROM scraper_type_parameters
                    WHERE scraper_type_id = $1
            `,
            [ scraperTypeID ])
        return res.rows.map(convertObjectKeysToCamelCase) as ScraperTypeParameter[]
    }

    async fetchScraper(
        tenantID: Tenant["id"],
        scraperTypeID: ScraperType["id"],
        id: Scraper["id"],
    ): Promise<Scraper | null> {
        const res = await this.pg.query(`
                    SELECT tenant_id,
                           scraper_type_id,
                           id,
                           created_at,
                           updated_at,
                           display_name
                    FROM scrapers
                    WHERE tenant_id = $1 AND scraper_type_id = $2 AND id = $3
            `,
            [ tenantID, scraperTypeID, id ])
        if (res.rows.length == 0) {
            return null
        } else if (res.rows.length > 1) {
            throw new Error(`Too many rows returned: ${JSON.stringify(res)}`)
        } else {
            return res.rows.map(convertObjectKeysToCamelCase)[0] as Scraper
        }
    }

    async fetchScrapers(tenantID: Tenant["id"]): Promise<Scraper[]> {
        const res = await this.pg.query(`
                    SELECT tenant_id,
                           scraper_type_id,
                           id,
                           created_at,
                           updated_at,
                           display_name
                    FROM scrapers
                    WHERE tenant_id = $1
                    ORDER BY display_name, id
            `,
            [ tenantID ])
        return res.rows.map(convertObjectKeysToCamelCase) as Scraper[]
    }

    async fetchScraperParameters(
        tenantID: Tenant["id"],
        scraperTypeID: ScraperType["id"],
        scraperID: Scraper["id"],
    ): Promise<ScraperParameter[]> {
        const res = await this.pg.query(`
                    SELECT tenant_id,
                           scraper_type_id,
                           scraper_id,
                           scraper_type_parameter_id,
                           created_at,
                           updated_at,
                           value
                    FROM scraper_parameters
                    WHERE tenant_id = $1 AND scraper_type_id = $2 AND scraper_id = $3
                    ORDER BY tenant_id, scraper_type_id, scraper_id, scraper_type_parameter_id
            `,
            [ tenantID, scraperTypeID, scraperID ])
        return res.rows.map(convertObjectKeysToCamelCase) as ScraperParameter[]
    }
}