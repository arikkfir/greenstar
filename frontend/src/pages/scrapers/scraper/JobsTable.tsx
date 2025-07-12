import { useCallback, useMemo } from "react"
import { enqueueSnackbar } from "notistack"
import { dismissAction } from "../../../util/notistack-actions.tsx"
import { Scraper, Tenant } from "../../../graphql/graphql.ts"
import { ScraperJobRow, ScraperJobs, TriggerScraper } from "../ScrapersQueries.ts"
import { useMutation, useQuery } from "@apollo/client"
import { useTenantID } from "../../../hooks/tenant.ts"
import { Action, DataGrid } from "../../../components/DataGrid.tsx"
import { GridColDef } from "@mui/x-data-grid-premium"
import { temporalColumn } from "../../../util/datagrid.tsx"
import PlayArrowIcon from "@mui/icons-material/PlayArrow"
import EditIcon from "@mui/icons-material/Edit"
import { useLocation } from "wouter"

export interface JobsTableProps {
    scraperID: Scraper["id"]
}

export function JobsTable({ scraperID }: JobsTableProps) {
    const [ location, navigate ] = useLocation()

    const tenantID: Tenant["id"] = useTenantID()

    const [ triggerScraper, { loading: triggering } ]      = useMutation(TriggerScraper)
    const { data: jobs, error: errorLoadingJobs, loading } = useQuery(ScraperJobs, {
        pollInterval: location == "/" ? 1000 : 0,
        variables: { tenantID, scraperID },
    })

    const handleTrigger = useCallback(
        async () => {
            const result = await triggerScraper({ variables: { tenantID, scraperID } })
            if (!result.errors?.length && result.data?.triggerScraper) {
                // TODO: put link to job logs
                enqueueSnackbar(`New scraper job triggered <a>abc</a>`, { action: dismissAction })
            }
        },
        [ triggerScraper, tenantID, scraperID ],
    )

    const columns: GridColDef<ScraperJobRow>[] = useMemo(
        (): GridColDef<ScraperJobRow>[] => [
            { field: "id", headerName: "ID", flex: 1 },
            temporalColumn({ field: "createdAt", headerName: "Created" }),
            { field: "status", headerName: "Status", flex: 1 },
        ],
        [],
    )

    const actions: Action<ScraperJobRow>[] = useMemo(
        (): Action<ScraperJobRow>[] => [
            {
                key: "open",
                title: "Open",
                handler: row => navigate(`/jobs/${row.id}`),
                disabled: loading,
                icon: <EditIcon fontSize="small" />,
            },
        ],
        [ loading, navigate ],
    )

    return (
        <DataGrid<ScraperJobRow> columns={columns}
                                 newAction={handleTrigger}
                                 newActionIcon={<PlayArrowIcon fontSize="small" />}
                                 rowActions={actions}
                                 rows={jobs?.tenant?.scraper?.jobs || []}
                                 loading={loading || triggering}
                                 errorLoading={errorLoadingJobs}
                                 initialState={{
                                     columns: {
                                         columnVisibilityModel: {
                                             id: false,
                                         },
                                     },
                                 }} />
    )
}
