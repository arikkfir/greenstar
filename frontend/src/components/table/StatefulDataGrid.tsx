import {DataGridPremium, DataGridPremiumProps, GridInitialState, useGridApiRef} from "@mui/x-data-grid-premium";
import {useCallback, useLayoutEffect, useState} from "react";
import {CircularProgress} from "@mui/material";
import {GridValidRowModel} from "@mui/x-data-grid-pro";

export interface StatefulDataGridProps<R extends GridValidRowModel = any> extends Omit<DataGridPremiumProps<R>, 'apiRef'> {
    stateId: string
}

export function StatefulDataGrid<R extends GridValidRowModel = any>(props: StatefulDataGridProps) {
    const {stateId: givenStateId, ...dataGridProps} = props;
    const stateId = 'dataGridState-' + givenStateId

    const apiRef = useGridApiRef();

    const [initialState, setInitialState] = useState<GridInitialState>();

    const saveSnapshot = useCallback(() => {
        if (apiRef?.current?.exportState && localStorage) {
            const currentState = apiRef.current.exportState();
            localStorage.setItem(stateId, JSON.stringify(currentState));
        }
    }, [apiRef]);

    useLayoutEffect(() => {
        const stateFromLocalStorage = localStorage?.getItem(stateId);
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
        <DataGridPremium<R> apiRef={apiRef} {...dataGridProps} initialState={initialState}/>
    )
}