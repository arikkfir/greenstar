import {DataGridPremium, DataGridPremiumProps, GridInitialState, useGridApiRef} from "@mui/x-data-grid-premium";
import {useCallback, useLayoutEffect, useState} from "react";
import {CircularProgress} from "@mui/material";
import {GridValidRowModel} from "@mui/x-data-grid-pro";

/*
 * To ensure DataGrid resizes correctly inside Flex containers (especially on height decrease!), make sure to set the
 * "overflow" CSS property to "hidden" ON ALL ITS PARENT HIERARCHY!
 *
 * See 1: https://github.com/mui/mui-x/issues/11295
 * See 2: https://github.com/mui/mui-x/issues/8463
 */

export interface CustomDataGridProps<R extends GridValidRowModel = any> extends Omit<DataGridPremiumProps<R>, 'apiRef'> {
    stateId: string
}

export function CustomDataGrid<R extends GridValidRowModel = any>(props: CustomDataGridProps) {
    const {stateId, sx, ...dataGridProps} = props;
    const [initialState, setInitialState] = useState<GridInitialState>();
    const apiRef = useGridApiRef();

    const prefixedStateID = 'dataGridState-' + stateId
    const saveSnapshot = useCallback(() => {
        if (apiRef?.current?.exportState && localStorage) {
            const currentState = apiRef.current.exportState();
            localStorage.setItem(prefixedStateID, JSON.stringify(currentState));
        }
    }, [apiRef]);

    useLayoutEffect(() => {
        const stateFromLocalStorage = localStorage?.getItem(prefixedStateID);
        setInitialState(stateFromLocalStorage ? JSON.parse(stateFromLocalStorage) : {});
        window.addEventListener('beforeunload', saveSnapshot);
        return () => {
            window.removeEventListener('beforeunload', saveSnapshot);
            saveSnapshot();
        };
    }, [saveSnapshot]);

    if (!initialState) {
        return <CircularProgress/>;
    }

    return (
        <DataGridPremium<R> sx={sx} apiRef={apiRef}  {...dataGridProps} initialState={initialState}/>
    )
}
