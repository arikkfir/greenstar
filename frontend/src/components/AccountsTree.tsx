import "./AccountsTree.scss"
import { forwardRef, Ref, SyntheticEvent, useCallback, useContext, useMemo } from "react"
import { useLazyQuery, useMutation, useQuery } from "@apollo/client"
import { Backdrop, CircularProgress, Snackbar, Typography } from "@mui/material"
import {
    RichTreeViewPro,
    RichTreeViewProProps,
    TreeItem,
    TreeItemProps,
    TreeItemSlotProps,
    TreeViewBaseItem,
    TreeViewItemId,
    useTreeItemModel,
    useTreeViewApiRef,
} from "@mui/x-tree-view-pro"
import { useTenantID } from "../hooks/tenant.ts"
import { SelectedCurrencyContext } from "../contexts/SelectedCurrency.ts"
import { Account, Currency } from "../graphql/graphql.ts"
import { useCurrencyFormatter } from "../hooks/locale.ts"
import { gql } from "../graphql"
import classNames from "classnames"

interface TreeViewItemReorderPosition {
    parentId: string | null;
    index: number;
}

const RootAccounts = gql(`
    query RootAccountsQuery($tenantID: ID!, $currency: String!, $until: DateTime) {
        tenant(id: $tenantID) {
            id
            rootAccounts {
                id
                label: displayName
                icon
                childCount
                balance(currency: $currency, until: $until)
            }
        }
    }
`)

const ChildAccounts = gql(`
    query ChildAccountsQuery($tenantID: ID!, $parentID: ID!, $currency: String!, $until: DateTime) {
        tenant(id: $tenantID) {
            id
            account(id: $parentID) {
                id
                children {
                    id
                    label: displayName
                    icon
                    childCount
                    balance(currency: $currency, until: $until)
                }
            }
        }
    }
`)

const MoveAccount = gql(`
    mutation MoveAccount($tenantID: ID!, $accountID: ID!, $targetParentAccountID: ID) {
        moveAccount(tenantID: $tenantID, accountID: $accountID, targetParentAccountID: $targetParentAccountID) {
            id
        }
    }
`)

export type AccountNode = TreeViewBaseItem<{
    id: Account["id"]
    label: Account["displayName"]
    icon: Account["icon"]
    balance: Account["balance"]
    childCount: Account["childCount"]
}>

interface AccountNodeLabelProps {
    label: string,
    balance?: number,
    currency?: Currency,
}

function AccountNodeLabel({ label, balance, currency }: AccountNodeLabelProps) {
    const currencyFormatter = useCurrencyFormatter()
    if (typeof balance !== "undefined" && currency) {
        const formattedBalance = currencyFormatter(balance, currency.code)
        return (
            <div className="account-label">
                <Typography>{label}</Typography>
                <Typography variant="overline">{formattedBalance}</Typography>
            </div>
        )
    } else {
        return (
            <div className="account-label">
                <Typography>{label}</Typography>
            </div>
        )
    }
}

interface AccountNodeTreeItemProps extends TreeItemProps {
    showBalance?: boolean
}

const AccountNodeTreeItem = forwardRef(function(
    props: AccountNodeTreeItemProps,
    ref: Ref<HTMLLIElement>,
) {
    const item                              = useTreeItemModel<AccountNode>(props.itemId)!
    const { showBalance, ...treeItemProps } = props
    const selectedCurrency                  = useContext(SelectedCurrencyContext)
    if (!selectedCurrency.currency) {
        return <></>
    }

    return (
        <TreeItem
            {...treeItemProps}
            ref={ref}
            slots={{ label: AccountNodeLabel }}
            slotProps={
                {
                    label: {
                        balance: showBalance ? item.balance : undefined,
                        currency: { code: selectedCurrency.currency.code },
                        label: item.label,
                    } as AccountNodeLabelProps,
                } as TreeItemSlotProps
            }
        />
    )
})

type TreeProps<Multiple extends boolean> = Omit<RichTreeViewProProps<AccountNode, Multiple>, "onItemSelectionToggle" | "items" | "dataSource" | "slots">

interface AccountsTreeProps<Multiple extends boolean> extends TreeProps<Multiple> {
    onAccountSelectionChange?: (accountId: AccountNode[]) => void
    showBalance?: boolean
}

