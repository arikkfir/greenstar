import { PoolClient } from "pg"
import {
    MutationCreateScraperArgs,
    MutationCreateScraperRunArgs,
    MutationDeleteScraperArgs,
    Scraper,
    ScraperParameter,
    ScraperParameterType,
    ScraperRun,
    ScraperType,
    ScraperTypeParameter,
    ScraperTypeParameterSource,
    Tenant,
} from "../schema/graphql.js"
import { convertObjectKeysToCamelCase } from "./pg_client.js"
import { k8s } from "./k8s_client.js"
import { DateTime, Duration } from "luxon"
import { config } from "../config.js"
import { ScraperRow } from "../resolvers.js"

export class ScrapersDataAccessLayer {
    constructor(private readonly pg: PoolClient) {
    }

    async createScraper(args: MutationCreateScraperArgs): Promise<Scraper> {
        const { tenantID, id, displayName, scraperTypeID, parameters } = args

        const parametersForDB: ScraperParametersInDB = parameters.reduce(
            (prev, curr) => ({ ...prev, [curr.parameterID]: curr.value }),
            {} as ScraperParametersInDB,
        )

        // Save scraper with empty parameters
        const res = await this.pg.query(`
            INSERT INTO scrapers (id, display_name, scraper_type_id, parameters, tenant_id)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id
        `, [ id || null, displayName, scraperTypeID, parametersForDB, tenantID ])
        if (res.rowCount != 1 || res.rows.length != 1) {
            throw new Error(`Incorrect number of affected or rows returned: ${JSON.stringify(res)}`)
        }
        const scraperID = res.rows[0].id as string

        // Delete (if existing) the CronJob associated with this scraper ID (can happen when scraper ID is reused)
        await this.deleteCronJobs(scraperJobLabelSelector(tenantID, scraperID))

        // Create the corresponding CronJob
        await k8s.api.batch.createNamespacedCronJob({
            namespace: k8s.namespace,
            body: {
                apiVersion: "batch/v1",
                kind: "CronJob",
                metadata: {
                    name: `scraper-${tenantID}-${scraperID}`.toLowerCase(),
                    labels: scraperJobLabels(tenantID, scraperID),
                    ownerReferences: [
                        {
                            apiVersion: 'apps/v1',
                            kind: 'Deployment',
                            name: k8s.deployment.metadata!.name!,
                            uid: k8s.deployment.metadata!.uid!,
                            controller: true,
                        }
                    ]
                },
                spec: {
                    schedule: "0 2 * * *",
                    suspend: true, // TODO: make suspension controllable (property of the scraper)
                    concurrencyPolicy: "Forbid",
                    failedJobsHistoryLimit: 2,
                    successfulJobsHistoryLimit: 1,
                    jobTemplate: {
                        metadata: {
                            labels: scraperJobLabels(tenantID, scraperID),
                        },
                        spec: {
                            backoffLimit: 0,
                            activeDeadlineSeconds: Duration.fromObject({ minutes: 15 }).seconds,
                            template: {
                                metadata: {
                                    labels: scraperJobLabels(tenantID, scraperID),
                                },
                                spec: {
                                    enableServiceLinks: false,
                                    restartPolicy: "OnFailure",
                                    preemptionPolicy: "Never",
                                    containers: [
                                        {
                                            args: [ "--", `scrapers/${scraperTypeID}/scraper.ts` ],
                                            env: [
                                                { name: "NODE_ENV", value: process.env.NODE_ENV || "development" },
                                                { name: "TENANT_ID", value: tenantID },
                                                { name: "SCRAPER_ID", value: scraperID },
                                                ...args.parameters.map(
                                                    ({ parameterID: name, value }) => ({
                                                        name,
                                                        value,
                                                    })),
                                            ],
                                            image: `${config.scrapers.image.repository}:${config.scrapers.image.tag}}`,
                                            imagePullPolicy: config.scrapers.image.pullPolicy,
                                            name: "scraper",
                                            resources: {
                                                requests: {
                                                    cpu: "1000m",
                                                    memory: "1Gi",
                                                },
                                                limits: {
                                                    cpu: "2000m",
                                                    memory: "2Gi",
                                                },
                                            },
                                        },
                                    ],
                                },
                            },
                        },
                    },
                },
            },
        })

        return (await this.fetchScraper(tenantID, scraperID))!
    }

    async deleteScraper(args: MutationDeleteScraperArgs): Promise<void> {
        const { tenantID, id } = args

        await this.deleteCronJobs(scraperJobLabelSelector(tenantID, id))

        const res = await this.pg.query(`
            DELETE
            FROM scrapers
            WHERE tenant_id = $1 AND id = $2
        `, [ tenantID, id ])
        if (res.rowCount != 1) {
            throw new Error(`Incorrect number of rows affected: ${JSON.stringify(res)}`)
        }
    }

