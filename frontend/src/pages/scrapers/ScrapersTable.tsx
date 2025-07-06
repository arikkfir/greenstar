import { useMutation, useQuery } from "@apollo/client"
import { useCallback, useMemo, useRef, useState } from "react"
import {
    ColumnsPanelTrigger,
    DataGridPremium,
    GridColDef,
    GridRenderCellParams,
    Toolbar,
    ToolbarButton,
} from "@mui/x-data-grid-premium"
import { gql } from "../../graphql"
import { useTenantID } from "../../hooks/tenant.ts"
import { Drawer, IconButton, Snackbar, Tooltip } from "@mui/material"
import DeleteIcon from "@mui/icons-material/Delete"
import PlayArrowIcon from "@mui/icons-material/PlayArrow"
import EditIcon from "@mui/icons-material/Edit"
import { ScraperDrawer } from "./ScraperDrawer.tsx"
import { FetchScrapers, ScraperRow } from "./FetchScrapersQuery.ts"
import AddIcon from "@mui/icons-material/Add"
import ViewColumnIcon from "@mui/icons-material/ViewColumn"

const TriggerScraper = gql(`
    mutation TriggerScraper($tenantID: ID!, $scraperID: ID!) {
        triggerScraper(tenantID: $tenantID, id: $scraperID) {
            id
            createdAt
            status
            parameters {
                parameter {
                    id
                    displayName
                    type
                }
                value
            }
        }
    }
`)

const DeleteScraper = gql(`
    mutation DeleteScraper($tenantID: ID!, $id: ID!) {
        deleteScraper(tenantID: $tenantID, id: $id)
    }
`)

