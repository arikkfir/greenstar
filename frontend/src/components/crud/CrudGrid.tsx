import {Alert, AlertProps, AlertTitle, Box, Button, Snackbar, Typography} from "@mui/material";
import {
    DataGridPremium,
    GridActionsCellItem,
    GridColDef,
    GridRowModes,
    GridRowModesModel,
    GridToolbarContainer,
    useGridRootProps
} from "@mui/x-data-grid-premium";
import {useCallback, useMemo, useState} from "react";
import {GridValidRowModel} from "@mui/x-data-grid/models/gridRows";
import {Add as AddIcon, Cancel as CancelIcon, Delete as DeleteIcon, Save as SaveIcon} from "@mui/icons-material";

export type Row = GridValidRowModel & { id: string, isNew: boolean }

interface NoRowsOverlayMessageProps {
    message?: string
}

function NoRowsOverlayMessage({message = "No rows found"}: NoRowsOverlayMessageProps) {
    return (<Box sx={{padding: 3, alignItems: "center"}}>
        <Typography>{message}</Typography>
    </Box>)
}

interface GridToolbarProps {
    createRow: () => void
    creationDisabled: boolean
}

function Toolbar({createRow, creationDisabled}: GridToolbarProps) {
    const rootProps = useGridRootProps();
    return (<GridToolbarContainer sx={{paddingLeft: 1}}>
        <Button disabled={creationDisabled}
                onClick={() => createRow()}
                startIcon={<AddIcon/>}
                size="small"
                {...rootProps.slotProps?.baseButton}>
            Create
        </Button>
    </GridToolbarContainer>);
}

export interface GridProps<R extends Row> {
    newRow: () => R
    defaultFieldToFocus: string
    columns: Array<GridColDef<R>>
    rows: Array<R> | undefined
    loading: boolean
    loadingError: Error | undefined
    createRow: (row: R) => Promise<R>
    creating: boolean
    resetCreation: () => void,
    updateRow: (row: R) => Promise<R>
    updating: boolean
    resetUpdate: () => void,
    deleteRow: (row: R) => Promise<any>
    deleting: boolean
    resetDeletion: () => void,
}

export function CrudGrid<R extends Row>(props: GridProps<R>) {
    const [message, setMessage] = useState<Pick<AlertProps, 'children' | 'severity'> | null>(null);
    const handleCloseMessageSnackbar = () => setMessage(null);
    const [overlayRows, setOverlayRows] = useState<R[]>([]);
    const [rowModesModel, setRowModesModel] = useState<GridRowModesModel>({});
    const {
        newRow,
        defaultFieldToFocus,
        columns,
        rows,
        createRow,
        resetCreation,
        updateRow,
        resetUpdate,
        deleteRow,
        resetDeletion
    } = props

    const removeOverlayRow = useCallback((row: R): R => {
        const {[row.id]: removedId, ...newRowModesModel} = rowModesModel
        setRowModesModel(newRowModesModel)
        setOverlayRows(overlayRows.filter(r => row.id !== r.id))
        return row
    }, [rowModesModel, overlayRows])

    const onDeleteClicked = useCallback(
        (row: R) => deleteRow(row)
            .then(() => removeOverlayRow(row))
            .catch(e => setMessage({children: e.message, severity: "error"}))
            .finally(resetDeletion),
        [removeOverlayRow, deleteRow, resetDeletion])

    const addNewRowToGrid = useCallback(() => {
        const r = newRow();
        setOverlayRows((oldRows) => [...oldRows, {...r, isNew: true}]);
        setRowModesModel((oldModel) => ({
            ...oldModel,
            [r.id]: {mode: GridRowModes.Edit, fieldToFocus: defaultFieldToFocus},
        }));
    }, [newRow, setOverlayRows, setRowModesModel, defaultFieldToFocus])

    const saveNewRow = useCallback((row: R): Promise<R> => (
        createRow(row)
            .then(row => ({...row, isNew: false}))
            .then(removeOverlayRow)
            .catch(e => {
                setMessage({children: e.message, severity: "error"})
                return row
            })
            .finally(resetCreation)
    ), [createRow, removeOverlayRow, resetCreation])

    const saveUpdatedRow = useCallback((row: R) => (
        updateRow(row)
            .then(row => ({...row, isNew: false}))
            .then(removeOverlayRow)
            .catch(e => {
                setMessage({children: e.message, severity: "error"})
                return row
            })
            .finally(resetUpdate)
    ), [updateRow, removeOverlayRow, resetUpdate])

    const handleRowUpdate = useCallback((row: R) => {
        if (row.isNew) {
            return saveNewRow(row)
        } else {
            return saveUpdatedRow(row)
        }
    }, [saveNewRow, saveUpdatedRow])

    const columnsWithActions = useMemo(() => {
        return columns.concat(
            [
                {
                    field: 'actions',
                    type: 'actions',
                    getActions: ({id, row}) => {
                        const mode: GridRowModes = rowModesModel[id]?.mode || GridRowModes.View
                        return [
                            <GridActionsCellItem icon={<SaveIcon/>}
                                                 disabled={mode !== GridRowModes.Edit}
                                                 label="Save"
                                                 onClick={() => removeOverlayRow(row)}/>,
                            <GridActionsCellItem icon={<CancelIcon/>}
                                                 disabled={mode !== GridRowModes.Edit}
                                                 label="Cancel"
                                                 onClick={() => removeOverlayRow(row)}/>,
                            <GridActionsCellItem icon={<DeleteIcon/>}
                                                 disabled={mode !== GridRowModes.View}
                                                 label="Delete"
                                                 onClick={() => onDeleteClicked(row)}/>,
                        ]
                    },
                },
            ]
        )
    }, [rowModesModel, removeOverlayRow, onDeleteClicked, columns])

    const mergedRows = useMemo(
        () => (rows?.filter((sr) => !overlayRows.find((or) => or.id === sr.id)).concat(overlayRows) || []),
        [rows, overlayRows],
    )

    const disableRowCreation =
        mergedRows.some(r => r.isNew) ||
        Object.values(rowModesModel).some(m => m.mode === GridRowModes.Edit)

    const loading = props.loading || props.creating || props.updating || props.deleting;
    return (
        <Box>
            <DataGridPremium<R> autoHeight={true}
                                rowSelection={false}
                                loading={loading}
                                columns={columnsWithActions}
                                rows={mergedRows}
                                editMode="row"
                                rowModesModel={rowModesModel}
                                onRowModesModelChange={model => setRowModesModel(model)}
                                processRowUpdate={(newRow, _) => handleRowUpdate(newRow)}
                                slots={{
                                    noRowsOverlay: NoRowsOverlayMessage,
                                    toolbar: Toolbar,
                                }}
                                slotProps={{
                                    noRowsOverlay: {
                                        message: props.loadingError?.message || "No tenants found"
                                    },
                                    toolbar: {
                                        creationDisabled: disableRowCreation || loading,
                                        createRow: addNewRowToGrid,
                                    },
                                }}
            />
            {!!message && (<Snackbar
                open
                anchorOrigin={{
                    vertical: 'top',
                    horizontal: 'center'
                }}
                onClose={handleCloseMessageSnackbar}
                autoHideDuration={1000 * 10}>
                <Alert severity={message.severity} onClose={handleCloseMessageSnackbar}>
                    {!!message.severity && (
                        <AlertTitle>{message.severity[0].toUpperCase() + message.severity.substring(1)}</AlertTitle>
                    )}
                    {message.children}
                </Alert>
            </Snackbar>)}
        </Box>
    )
}
