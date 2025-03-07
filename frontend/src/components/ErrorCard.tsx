import {
    Button,
    Card,
    CardActions,
    CardContent,
    CardHeader,
    Collapse,
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableRow,
    Typography,
} from "@mui/material"
import ErrorIcon from "@mui/icons-material/Error"
import ExpandMoreIcon from "@mui/icons-material/ExpandMore"
import { ReactNode, useMemo, useState } from "react"
import ExpandLessIcon from "@mui/icons-material/ExpandLess"
import { ApolloError, ServerError } from "@apollo/client"
import { GraphQLFormattedError } from "graphql/error"
import type { ServerParseError } from "@apollo/client/link/http"

export interface ExtraItem {
    key: string,
    label: string,
    value: string | ReactNode
}

export type ExtraItems = ExtraItem[]

function apolloErrorExtra(error: ApolloError): ExtraItems {
    const extra: ExtraItems = []

    if (error.graphQLErrors) {
        for (let i = 0; i < error.graphQLErrors.length; i++) {
            const gqlErr: GraphQLFormattedError = error.graphQLErrors[i]
            const keyPrefix                     = `error.graphQLErrors[${i}]`
            const labelPrefix                   = `GraphQL Error ${i}`
            extra.push(
                {
                    key: `${keyPrefix}.message`,
                    label: `${labelPrefix} message`,
                    value: gqlErr.message,
                },
                {
                    key: `${keyPrefix}.path`,
                    label: `${labelPrefix} path`,
                    value: <pre>{gqlErr.path?.join("/")}</pre>,
                },
            )
            if (gqlErr.locations?.length) {
                for (let j = 0; j < gqlErr.locations.length; j++) {
                    const l = gqlErr.locations[j]
                    extra.push({
                        key: `${keyPrefix}.locations[${j}]`,
                        label: `${labelPrefix} Location ${j}`,
                        value: l.line + ":" + l.column,
                    })
                }
            }
            if (gqlErr.extensions?.length) {
                for (let e in gqlErr.extensions) {
                    if (gqlErr.extensions.hasOwnProperty(e)) {
                        const ext = gqlErr.extensions[e]
                        extra.push({
                            key: `${keyPrefix}.extensions[${e}]`,
                            label: `${labelPrefix} Extension ${e}`,
                            value: JSON.stringify(ext, null, 2),
                        })
                    }
                }
            }
        }
    }

    if (error.protocolErrors) {
        for (let i = 0; i < error.protocolErrors.length; i++) {
            const protocolErr: GraphQLFormattedError = error.protocolErrors[i]
            const keyPrefix                          = `error.protocolErrors[${i}]`
            const labelPrefix                        = `Protocol Error ${i}`
            extra.push(
                {
                    key: `${keyPrefix}.message`,
                    label: `${labelPrefix} message`,
                    value: protocolErr.message,
                },
                {
                    key: `${keyPrefix}.path`,
                    label: `${labelPrefix} path`,
                    value: protocolErr.path,
                },
            )
            if (protocolErr.locations?.length) {
                for (let j = 0; j < protocolErr.locations.length; j++) {
                    const l = protocolErr.locations[j]
                    extra.push({
                        key: `${keyPrefix}.locations[${j}]`,
                        label: `${labelPrefix} Location ${j}`,
                        value: l.line + ":" + l.column,
                    })
                }
            }
            if (protocolErr.extensions?.length) {
                for (let e in protocolErr.extensions) {
                    if (protocolErr.extensions.hasOwnProperty(e)) {
                        const ext = protocolErr.extensions[e]
                        extra.push({
                            key: `${keyPrefix}.extensions[${e}]`,
                            label: `${labelPrefix} Extension ${e}`,
                            value: JSON.stringify(ext, null, 2),
                        })
                    }
                }
            }
        }
    }

    if (error.clientErrors) {
        for (let i = 0; i < error.clientErrors.length; i++) {
            const clientErr: Error = error.clientErrors[i]
            const keyPrefix        = `error.clientErrors[${i}]`
            const labelPrefix      = `Client Error ${i}`
            extra.push(
                {
                    key: `${keyPrefix}.name`,
                    label: `${labelPrefix} name`,
                    value: clientErr.name,
                },
                {
                    key: `${keyPrefix}.message`,
                    label: `${labelPrefix} message`,
                    value: clientErr.message,
                },
                {
                    key: `${keyPrefix}.stack`,
                    label: `${labelPrefix} stack`,
                    value: <pre>{clientErr.stack}</pre>,
                },
                {
                    key: `${keyPrefix}.cause`,
                    label: `${labelPrefix} Cause`,
                    value: <pre>{clientErr.cause + ""}</pre>,
                },
            )
        }
    }

    if (error.networkError) {
        const netErr = error.networkError
        if ("response" in netErr) {
            const response = netErr["response"] as Response
            response.headers.forEach((v, k) => {
                extra.push({
                    key: `error.networkError.response.headers[${k}]`,
                    label: `Network Error Response Header '${k}'`,
                    value: <pre>{v}</pre>,
                })
            })
        }
        if ("statusCode" in netErr) {
            const statusCode = netErr["statusCode"] as number
            extra.push({
                key: `error.networkError.statusCode`,
                label: `Network Error Status code`,
                value: <pre>{statusCode}</pre>,
            })
        }
        if ("result" in netErr) {
            const result = netErr["result"] as ServerError["result"]
            extra.push({
                key: `error.networkError.result`,
                label: `Network Error Result`,
                value: <pre>{JSON.stringify(result, null, 2)}</pre>,
            })
        }
        if ("bodyText" in netErr) {
            const bodyText = netErr["bodyText"] as ServerParseError["bodyText"]
            extra.push({
                key: `error.networkError.bodyText`,
                label: `Network Error Body text`,
                value: <pre>{bodyText}</pre>,
            })
        }
    }

    if (error.cause) {
        const causeErr = error.cause
        if (causeErr.extensions) {
            if (Array.isArray(causeErr.extensions)) {
                for (let i = 0; i < causeErr.extensions.length; i++) {
                    const ext = causeErr.extensions[i]
                    extra.push({
                        key: `error.cause.extensions[${i}]`,
                        label: `Error Cause Extension ${i}`,
                        value: JSON.stringify(ext, null, 2),
                    })
                }
            } else {
                const ext = causeErr.extensions
                extra.push({
                    key: `error.cause.extensions`,
                    label: `Cause Error Extensions`,
                    value: JSON.stringify(ext, null, 2),
                })
            }
        }
    }

    if (error.extraInfo) {
        const extraInfo = error.extraInfo
        extra.push({
            key: `error.extraInfo`,
            label: `Error Extra Info`,
            value: JSON.stringify(extraInfo, null, 2),
        })
    }

    return extra
}

