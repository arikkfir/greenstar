import { expressMiddleware } from "@apollo/server/express4"
import cors from "cors"
import express from "express"
import * as http from "node:http"
import { Context } from "./context.js"
import type { GraphQLRequestContextExecutionDidStart } from "@apollo/server"
import {
    ApolloServer,
    ApolloServerPlugin,
    GraphQLRequestExecutionListener,
    GraphQLRequestListener,
} from "@apollo/server"
import { TypeDefinitions } from "./schema/index.js"
import { GraphResolvers } from "./resolvers.js"
import { GraphQLFormattedError, OperationTypeNode } from "graphql/index.js"
import { ApolloServerErrorCode } from "@apollo/server/errors"
import { ApolloServerPluginDrainHttpServer } from "@apollo/server/plugin/drainHttpServer"
import { pgPool } from "./data/pg_client.js"
import { DataLayerImpl, NoOpDataLayer } from "./data/data_layer.js"
import { Pool } from "pg"
import { AccountsDataAccessLayer } from "./data/accounts.js"
import { CurrenciesDataAccessLayer } from "./data/currencies.js"
import { ScrapersDataAccessLayer } from "./data/scrapers.js"
import { TenantsDataAccessLayer } from "./data/tenants.js"
import { TransactionsDataAccessLayer } from "./data/transactions.js"

function formatError(formattedError: GraphQLFormattedError, error: unknown): GraphQLFormattedError {
    if (formattedError.extensions?.code === ApolloServerErrorCode.INTERNAL_SERVER_ERROR) {
        const jsonFormattedGraphQLError = JSON.stringify(formattedError, null, 2)
        const jsonFormattedError        = JSON.stringify(error, null, 2)
        console.error(`Internal error has occurred:\n${jsonFormattedGraphQLError}\n${jsonFormattedError}`)
        return {
            message: "An internal error has occurred",
            locations: formattedError.locations,
            path: formattedError.path,
            extensions: {
                code: ApolloServerErrorCode.INTERNAL_SERVER_ERROR,
            },
        }
    }
    return formattedError
}

function createApolloServerPluginTransactional(pool: Pool): ApolloServerPlugin<Context> {
    return {
        async requestDidStart(): Promise<GraphQLRequestListener<Context> | void> {
            return {
                async executionDidStart(ctx: GraphQLRequestContextExecutionDidStart<Context>): Promise<GraphQLRequestExecutionListener<Context> | void> {
                    if (ctx.source !== `{__typename}`) {
                        const client = await pool.connect()
                        await client.query("BEGIN")

                        ctx.contextValue.data = new DataLayerImpl(
                            new AccountsDataAccessLayer(client),
                            new CurrenciesDataAccessLayer(client),
                            new ScrapersDataAccessLayer(client),
                            new TenantsDataAccessLayer(client),
                            new TransactionsDataAccessLayer(client),
                        )

                        switch (ctx.operation.operation) {
                            case OperationTypeNode.MUTATION:
                                return {
                                    async executionDidEnd(err?: Error) {
                                        try {
                                            if (err) {
                                                await client.query("ROLLBACK")
                                            } else {
                                                await client.query("COMMIT")
                                            }
                                        } finally {
                                            client.release()
                                        }
                                    },
                                }
                            case OperationTypeNode.SUBSCRIPTION:
                            case OperationTypeNode.QUERY:
                                return {
                                    async executionDidEnd() {
                                        try {
                                            await client.query("ROLLBACK")
                                        } finally {
                                            client.release()
                                        }
                                    },
                                }
                        }
                    }
                },
            }
        },
    }
}

