import { expressMiddleware } from "@apollo/server/express4"
import cors from "cors"
import express from "express"
import * as http from "node:http"
import { Context } from "./context.js"
import { initWebSocketServer } from "./websocket.js"
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
import { config } from "./config.js"
import multer from "multer"
import fs from "fs"
import { LargeObjectManager } from "pg-large-object"
import { rateLimit } from "express-rate-limit"
import path from "path"
import * as os from "node:os"

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

    const expressApp = express()
    const httpServer = http.createServer(expressApp)

    // Initialize WebSocket server
    initWebSocketServer(httpServer)

    const apolloServer = new ApolloServer<Context>({
        typeDefs: TypeDefinitions,
        resolvers: GraphResolvers,
        includeStacktraceInErrorResponses: true,
        formatError,
        plugins: [ ApolloServerPluginDrainHttpServer({ httpServer }), createApolloServerPluginTransactional(pgPool) ],
    })
    await apolloServer.start()

    console.info(`File uploads directory set to: ${config.uploads.path}`)
    const storage = multer.diskStorage({ destination: config.uploads.path })
    const upload  = multer({ storage, limits: { fileSize: 1024 * 1024 * 256 } })

    expressApp.use(rateLimit({
        windowMs: 60 * 1000,        // 15 minutes
        limit: 100,                 // Limit each IP to 100 requests per `window` (here, per 15 minutes).
        standardHeaders: "draft-8", // draft-6: `RateLimit-*` headers; draft-7 & draft-8: combined `RateLimit` header
        legacyHeaders: false,       // Disable the `X-RateLimit-*` headers.
        ipv6Subnet: 56,             // Set to 60 or 64 to be less aggressive, or 52 or 48 to be more aggressive
    }))

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
                    SELECT id, data_oid, size, content_type
                    FROM files AS f
                    WHERE f.tenant_id = $1
                      AND f.name = $2
                `, [ req.params.tenantID, req.params.fileName ])

                if (!rs.rows.length || !rs.rows[0].data_oid) {
                    res.status(404).type("text/plain").send("Not found")
                } else {
                    const row               = rs.rows[0]
                    const lom               = new LargeObjectManager({ pg: client })
                    const [ _size, stream ] = await lom.openAndReadableStreamAsync(row.data_oid, 16384)

                    res.status(200).type(row.content_type)

                    await new Promise<void>((resolve, reject) => {
                        stream.pipe(res)
                        stream.on("finish", resolve)
                        stream.on("error", reject)
                        res.on("close", () => {
                            stream.destroy(new Error("Client closed connection"))
                        })
                    })
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

    expressApp.post("/static/:tenantID/:scraperID/:scraperJobID",
        upload.single("file"),
        async (req, res) => {
            if (!req.file) {
                return res.status(400).type("text/plain").send("No file was uploaded.")
            }

            const tenantID = req.params.tenantID

            const fileName = `${req.params.scraperJobID}-${req.file.originalname}`

            const client = await pgPool.connect()
            await client.query("BEGIN")
            try {
                const scrapersDAL = new ScrapersDataAccessLayer(client)
                const scraper     = await scrapersDAL.fetchScraper(tenantID, req.params.scraperID)
                if (!scraper) {
                    return res.status(404).type("text/plain").send("Scraper not found.")
                }

                const tempDir      = os.tmpdir()
                const resolvedPath = path.resolve(req.file.path)
                if (!resolvedPath.startsWith(tempDir)) {
                    return res.status(400).type("text/plain").send("Invalid file path.")
                }
                const fileStream = fs.createReadStream(resolvedPath)

                const lom             = new LargeObjectManager({ pg: client })
                const [ oid, stream ] = await lom.createAndWritableStreamAsync()

                await new Promise((resolve, reject) => {
                    fileStream.pipe(stream)
                    stream.on("finish", resolve)
                    stream.on("error", reject)
                })

                const rs = await client.query(`
                    INSERT INTO files (name, data_oid, size, content_type, tenant_id)
                    VALUES ($1, $2, $3, $4, $5)
                    ON CONFLICT ON CONSTRAINT uq_files DO UPDATE
                        SET data_oid     = $2,
                            size         = $3,
                            content_type = $4
                    RETURNING id
                `, [ fileName, oid, req.file.size, req.file.mimetype, tenantID ])
                if (rs.rowCount != 1) {
                    console.error(`Incorrect number of rows affected by file upload: ${rs.rowCount}`)
                    res.status(500).type("text/plain").send("Internal server error")
                    await client.query("ROLLBACK")
                }

                await client.query("COMMIT")
                res.status(201)
                   .setHeader("x-greenstar-file-id", rs.rows[0].id)
                   .location(`/static/${tenantID}/${fileName}`)
                   .type("text/plain")
                   .send(`File uploaded (ID: ${rs.rows[0].id})`)

            } catch (e) {
                console.error(`Error retrieving file: `, e)
                res.status(500).type("text/plain").send("Internal server error")
                await client.query("ROLLBACK")
            } finally {
                if (req.file) {
                    try {
                        await fs.promises.unlink(req.file.path)
                    } catch (e) {
                        console.error("Error deleting temp file: ", req.file.path, e)
                    }
                }
                client.release()
            }
        },
    )

    expressApp.use((error: any, req: any, res: any, next: any) => {
        if (error instanceof multer.MulterError) {
            if (error.code === "LIMIT_FILE_SIZE") {
                return res.status(413).type("text/plain").send("File is too large.")
            }
        }
        next(error)
    })

    httpServer.listen(port, () => console.log(`🚀 Server started!`))
}