export function buildErrorExtra(error: Error): ExtraItems {
    const extra: ExtraItems = []

    extra.push(
        {
            key: "error.name",
            label: "Error name",
            value: error.name,
        },
        {
            key: "error.message",
            label: "Error message",
            value: error.message,
        },
        {
            key: "error.stack",
            label: "Error stack",
            value: <pre style={{ textAlign: "left" }}>{error.stack}</pre>,
        },
    )

    if (error.cause instanceof Error) {
        const causeErr = error.cause
        extra.push({
            key: `error.cause.message`,
            label: `Cause Error Message`,
            value: buildErrorExtra(causeErr).map((kv) => `${kv.key}: ${kv.value}\n`).join("\n"),
        })
    }

    if (error instanceof ApolloError) {
        extra.push(...apolloErrorExtra(error))
    }

    return extra
}

export interface ErrorCardProps {
    title: string
    subTitle?: string
    text?: string
    extra?: ExtraItems
}

export function ErrorCard({ title, subTitle, text, extra }: ErrorCardProps) {
    const [ expanded, setExpanded ] = useState(false)
    const handleExpandClick         = () => setExpanded(!expanded)

    const expandIcon = useMemo(
        () => (expanded ? <ExpandLessIcon fontSize="inherit" /> : <ExpandMoreIcon fontSize="inherit" />),
        [ expanded ],
    )
    return (
        <Card elevation={3}>
            <CardHeader
                avatar={<ErrorIcon color="error" />}
                title={title}
                subheader={subTitle}
            />
            <CardContent>
                <Typography variant="body2" sx={{ color: "text.secondary", fontWeight: "bold" }}>{text}</Typography>
            </CardContent>
            {extra && extra.length > 0 && (
                <>
                    <CardActions disableSpacing>
                        <div style={{ marginLeft: "auto" }}>
                            <Button variant="contained" endIcon={expandIcon} onClick={handleExpandClick}>
                                Details
                            </Button>
                        </div>
                    </CardActions>
                    <Collapse in={expanded} timeout="auto" unmountOnExit>
                        <CardContent>
                            <Table>
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Key</TableCell>
                                        <TableCell>Name</TableCell>
                                        <TableCell>Value</TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {extra.map((kv) => (
                                        <TableRow key={kv.key}>
                                            <TableCell>
                                                <pre>{kv.key}</pre>
                                            </TableCell>
                                            <TableCell>
                                                <pre>{kv.label}</pre>
                                            </TableCell>
                                            <TableCell>{kv.value}</TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </CardContent>
                    </Collapse>
                </>
            )}
        </Card>
    )
}