export async function startServer() {
    if (!process.env.PORT) {
        throw new Error("environment variable PORT is required")
    }
    const port = parseInt(process.env.PORT || "bad-port")

    const expressApp   = express()
    const httpServer   = http.createServer(expressApp)
    const apolloServer = new ApolloServer<Context>({
        typeDefs: TypeDefinitions,
        resolvers: GraphResolvers,
        includeStacktraceInErrorResponses: true,
        formatError,
        plugins: [ ApolloServerPluginDrainHttpServer({ httpServer }), createApolloServerPluginTransactional(pgPool) ],
    })
    await apolloServer.start()

    expressApp.use("/graphql",
        cors<cors.CorsRequest>(),
        express.json(),
        expressMiddleware<Context>(apolloServer, {
            context: (): Promise<Context> => (Promise.resolve({ data: new NoOpDataLayer() })),
        }))

    expressApp.get("/static/:tenantID/:fileName",
        cors<cors.CorsRequest>(),
        async (req, res) => {
            const client = await pgPool.connect()
            await client.query("BEGIN")
            try {
                const rs = await client.query(`
                    SELECT id, data, size, content_type
                    FROM files AS f
                    WHERE f.tenant_id = $1
                      AND f.name = $2
                `, [ req.params.tenantID, req.params.fileName ])

                if (!rs.rows.length) {
                    res.status(404).type("text/plain").send("Not found")
                } else {
                    const row = rs.rows[0]
                    res.status(200).type(row.content_type).send(row.data)
                }
                await client.query("ROLLBACK")
            } catch (e) {
                console.error(`Error retrieving file: `, e)
                res.status(500).type("text/plain").send("Internal server error")
                await client.query("ROLLBACK")
            } finally {
                client.release()
            }
        },
    )

    expressApp.get("/static/:id([^/]+)",
        cors<cors.CorsRequest>(),
        async (req, res) => {
            const client = await pgPool.connect()
            await client.query("BEGIN")
            try {
                const rs = await client.query(`
                    SELECT id, data, size, content_type
                    FROM files AS f
                    WHERE f.id = $1
                `, [ req.params.id ])

                if (!rs.rows.length) {
                    res.status(404).type("text/plain").send("Not found")
                } else {
                    const row = rs.rows[0]
                    res.status(200).type(row.content_type).send(row.data)
                }
                await client.query("ROLLBACK")
            } catch (e) {
                console.error(`Error retrieving file: `, e)
                res.status(500).type("text/plain").send("Internal server error")
                await client.query("ROLLBACK")
            } finally {
                client.release()
            }
        },
    )

    expressApp.post("/static/:tenantID/:fileName",
        cors<cors.CorsRequest>(),
        express.raw({
            type: "*/*",
            limit: "50mb",
        }),
        async (req, res) => {
            const contentType = req.header("content-type")
            if (!contentType) {
                res.status(400).type("text/plain").send("Content-type header is missing")
            }

            const client = await pgPool.connect()
            await client.query("BEGIN")
            try {
                const rs = await client.query(`
                    INSERT INTO files (name, data, size, content_type, tenant_id)
                    VALUES ($1, $2, $3, $4, $5)
                    ON CONFLICT ON CONSTRAINT uq_files DO UPDATE
                        SET data = $2,
                            size = $3,
                            content_type = $4
                    RETURNING id
                `, [ req.params.fileName, req.body, req.body.length, contentType, req.params.tenantID ])

                if (rs.rowCount != 1) {
                    console.error(`Incorrect number of rows affected by file upload: ${rs.rowCount}`)
                    res.status(500).type("text/plain").send("Internal server error")
                    await client.query("ROLLBACK")
                }

                await client.query("COMMIT")
                res.status(201)
                   .setHeader('x-greenstar-file-id', rs.rows[0].id)
                   .location(`/static/${req.params.tenantID}/${req.params.fileName}`)
                   .type("text/plain")
                   .send(`File uploaded (ID: ${rs.rows[0].id})`)

            } catch (e) {
                console.error(`Error retrieving file: `, e)
                res.status(500).type("text/plain").send("Internal server error")
                await client.query("ROLLBACK")
            } finally {
                client.release()
            }
        },
    )

    httpServer.listen(port, () => console.log(`🚀 Server started!`))
}
