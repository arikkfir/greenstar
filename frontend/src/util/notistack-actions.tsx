import { Button } from "@mui/material"
import { closeSnackbar, SnackbarAction, SnackbarKey } from "notistack"

export const dismissAction: SnackbarAction = (snackbarId: SnackbarKey) => (
    <Button variant="text" onClick={() => { closeSnackbar(snackbarId) }}>Dismiss</Button>
)