    async fetchScraperType(id: ScraperType["id"]): Promise<ScraperType | null> {
        const res = await this.pg.query(`
                    SELECT st.id,
                           st.created_at,
                           st.updated_at,
                           st.display_name,
                           st.parameters
                    FROM scraper_types AS st
                    WHERE id = $1
            `,
            [ id ])
        if (res.rows.length == 0) {
            return null
        } else if (res.rows.length > 1) {
            throw new Error(`Too many rows returned: ${JSON.stringify(res)}`)
        } else {
            return res.rows
                      .map(convertObjectKeysToCamelCase)
                      .map((r): ScraperType => ({
                          ...(r as ScraperType),
                          parameters: mapScraperTypeParameters(r.parameters),
                      }))[0]
        }
    }

    async fetchScraperTypes(): Promise<ScraperType[]> {
        const res = await this.pg.query(`
                    SELECT st.id,
                           st.created_at,
                           st.updated_at,
                           st.display_name,
                           st.parameters
                    FROM scraper_types AS st
                    ORDER BY st.display_name, st.id
            `,
            [])
        return res.rows
                  .map(convertObjectKeysToCamelCase)
                  .map((r): ScraperType => ({
                      ...(r as ScraperType),
                      parameters: mapScraperTypeParameters(r.parameters),
                  }))
    }

    async fetchScraper(tenantID: Tenant["id"], id: Scraper["id"]): Promise<Scraper | null> {
        const res = await this.pg.query(`
                    SELECT s.id              AS s_id,
                           s.created_at      AS s_created_at,
                           s.updated_at      AS s_updated_at,
                           s.display_name    AS s_display_name,
                           s.scraper_type_id AS s_scraper_type_id,
                           s.parameters      AS s_parameters,
                           s.tenant_id       AS s_tenant_id,
                           st.id             AS st_id,
                           st.created_at     AS st_created_at,
                           st.updated_at     AS st_updated_at,
                           st.display_name   AS st_display_name,
                           st.parameters     AS st_parameters
                    FROM scrapers AS s
                             JOIN scraper_types AS st ON s.scraper_type_id = st.id
                    WHERE s.tenant_id = $1 AND s.id = $2
            `,
            [ tenantID, id ])
        if (res.rows.length == 0) {
            return null
        } else if (res.rows.length > 1) {
            throw new Error(`Too many rows returned: ${JSON.stringify(res)}`)
        } else {
            return res.rows.map(mapScraper)[0]
        }
    }

    async fetchScrapers(tenantID: Tenant["id"], scraperTypeID?: ScraperType["id"]): Promise<Scraper[]> {
        const res = await this.pg.query(`
                    SELECT s.id              AS s_id,
                           s.created_at      AS s_created_at,
                           s.updated_at      AS s_updated_at,
                           s.display_name    AS s_display_name,
                           s.scraper_type_id AS s_scraper_type_id,
                           s.parameters      AS s_parameters,
                           s.tenant_id       AS s_tenant_id,
                           st.id             AS st_id,
                           st.created_at     AS st_created_at,
                           st.updated_at     AS st_updated_at,
                           st.parameters     AS st_parameters,
                           st.display_name   AS st_display_name
                    FROM scrapers AS s
                             JOIN scraper_types AS st ON s.scraper_type_id = st.id
                    WHERE TRUE
                      AND tenant_id = $1
                        ${scraperTypeID ? `AND scraper_type_id = $2` : ""}
                    ORDER BY s.display_name, st.display_name, s.id
            `,
            scraperTypeID ? [ tenantID, scraperTypeID ] : [ tenantID ])
        return res.rows.map(mapScraper)
    }

    async createScraperRun(args: MutationCreateScraperRunArgs): Promise<ScraperRun> {
        const { tenantID, scraperID, parameters: parametersInput } = args

        // Find our scraper
        const scraper = await this.fetchScraper(tenantID, scraperID)
        if (!scraper) {
            throw new Error(`Could not find scraper ${scraperID}`)
        }
        const scraperParametersMap: ScraperParametersInDB = scraper.parameters.reduce(
            (prev, curr) => ({ ...prev, [curr.parameter.id]: curr.value }),
            {} as ScraperParametersInDB,
        )

        // Convert given scraper run parameters array into a map of key->value
        const scraperRunParametersMap: ScraperParametersInDB = parametersInput.reduce(
            (prev, curr) => ({ ...prev, [curr.parameterID]: curr.value }),
            {} as ScraperParametersInDB,
        )

        // Save scraper run with scraper+scraperRun parameters
        const res = await this.pg.query(`
            INSERT INTO scraper_runs (scraper_id, parameters)
            VALUES ($1, $2)
            RETURNING id
        `, [ scraperID, { ...scraperParametersMap, ...scraperRunParametersMap } ])
        if (res.rowCount != 1 || res.rows.length != 1) {
            throw new Error(`Incorrect number of affected or rows returned: ${JSON.stringify(res)}`)
        }

        return (await this.fetchScraperRun(tenantID, res.rows[0].id))!
    }

