import {
    Account,
    FetchScrapersQuery,
    ScraperParameterType,
    ScraperType,
    UpsertScraperMutation,
} from "../../graphql/graphql.ts"
import {
    Alert,
    Button,
    Checkbox,
    Divider,
    FormControl,
    FormControlLabel,
    InputLabel,
    MenuItem,
    Select,
    SelectChangeEvent,
    TextField,
    Typography,
} from "@mui/material"
import { useTenantID } from "../../hooks/tenant.ts"
import { useMutation, useQuery } from "@apollo/client"
import { FormEvent, useState } from "react"
import "./ScraperDrawer.scss"
import { FetchScraperTypes } from "./FetchScraperTypesQuery.ts"
import { FetchAccounts, ScraperRow } from "./FetchScrapersQuery.ts"
import { UpsertScraper } from "./UpsertScraperQuery.ts"

export interface ScraperDrawerProps {
    scraper?: ScraperRow,
    onScraperUpserted: (scraper: UpsertScraperMutation["upsertScraper"]) => void
    onClose: () => void
}

type ScraperParameter = NonNullable<FetchScrapersQuery["tenant"]> [ "scrapers" ][number]["parameters"][number]
//type ScraperTypeParameter = NonNullable<FetchScrapersQuery["tenant"]> [ "scrapers"
// ][number]["type"]["parameters"][number]

