import { ColumnsPanelTrigger, Toolbar, ToolbarButton, ToolbarProps } from "@mui/x-data-grid-premium"
import { useRef } from "react"
import { Tooltip } from "@mui/material"
import AddIcon from "@mui/icons-material/Add"
import ViewColumnIcon from "@mui/icons-material/ViewColumn"

export interface CustomToolbarProps extends ToolbarProps {
    onAddScraperClicked: () => void
}

export function CustomToolbar({ onAddScraperClicked, ...rest }: CustomToolbarProps) {
    const newScraperPanelTriggerRef = useRef<HTMLButtonElement>(null)
    return (
        <Toolbar {...rest}>
            <Tooltip title="New scraper">
                <ToolbarButton ref={newScraperPanelTriggerRef} onClick={onAddScraperClicked}>
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
}