export function ScrapersTable() {
    const tenantID: string                            = useTenantID()
    const [ snackbarMessage, setSnackbarMessage ]     = useState<string | null>(null)
    const [ scraperForDrawer, setScraperForDrawer ]   = useState<ScraperRow | undefined>()
    const [ scraperDrawerOpen, setScraperDrawerOpen ] = useState(false)

    const { data: scrapers, loading: loadingScrapers, error: errorLoadingScrapers, refetch: refetchScrapers } =
              useQuery(FetchScrapers, { variables: { tenantID } })

    const [ deleteScraper, { loading: deleting } ] = useMutation(DeleteScraper, {
        onCompleted: async () => {
            setSnackbarMessage("Scraper deleted successfully")
            await refetchScrapers()
        },
        onError: (error) => {
            setSnackbarMessage(`Error deleting scraper: ${error.message}`)
        },
    })

    const [ triggerScraper, { loading: triggering } ] = useMutation(TriggerScraper, {
        onCompleted: async () => {
            // TODO: show link to the scraper run page (which will show the logs) in the snackbar message
            setSnackbarMessage("Scraper triggerred successfully! Click here to track it")
            await refetchScrapers()
        },
        onError: (error) => {
            setSnackbarMessage(`Error triggering a run for the scraper: ${error.message}`)
        },
    })

    const handleEdit = useCallback(async (row: ScraperRow) => {
        setScraperForDrawer(row)
        setScraperDrawerOpen(true)
    }, [ setScraperForDrawer, setScraperDrawerOpen ])

    const handleScraperUpserted = async (_: any) => {
        await refetchScrapers()
    }

    const handleAddScraper = useCallback(() => {
        setScraperForDrawer(undefined) // undefined means creating a new scraper
        setScraperDrawerOpen(true)
    }, [ setScraperForDrawer, setScraperDrawerOpen ])


    const handleScraperDrawerClosed = useCallback(async () => {
        setScraperForDrawer(undefined)
        setScraperDrawerOpen(false)
    }, [ setScraperForDrawer, setScraperDrawerOpen ])

    const handleDelete = useCallback(async (id: string) => {
        await deleteScraper({
            variables: { tenantID, id },
        })
    }, [ tenantID, deleteScraper ])

    const handleRun = useCallback(async (id: string) => {
        await triggerScraper({
            variables: {
                tenantID,
                scraperID: id,
            },
        })
    }, [ tenantID, triggerScraper ])

    const renderActionsCell = useCallback(
        (p: GridRenderCellParams<ScraperRow>) => {
            return (
                <div className="actions-cell">
                    <Tooltip title="Edit scraper">
                        <IconButton size="small" onClick={() => handleEdit(p.row)} disabled={triggering || deleting}>
                            <EditIcon fontSize="small" />
                        </IconButton>
                    </Tooltip>
                    <Tooltip title="Delete scraper">
                        <IconButton size="small" onClick={() => handleDelete(p.row.id)}
                                    disabled={triggering || deleting}>
                            <DeleteIcon fontSize="small" />
                        </IconButton>
                    </Tooltip>
                    <Tooltip title="Run scraper">
                        <IconButton size="small" onClick={() => handleRun(p.row.id)} disabled={triggering || deleting}>
                            <PlayArrowIcon fontSize="small" />
                        </IconButton>
                    </Tooltip>
                </div>
            )
        },
        [ handleEdit, handleDelete, handleRun, deleting, triggering ],
    )

    const scraperColumns: GridColDef<ScraperRow>[] = useMemo(
        (): GridColDef<ScraperRow>[] => [
            {
                field: "id",
                headerName: "ID",
                flex: 1,
            },
            {
                field: "displayName",
                headerName: "Name",
                flex: 1,
            },
            {
                field: "type",
                headerName: "Type",
                flex: 1,
                valueGetter: (value: { displayName: string }) => value.displayName,
            },
            {
                field: "createdAt",
                headerName: "Created",
                flex: 1,
            },
            {
                field: "updatedAt",
                headerName: "Updated",
                flex: 1,
            },
            {
                field: "actions",
                headerName: "Actions",
                flex: 1,
                sortable: false,
                filterable: false,
                renderCell: renderActionsCell,
                groupable: false,
                pivotable: false,
                editable: false,
                aggregable: false,
                pinnable: false,
            },
        ],
        [ renderActionsCell ],
    )

    const newScraperPanelTriggerRef = useRef<HTMLButtonElement>(null)
    const toolbar                   = () => (
        <Toolbar>
            <Tooltip title="New scraper">
                <ToolbarButton ref={newScraperPanelTriggerRef} onClick={handleAddScraper}>
                    <AddIcon fontSize="small" />
                </ToolbarButton>
            </Tooltip>
            <Tooltip title="Columns">
                <ColumnsPanelTrigger render={<ToolbarButton />}>
                    <ViewColumnIcon fontSize="small" />
                </ColumnsPanelTrigger>
            </Tooltip>
        </Toolbar>
    )

    return (
        <div className="scrapers-grid-content">
            <Snackbar open={!!errorLoadingScrapers} message={errorLoadingScrapers?.message} />
            <Snackbar
                open={snackbarMessage !== null}
                autoHideDuration={6000}
                onClose={() => setSnackbarMessage(null)}
                message={snackbarMessage}
            />
            <Drawer anchor="right" open={scraperDrawerOpen} onClose={handleScraperDrawerClosed}>
                <ScraperDrawer scraper={scraperForDrawer}
                               onScraperUpserted={handleScraperUpserted}
                               onClose={handleScraperDrawerClosed} />
            </Drawer>
            <DataGridPremium<ScraperRow>
                rows={scrapers?.tenant?.scrapers || []}
                showToolbar
                columns={scraperColumns}
                disablePivoting
                disableRowGrouping
                disableAggregation
                loading={loadingScrapers}
                pagination
                pageSizeOptions={[ 10, 25, 50, 100 ]}
                slots={{ toolbar: toolbar }}
                initialState={{
                    columns: {
                        columnVisibilityModel: {
                            id: false,
                        },
                    },
                    pagination: {
                        paginationModel: { pageSize: 25, page: 0 },
                    },
                }}
            />
        </div>
    )
}
