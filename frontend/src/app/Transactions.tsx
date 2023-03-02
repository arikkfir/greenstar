import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import TreeItem from '@mui/lab/TreeItem';
import TreeView from '@mui/lab/TreeView';
import {Paper, Stack} from "@mui/material";
import React from "react";

export function Transactions() {
    return (
        <Stack direction="row">
            <Paper elevation={3}>
                <TreeView defaultCollapseIcon={<ExpandMoreIcon/>}
                          defaultExpandIcon={<ChevronRightIcon/>}
                          sx={{height: 240, flexGrow: 1, maxWidth: 400, overflowY: 'auto'}}>
                    <TreeItem nodeId="1" label="Applications">
                        <TreeItem nodeId="2" label="Calendar"/>
                    </TreeItem>
                    <TreeItem nodeId="5" label="Documents">
                        <TreeItem nodeId="10" label="OSS"/>
                        <TreeItem nodeId="6" label="MUI">
                            <TreeItem nodeId="8" label="index.js"/>
                        </TreeItem>
                    </TreeItem>
                </TreeView>
            </Paper>
        </Stack>
    );
}