export function ScraperDrawer({ scraper, onScraperUpserted, onClose }: ScraperDrawerProps) {
    const tenantID = useTenantID()

    const { data: scraperTypes, error: errorLoadingScraperTypes, loading: loadingScraperTypes } =
              useQuery(FetchScraperTypes)

    const { data: accounts, error: errorLoadingAccounts, loading: loadingAccounts } =
              useQuery(FetchAccounts, { variables: { tenantID } })

    const [ upsertScraper, { loading: upserting, error: errorUpserting } ] = useMutation(UpsertScraper)

    const [ scraperTypeID, setScraperTypeID ] = useState<string | undefined>(scraper?.type?.id)
    const [ displayName, setDisplayName ]     = useState<string>(scraper?.displayName || "")
    const [ parameters, setParameters ]       = useState<ScraperParameter[]>(scraper?.parameters || [])

    const handleScraperTypeChange = (e: SelectChangeEvent<ScraperType["id"]>) => {
        const scraperTypeID = e.target.value
        setScraperTypeID(scraperTypeID)
        const st = scraperTypes?.scraperTypes?.find((st) => st.id == scraperTypeID)
        if (!st) {
            throw new Error(`Scraper type '${scraperTypeID}' not found`)
        }

        const newParameters: ScraperParameter[] = []
        for (let stp of st.parameters) {
            const pv = parameters.find(p => p.parameter.id === stp.id)
            if (pv) {
                newParameters.push(pv)
            } else {
                newParameters.push({ parameter: stp, value: "" })
            }
        }
        setParameters(newParameters)
    }

    function handleValueChange(p: ScraperParameter, value: string): void {
        setParameters(
            (prev) =>
                (
                    prev.map(old => old.parameter.id === p.parameter.id
                        ? { ...old, value }
                        : old)
                ),
        )
    }

    const handleSubmit = async (event: FormEvent) => {
        event.preventDefault()
        if (!loadingScraperTypes && !loadingAccounts && !upserting && displayName && scraperTypeID && !parameters.some(p => !p.value)) {
            const result = await upsertScraper({
                variables: {
                    tenantID,
                    scraperID: scraper?.id,
                    displayName,
                    scraperTypeID,
                    parameters: parameters.map(p => ({ parameterID: p.parameter.id, value: p.value })),
                },
            })
            if (result.data?.upsertScraper?.id) {
                onScraperUpserted(result.data.upsertScraper)
            }
            onClose()
        }
    }

    return (
        <article className="scraper-drawer">
            <Typography variant="h2">Scraper details</Typography>
            <Divider />
            {errorUpserting && (<Alert severity="error">{errorUpserting.message}</Alert>)}
            {errorLoadingScraperTypes && (<Alert severity="error">{errorLoadingScraperTypes.message}</Alert>)}
            {errorLoadingAccounts && (<Alert severity="error">{errorLoadingAccounts.message}</Alert>)}
            <form onSubmit={handleSubmit}>
                <FormControl fullWidth>
                    <InputLabel id="scraperType-label">Type</InputLabel>
                    <Select<ScraperType["id"]>
                        labelId="scraperType-label"
                        id="scraperType"
                        disabled={loadingScraperTypes}
                        value={scraperTypeID || ""}
                        label="Type"
                        onChange={(e) => handleScraperTypeChange(e)}
                    >
                        {scraperTypes?.scraperTypes?.map(type => (
                            <MenuItem value={type.id}>{type.displayName}</MenuItem>
                        ))}
                    </Select>
                </FormControl>
                <TextField
                    error={!displayName}
                    label="Display name"
                    required
                    disabled={loadingScraperTypes}
                    fullWidth
                    value={displayName}
                    onChange={(event) => setDisplayName(event.target.value)}
                    helperText={!displayName ? "Required" : ""}
                />
                {parameters.map((p) => {
                    switch (p.parameter.type) {
                        case ScraperParameterType.Account:
                            return (
                                <FormControl fullWidth>
                                    <InputLabel
                                        id={p.parameter.id + "-account-label"}>{p.parameter.displayName}</InputLabel>
                                    <Select<Account["id"]>
                                        labelId={p.parameter.id + "-account-label"}
                                        disabled={loadingAccounts}
                                        value={p.value}
                                        label={p.parameter.displayName}
                                        onChange={(event) => handleValueChange(p, event.target.value)}
                                    >
                                        {accounts?.tenant?.accounts?.map(account => (
                                            <MenuItem value={account.id}>{account.displayName}</MenuItem>
                                        ))}
                                    </Select>
                                </FormControl>
                            )
                        case ScraperParameterType.Boolean:
                            const flippedValue = p.value === "true" ? "false" : "true"
                            return <FormControlLabel
                                control={<Checkbox checked={p.value === "true"}
                                                   onChange={() => handleValueChange(p, flippedValue)} />}
                                label={p.parameter.displayName} />
                        case ScraperParameterType.Date:
                            throw new Error("Account parameter type not implemented")
                        case ScraperParameterType.Float:
                            return <TextField
                                required
                                value={p.value}
                                error={isNaN(parseFloat(p.value))}
                                label={p.parameter.displayName}
                                onChange={(event) => handleValueChange(p, event.target.value)}
                                helperText={isNaN(parseFloat(p.value)) ? "Invalid number" : ""}
                            />
                        case ScraperParameterType.Integer:
                            return <TextField
                                required
                                value={p.value}
                                error={isNaN(parseInt(p.value))}
                                label={p.parameter.displayName}
                                onChange={(event) => handleValueChange(p, event.target.value)}
                                helperText={isNaN(parseInt(p.value)) ? "Invalid number" : ""}
                            />
                        case ScraperParameterType.Password:
                            return <TextField
                                required
                                value={p.value}
                                type="password"
                                error={!p.value}
                                label={p.parameter.displayName}
                                onChange={(event) => handleValueChange(p, event.target.value)}
                                helperText={!p.value ? "Value must not be empty" : ""}
                            />
                        case ScraperParameterType.String:
                            return <TextField
                                required
                                value={p.value}
                                error={!p.value}
                                label={p.parameter.displayName}
                                onChange={(event) => handleValueChange(p, event.target.value)}
                                helperText={!p.value ? "Value must not be empty" : ""}
                            />
                        default:
                            throw new Error(`Unsupported parameter type '${p.parameter.type}'`)
                    }
                })}
                <Button type="submit"
                        loading={upserting}
                        variant="contained"
                        fullWidth
                        disabled={loadingScraperTypes || !displayName || !scraperTypeID}>Save</Button>
            </form>
        </article>
    )
}