export function AccountsTree<Multiple extends boolean>({
    showBalance,
    onAccountSelectionChange,
    className,
    ...rest
}: AccountsTreeProps<Multiple>) {
    const tenantID         = useTenantID()
    const selectedCurrency = useContext(SelectedCurrencyContext)
    const treeViewApiRef   = useTreeViewApiRef()

    // Initial query to load the root accounts
    const {
              data: rootAccounts,
              loading: rootAccountsLoading,
              error: errorLoadingRootAccounts,
          } = useQuery(RootAccounts, { variables: { tenantID, currency: selectedCurrency.currency?.code || "USD" } })

    // Mapping of the root accounts into tree nodes (refreshes whenever the root accounts change)
    const rootAccountNodes = useMemo(
        (): AccountNode[] => rootAccounts?.tenant?.rootAccounts || [],
        [ rootAccounts?.tenant?.rootAccounts ],
    )

    // Lazy GraphQL query to load child accounts of a given parent account
    const [ loadChildAccounts, loadChildAccountsResult ] = useLazyQuery(ChildAccounts)

    // Mutation for moving accounts in the account tree
    const [ moveAccount, { loading: movingAccount, error: accountMoveError } ] = useMutation(MoveAccount)

    // Callback passed to the TreeView to load child accounts for a given parent account
    const loadChildAccountNodes = useCallback(
        async (parentID?: TreeViewItemId): Promise<AccountNode[]> => {
            if (!parentID) {
                return []
            }

            const r = await loadChildAccounts({
                variables: {
                    tenantID,
                    parentID,
                    currency: selectedCurrency.currency?.code || "USD",
                },
            })
            if (r.error) {
                throw r.error
            }

            return r.data?.tenant?.account?.children || []
        },
        [ loadChildAccounts, tenantID, selectedCurrency.currency?.code ],
    )

    // Callback to handle account selection
    const handleAccountSelectionToggle = useCallback(
        (_event: SyntheticEvent | null, idOrIDs: any): void => {
            if (treeViewApiRef && treeViewApiRef.current && onAccountSelectionChange) {
                const ids      = Array.isArray(idOrIDs) ? idOrIDs as TreeViewItemId[] : [ idOrIDs as TreeViewItemId ]
                const accounts = ids.map(id => treeViewApiRef.current!.getItem(id)!)
                onAccountSelectionChange(accounts)
            }
        },
        [ onAccountSelectionChange, treeViewApiRef ],
    )

    // Callback to handle account dragging to new parents
    const handleAccountMove = useCallback(
        async (params: {
            itemId: string,
            oldPosition: TreeViewItemReorderPosition,
            newPosition: TreeViewItemReorderPosition
        }) => await moveAccount({
            variables: { tenantID, accountID: params.itemId, targetParentAccountID: params.newPosition.parentId },
        }),
        [ tenantID, moveAccount ],
    )

    // Data source for the tree view
    // noinspection JSUnusedGlobalSymbols
    const accountsDataSource = useMemo(
        () => ({
            getChildrenCount: (parentAccount: AccountNode): number => parentAccount.childCount,
            getTreeItems: loadChildAccountNodes,
        }),
        [ loadChildAccountNodes ],
    )

    if (rootAccountsLoading) {
        return (
            <Backdrop open={true}>
                <CircularProgress />
            </Backdrop>
        )
    }

    return (
        <>
            <Snackbar open={errorLoadingRootAccounts !== undefined} message={errorLoadingRootAccounts?.message} />
            <Snackbar open={loadChildAccountsResult.error !== undefined}
                      autoHideDuration={10_000}
                      message={loadChildAccountsResult.error?.message}
            />
            <Snackbar open={accountMoveError !== undefined}
                      autoHideDuration={10_000}
                      message={accountMoveError?.message}
            />
            <RichTreeViewPro<AccountNode, Multiple>
                {...rest}
                apiRef={treeViewApiRef}
                className={classNames("accounts-tree", className)}
                items={rootAccountNodes}
                slots={{ item: AccountNodeTreeItem }}
                slotProps={{
                    item: { showBalance } as AccountNodeTreeItemProps,
                }}
                itemsReordering
                canMoveItemToNewPosition={(_params) => !movingAccount}
                onItemPositionChange={handleAccountMove}
                onSelectedItemsChange={handleAccountSelectionToggle}
                dataSource={accountsDataSource}
            />
        </>
    )
}
