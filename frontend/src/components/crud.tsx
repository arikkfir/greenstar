import {
    DataGridPremium,
    GridActionsCellItem,
    GridRowModes,
    GridRowModesModel,
    GridRowParams,
    GridToolbarContainer,
    GridToolbarProps,
    GridValidRowModel,
    useGridRootProps
} from "@mui/x-data-grid-premium";
import {useSnackbar} from "notistack";
import {DataGridPremiumProps} from "@mui/x-data-grid-premium/models/dataGridPremiumProps";
import {GridColDef} from "@mui/x-data-grid/models/colDef/gridColDef";
import {useCallback, useMemo, useState} from "react";
import SaveIcon from "@mui/icons-material/Save";
import CancelIcon from "@mui/icons-material/Cancel";
import DeleteIcon from "@mui/icons-material/Delete";
import {GridRowModesModelProps} from "@mui/x-data-grid/models/api/gridEditingApi";
import {Box, Button, Typography} from "@mui/material";
import {ToolbarPropsOverrides} from "@mui/x-data-grid/models/gridSlotsComponentsProps";
import AddIcon from "@mui/icons-material/Add";

export type Row = GridValidRowModel & {
    id: string | number
}

type GridRow<R> = R & Row & {
    isNew: boolean
}

interface NoRowsOverlayMessageProps {
    message?: string
}

function NoRowsOverlayMessage({message = "No rows found"}: NoRowsOverlayMessageProps) {
    return <Box sx={{padding: 3, alignItems: "center"}}><Typography>{message}</Typography></Box>
}

type ToolbarProps = Partial<GridToolbarProps | ToolbarPropsOverrides>

function Toolbar({onCreate, creationDisabled}: ToolbarProps) {
    const rootProps = useGridRootProps();
    return (
        <GridToolbarContainer sx={{paddingLeft: 1}}>
            <Button disabled={creationDisabled}
                    onClick={onCreate}
                    startIcon={<AddIcon/>}
                    size="small"
                    {...rootProps.slotProps?.baseButton}>
                Create
            </Button>
        </GridToolbarContainer>
    );
}

type ExposedDataGridProps<R extends Row> = Omit<DataGridPremiumProps<R>,
    "columns" | "editMode" | "loading" | "onRowModesModelChange" | "processRowUpdate" | "rows" | "rowModesModel" |
    "rowSelection" | "slots" | "slotProps">

interface CrudGridProps<R extends Row> extends ExposedDataGridProps<GridRow<R>> {
    columns: GridColDef<R>[]
    defaultFieldToFocus?: string
    query: {
        rows: R[]
        loading: boolean
        error: Error | undefined
    }
    creation: {
        newRow: () => R
        create: (row: R) => Promise<R>
        creating: boolean
    }
    update: {
        update: (row: R) => Promise<R>
        updating: boolean
    }
    deletion: {
        delete: (id: R['id']) => Promise<any>
        deleting: boolean
    }
}

export function CrudGrid<R extends Row>(props: CrudGridProps<R>) {
    const {enqueueSnackbar} = useSnackbar();
    const {
        columns,
        defaultFieldToFocus,
        query,
        creation,
        update,
        deletion,
        ...rest
    } = props
    const [rowModesModel, setRowModesModel] = useState<GridRowModesModel>({});
    const [overlayRows, setOverlayRows] = useState<GridRow<R>[]>([]);
    const mergedRows = useMemo(
        (): GridRow<R>[] =>
            query.rows
                .filter(r => !overlayRows.find(or => or.id === r.id))
                .map(r => Object.assign({}, r, {isNew: false}))
                .concat(overlayRows),
        [query.rows, overlayRows],
    )
    const removeOverlayRow = useCallback((id: R['id']): void => {
        setRowModesModel(Object.assign({}, rowModesModel, {id: {mode: GridRowModes.View}}))
        setOverlayRows(overlayRows.filter(r => id !== r.id))
    }, [rowModesModel, overlayRows])
    const onRowUpdated = useCallback(async (newRow: GridRow<R>, oldRow: GridRow<R>): Promise<GridRow<R>> => {
        if (newRow.isNew) {
            try {
                let row = await creation.create(newRow);
                removeOverlayRow(row.id)
                return Object.assign({}, row, {isNew: false})
            } catch (e) {
                console.error("Creation failed: ", e)
                enqueueSnackbar(e instanceof Error ? e.message : "An internal error occurred.", {variant: "error"})
                return oldRow
            }
        } else {
            try {
                let row = await update.update(newRow);
                removeOverlayRow(row.id)
                return Object.assign({}, row, {isNew: false})
            } catch (e) {
                console.error("Update failed: ", e)
                enqueueSnackbar(e instanceof Error ? e.message : "An internal error occurred.", {variant: "error"})
                return oldRow
            }
        }
    }, [enqueueSnackbar, creation.create, update.update, removeOverlayRow])
    const onDeleteClicked = useCallback(async (id: R['id']) => {
        try {
            await deletion.delete(id);
            removeOverlayRow(id)
        } catch (e) {
            console.error("Deletion failed: ", e)
            enqueueSnackbar(e instanceof Error ? e.message : "An internal error occurred.", {variant: "error"})
        }
    }, [deletion.delete, enqueueSnackbar, removeOverlayRow])
    const rowActions = useCallback((id: R['id']) => {
        const mode: GridRowModes = rowModesModel[id]?.mode || GridRowModes.View
        return [
            <GridActionsCellItem icon={<SaveIcon/>}
                                 disabled={mode !== GridRowModes.Edit}
                                 label="Save"
                                 onClick={() => removeOverlayRow(id)}/>,
            <GridActionsCellItem icon={<CancelIcon/>}
                                 disabled={mode !== GridRowModes.Edit}
                                 label="Cancel"
                                 onClick={() => removeOverlayRow(id)}/>,
            <GridActionsCellItem icon={<DeleteIcon/>}
                                 disabled={mode !== GridRowModes.View}
                                 label="Delete"
                                 onClick={() => onDeleteClicked(id)}/>,
        ]
    }, [rowModesModel, removeOverlayRow, onDeleteClicked])
    const columnsWithActions = useMemo((): GridColDef<GridRow<R>>[] => {
        return [
            ...columns as GridColDef<GridRow<R>>[],
            {field: 'actions', type: 'actions', getActions: ({id}: GridRowParams<R>) => rowActions(id as R['id'])},
        ]
    }, [columns])
    const loading = query.loading || creation.creating || update.updating || deletion.deleting;
    return (
        <DataGridPremium<GridRow<R>> columns={columnsWithActions}
                                     editMode="row"
                                     loading={loading}
                                     onRowModesModelChange={setRowModesModel}
                                     processRowUpdate={onRowUpdated}
                                     rows={mergedRows}
                                     rowModesModel={rowModesModel}
                                     rowSelection={false}
                                     slots={{
                                         noRowsOverlay: NoRowsOverlayMessage,
                                         toolbar: Toolbar,
                                     }}
                                     slotProps={{
                                         noRowsOverlay: {
                                             message: query.error?.message || "No rows found",
                                         },
                                         toolbar: {
                                             creationDisabled: loading || overlayRows.length > 0,
                                             onCreate: (): void => {
                                                 const row = Object.assign(creation.newRow(), {isNew: true})
                                                 const rowModel: GridRowModesModelProps = {mode: GridRowModes.Edit, fieldToFocus: defaultFieldToFocus};
                                                 setOverlayRows([...overlayRows, row])
                                                 setRowModesModel(Object.assign({}, rowModesModel, {[row.id]: rowModel}))
                                             },
                                         },
                                     }}
                                     {...rest}/>
    )
}
