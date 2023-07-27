import {CrudGrid} from "../components/crud/CrudGrid";
import {useMemo} from "react";
import {
    allTenantsQuery,
    createTenantMutation,
    deleteTenantMutation,
    randomTenantID,
    updateTenantMutation
} from "../services/tenants";
import {useMutation, useQuery} from "@apollo/client";
import {Scalars} from "../gql/graphql";

const idColumn = {
    field: 'id',
    editable: false,
    headerName: 'ID',
    description: 'Tenant ID',
    type: 'string'
}

const displayNameColumn = {
    field: 'displayName',
    headerName: 'Name',
    description: 'Tenant name',
    editable: true
}

export function Tenants() {
    const columns = useMemo(() => ([
        idColumn,
        displayNameColumn,
    ]), [])

    const {data, loading, error: loadingError} = useQuery(allTenantsQuery);
    const [createTenant, {loading: creatingRow, reset: resetCreation}] = useMutation(createTenantMutation);
    const [updateTenant, {loading: updatingRow, reset: resetUpdate}] = useMutation(updateTenantMutation);
    const [deleteTenant, {loading: deletingRow, reset: resetDeletion}] = useMutation(deleteTenantMutation);

    type TenantRow = {
        isNew: boolean,
        id: Scalars['ID']['output'];
        displayName: Scalars['String']['output'];
    }

    const createRow = ({id, displayName}: TenantRow) =>
        createTenant({refetchQueries: [allTenantsQuery], variables: {id, displayName}})
            .then(r => ({
                id: r.data?.createTenant.id || "",
                displayName: r.data?.createTenant.displayName || "",
                isNew: false,
            }))
    const updateRow = ({id, displayName}: TenantRow) =>
        updateTenant({refetchQueries: [allTenantsQuery], variables: {id, displayName}})
            .then(r => ({
                id: r.data?.updateTenant.id || "",
                displayName: r.data?.updateTenant.displayName || "",
                isNew: false,
            }))
    const deleteRow = ({id}: TenantRow) => deleteTenant({refetchQueries: [allTenantsQuery], variables: {id}})

    return (
        <CrudGrid<TenantRow> newRow={() => ({id: randomTenantID(7), displayName: '', isNew: false})}
                             defaultFieldToFocus='displayName'
                             columns={columns}
                             rows={data?.tenants.map(r => ({...r, isNew: false})) || []}
                             loading={loading}
                             loadingError={loadingError}
                             createRow={createRow}
                             creating={creatingRow}
                             resetCreation={resetCreation}
                             updateRow={updateRow}
                             updating={updatingRow}
                             resetUpdate={resetUpdate}
                             deleteRow={deleteRow}
                             deleting={deletingRow}
                             resetDeletion={resetDeletion}/>
    )
}
