import {
    TreeItem2Checkbox,
    TreeItem2Content,
    TreeItem2GroupTransition,
    TreeItem2IconContainer,
    TreeItem2Label,
    TreeItem2Root
} from "@mui/x-tree-view/TreeItem2";
import {UseTreeItem2Parameters} from "@mui/x-tree-view/useTreeItem2";
import {forwardRef, HTMLAttributes, Ref, useContext} from "react";
import {useTreeItem2} from "@mui/x-tree-view/useTreeItem2/useTreeItem2";
import {TreeItem2Provider} from "@mui/x-tree-view/TreeItem2Provider";
import {TreeItem2Icon} from "@mui/x-tree-view/TreeItem2Icon";
import {Box} from "@mui/material";
import {DynamicIcon} from "../../DynamicIcon.tsx";
import {AccountNode} from "./AccountNode.ts";
import {LocaleContext} from "../../../providers/LocaleProvider.tsx";

interface TreeItemExProps
    extends Omit<UseTreeItem2Parameters, 'rootRef'>, Omit<HTMLAttributes<HTMLLIElement>, 'onFocus'> {
}

export const TreeItemEx = forwardRef(function CustomTreeItem(props: TreeItemExProps, ref: Ref<HTMLLIElement>) {
    const {id, itemId, label, disabled, children, ...other} = props;
    const {
        getRootProps,
        getContentProps,
        getIconContainerProps,
        getCheckboxProps,
        getLabelProps,
        getGroupTransitionProps,
        publicAPI,
        status,
    } = useTreeItem2({id, itemId, children, label, disabled, rootRef: ref});
    const locale = useContext(LocaleContext)

    const item: AccountNode = publicAPI.getItem(itemId) as AccountNode
    const currencyFormatter = new Intl.NumberFormat(navigator.language, {
        style: 'currency',
        currency: locale.currency,
    })

    return (
        <TreeItem2Provider itemId={itemId}>
            <TreeItem2Root {...getRootProps(other)}>
                <TreeItem2Content {...getContentProps()}>
                    <TreeItem2IconContainer {...getIconContainerProps()}>
                        <TreeItem2Icon status={status}/>
                    </TreeItem2IconContainer>
                    <TreeItem2Checkbox {...getCheckboxProps()} />
                    <Box sx={{flexGrow: 1, display: 'flex', gap: 1}}>
                        <DynamicIcon icon={item.icon || 'Work'}/>
                        <TreeItem2Label sx={{flexGrow: 1}} {...getLabelProps()} />
                        <Box sx={{
                            flexGrow: 0,
                            flexShrink: 0
                        }}>{item.balance ? currencyFormatter.format(item.balance) : "-"}</Box>
                    </Box>
                </TreeItem2Content>
                {children && <TreeItem2GroupTransition {...getGroupTransitionProps()} />}
            </TreeItem2Root>
        </TreeItem2Provider>
    );
});
