import {Tenant} from "../gql/graphql.ts";
import {useMutation, useQuery} from "@apollo/client";
import {graphql} from "../gql";
import {CrudGrid} from "../components/crud.tsx";

const tenantIDHashLetters = "abcdefghijklmnopqrstuvwxyz"

function randomTenantID(length: number): string {
    let hash = ""
    for (let i = 0; i < length; i++) {
        hash += tenantIDHashLetters.charAt(Math.floor(Math.random() * tenantIDHashLetters.length))
    }
    return hash
}

const allTenantsQuery = graphql(/* GraphQL */ `
    query allTenants {
        tenants {
            id
            displayName
        }
    }
`)

const createTenantMutation = graphql(/* GraphQL */`
    mutation CreateTenant($id: String, $displayName: String!) {
        tenant: createTenant(tenantID: $id, tenant: {displayName: $displayName}) {
            id
            displayName
        }
    }
`)

const deleteTenantMutation = graphql(/* GraphQL */`
    mutation DeleteTenant($id: String!) {
        id: deleteTenant(tenantID: $id)
    }
`)

const updateTenantMutation = graphql(/* GraphQL */`
    mutation UpdateTenant($id: String!, $displayName: String!) {
        tenant: updateTenant(tenantID: $id, tenant: {displayName: $displayName}) {
            id
            displayName
        }
    }
`)

export function Tenants() {
    const {data: tenants, loading: loadingTenants, error: tenantsLoadingError} = useQuery(allTenantsQuery);
    const [createTenant, {loading: creatingTenant}] = useMutation(createTenantMutation, {refetchQueries: [allTenantsQuery]});
    const [updateTenant, {loading: updatingTenant}] = useMutation(updateTenantMutation, {refetchQueries: [allTenantsQuery]});
    const [deleteTenant, {loading: deletingTenant}] = useMutation(deleteTenantMutation, {refetchQueries: [allTenantsQuery]});
    const columns = [
        {field: 'id', editable: false, headerName: 'ID', type: 'string'},
        {field: 'displayName', headerName: 'Name', editable: true},
    ]
    type Row = {
        id: Tenant['id'],
        displayName: Tenant['displayName']
    }
    return (
        <CrudGrid<Row> columns={columns}
                       defaultFieldToFocus="displayName"
                       query={{
                           rows: tenants?.tenants || [],
                           loading: loadingTenants,
                           error: tenantsLoadingError,
                       }}
                       creation={{
                           newRow: () => ({id: randomTenantID(7), displayName: ''}),
                           create: async ({id, displayName}) => {
                               const result = await createTenant({variables: {id, displayName}})
                               if (!result.data) {
                                   throw new Error("Empty result returned from server")
                               } else {
                                   return result.data.tenant
                               }
                           },
                           creating: creatingTenant,
                       }}
                       update={{
                           update: async ({id, displayName}) => {
                               const result = await updateTenant({variables: {id, displayName}})
                               if (!result.data) {
                                   throw new Error("Empty result returned from server")
                               } else {
                                   return result.data.tenant
                               }
                           },
                           updating: updatingTenant,
                       }}
                       deletion={{
                           delete: async (id) => await deleteTenant({variables: {id}}),
                           deleting: deletingTenant,
                       }}
        />
    )
}
