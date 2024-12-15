import { Backdrop, CircularProgress } from "@mui/material"

export interface SpinnerBlockProps {
    open: boolean
}

export function SpinnerBlock({ open }: SpinnerBlockProps) {
    return (
        <Backdrop sx={(theme) => ({ color: "#fff", zIndex: theme.zIndex.drawer + 1 })} open={open}>
            <CircularProgress color="inherit" />
        </Backdrop>
    )
}