    async fetchScraperRun(tenantID: Tenant["id"], id: Scraper["id"]): Promise<ScraperRun | null> {
        const res = await this.pg.query(`
                    SELECT sr.id             AS sr_id,
                           sr.created_at     AS sr_created_at,
                           sr.updated_at     AS sr_updated_at,
                           sr.scraper_id     AS sr_scraper_id,
                           sr.parameters     AS sr_parameters,
                           s.id              AS s_id,
                           s.created_at      AS s_created_at,
                           s.updated_at      AS s_updated_at,
                           s.display_name    AS s_display_name,
                           s.scraper_type_id AS s_scraper_type_id,
                           s.tenant_id       AS s_tenant_id,
                           s.parameters      AS s_parameters,
                           st.id             AS st_id,
                           st.created_at     AS st_created_at,
                           st.updated_at     AS st_updated_at,
                           st.display_name   AS st_display_name,
                           st.parameters     AS st_parameters
                    FROM scraper_runs AS sr
                             JOIN scrapers AS s ON sr.scraper_id = s.id
                             JOIN scraper_types AS st ON s.scraper_type_id = st.id
                    WHERE TRUE
                      AND s.tenant_id = $1
                      AND sr.id = $2
            `,
            [ tenantID, id ])
        if (res.rows.length == 0) {
            return null
        } else if (res.rows.length > 1) {
            throw new Error(`Too many rows returned: ${JSON.stringify(res)}`)
        } else {
            return res.rows.map(mapScraperRun)[0]
        }
    }

    async fetchScraperRuns(tenantID: string, id: string): Promise<ScraperRun[]> {
        const res = await this.pg.query(`
                    SELECT sr.id             AS sr_id,
                           sr.created_at     AS sr_created_at,
                           sr.updated_at     AS sr_updated_at,
                           sr.scraper_id     AS sr_scraper_id,
                           sr.parameters     AS sr_parameters,
                           s.id              AS s_id,
                           s.created_at      AS s_created_at,
                           s.updated_at      AS s_updated_at,
                           s.display_name    AS s_display_name,
                           s.scraper_type_id AS s_scraper_type_id,
                           s.tenant_id       AS s_tenant_id,
                           s.parameters      AS s_parameters,
                           st.id             AS st_id,
                           st.created_at     AS st_created_at,
                           st.updated_at     AS st_updated_at,
                           st.display_name   AS st_display_name,
                           st.parameters     AS st_parameters
                    FROM scraper_runs AS sr
                             JOIN scrapers AS s ON sr.scraper_id = s.id AND s.tenant_id = $1
                             JOIN scraper_types AS st ON s.scraper_type_id = st.id
                    WHERE s.id = $2
            `,
            [ tenantID, id ])
        return res.rows.map(mapScraperRun)
    }

    private async deleteCronJobs(labelSelector: string) {
        let deletionStartTime = DateTime.now()
        while (true) {
            if (deletionStartTime.diffNow("seconds").seconds > 120) {
                throw new Error(`Timed out waiting for CronJob to be deleted`)
            }
            const existingCronJobsList = await k8s.api.batch.listNamespacedCronJob({
                namespace: k8s.namespace,
                labelSelector,
            })
            if (existingCronJobsList.items.length == 0) {
                break
            } else {
                for (let cronJob of existingCronJobsList.items) {
                    if (cronJob.metadata!.deletionTimestamp === undefined) {
                        try {
                            await k8s.api.batch.deleteNamespacedCronJob({
                                namespace: cronJob.metadata!.namespace!,
                                name: cronJob.metadata!.name!,
                            })
                        } catch (err) {
                            if (err instanceof k8s.APIException && err.code !== 404) {
                                throw err
                            }
                        }
                    }
                }
            }
            await new Promise(resolve => setTimeout(resolve, 1000))
        }
    }
}

interface ScraperTypeParametersInDB {
    [k: string]: ScraperTypeParameterInDB
}

interface ScraperTypeParameterInDB {
    displayName: string,
    type: string,
    source: string,
}

interface ScraperParametersInDB {
    [k: string]: string
}

function mapScraperTypeParameters(parameters: ScraperTypeParametersInDB): ScraperTypeParameter[] {
    return Object.entries(parameters)
                 .map(([ key, value ]) => mapScraperTypeParameter(key, value))
}

function mapScraperTypeParameter(id: string, p: ScraperTypeParameterInDB): ScraperTypeParameter {
    let source: ScraperTypeParameterSource
    if (p.source.toUpperCase() == ScraperTypeParameterSource.System.toUpperCase()) {
        source = ScraperTypeParameterSource.System
    } else if (p.source.toUpperCase() == ScraperTypeParameterSource.User.toUpperCase()) {
        source = ScraperTypeParameterSource.User
    } else {
        throw new Error(`Invalid scraper parameter source: ${p.source}`)
    }

    let type: ScraperParameterType
    if (p.type.toUpperCase() == ScraperParameterType.Account.toUpperCase()) {
        type = ScraperParameterType.Account
    } else if (p.type.toUpperCase() == ScraperParameterType.Boolean.toUpperCase()) {
        type = ScraperParameterType.Boolean
    } else if (p.type.toUpperCase() == ScraperParameterType.Date.toUpperCase()) {
        type = ScraperParameterType.Date
    } else if (p.type.toUpperCase() == ScraperParameterType.Float.toUpperCase()) {
        type = ScraperParameterType.Float
    } else if (p.type.toUpperCase() == ScraperParameterType.Integer.toUpperCase()) {
        type = ScraperParameterType.Integer
    } else if (p.type.toUpperCase() == ScraperParameterType.Password.toUpperCase()) {
        type = ScraperParameterType.Password
    } else if (p.type.toUpperCase() == ScraperParameterType.String.toUpperCase()) {
        type = ScraperParameterType.String
    } else {
        throw new Error(`Invalid scraper parameter type: ${p.type}`)
    }

    return { id, displayName: p.displayName, type, source }
}

function scraperJobLabels(tenantID: Tenant["id"], scraperID: Scraper["id"]): { [key: string]: string } {
    return {
        "greenstar.finance/tenant-id": tenantID,
        "greenstar.finance/scraper-id": scraperID,
    }
}

function scraperJobLabelSelector(tenantID: Tenant["id"], scraperID: Scraper["id"]): string {
    return Object.entries(scraperJobLabels(tenantID, scraperID))
                 .map(([ key, value ]) => `${key}=${value}`)
                 .join(",")
}

function mapScraper(r: any): ScraperRow {
    const dbScraperTypeParameters: ScraperTypeParametersInDB = r.st_parameters
    const dbScraperParametersMap: ScraperParametersInDB      = r.s_parameters

    const scraperTypeParametersMap: { [k: string]: ScraperTypeParameter } =
              Object.entries(dbScraperTypeParameters)
                    .map(([ id, p ]) => mapScraperTypeParameter(id, p))
                    .reduce(
                        (prev: { [k: string]: ScraperTypeParameter }, curr) =>
                            ({ ...prev, [curr.id]: curr }),
                        {},
                    )

    const scraperParametersArray: ScraperParameter[] =
              Object.entries(dbScraperParametersMap).map(([ id, value ]) => {
                  if (dbScraperTypeParameters[id] === undefined) {
                      throw new Error(`Scraper parameter ${id} not found in scraper type parameters`)
                  }
                  return { parameter: scraperTypeParametersMap[id], value }
              })

    return {
        id: r.s_id,
        createdAt: r.s_created_at,
        updatedAt: r.s_updated_at,
        displayName: r.s_display_name,
        type: {
            id: r.st_id,
            createdAt: r.st_created_at,
            updatedAt: r.st_updated_at,
            displayName: r.st_display_name,
            parameters: mapScraperTypeParameters(r.st_parameters),
        },
        parameters: scraperParametersArray,
        tenantID: r.s_tenant_id,
        scraperTypeID: r.s_scraper_type_id,
    } as ScraperRow
}

function mapScraperRun(r: any): ScraperRun {
    const dbScraperTypeParameters: ScraperTypeParametersInDB              = r.st_parameters
    const dbScraperRunParametersMap: ScraperParametersInDB                = r.sr_parameters
    const scraperTypeParametersMap: { [k: string]: ScraperTypeParameter } =
              Object.entries(dbScraperTypeParameters)
                    .map(([ id, p ]) => mapScraperTypeParameter(id, p))
                    .reduce(
                        (prev: { [k: string]: ScraperTypeParameter }, curr) =>
                            ({ ...prev, [curr.id]: curr }),
                        {},
                    )
    const scraperRunParametersArray: ScraperParameter[]                   =
              Object.entries(dbScraperRunParametersMap).map(([ id, value ]) => {
                  if (dbScraperTypeParameters[id] === undefined) {
                      throw new Error(`Scraper run parameter ${id} not found in scraper type parameters`)
                  }
                  return { parameter: scraperTypeParametersMap[id], value }
              })
    return {
        id: r.sr_id,
        createdAt: r.sr_created_at,
        updatedAt: r.sr_updated_at,
        scraper: mapScraper(r),
        parameters: scraperRunParametersArray,
    }
}
