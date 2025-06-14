import { createRoot } from "react-dom/client"
import { StrictMode } from "react"
import { muiXTelemetrySettings } from "@mui/x-license"
import { App } from "./App.tsx"
import type {} from "@mui/x-tree-view-pro/themeAugmentation"
import type {} from "@mui/x-date-pickers-pro/themeAugmentation"
import type {} from "@mui/x-date-pickers-pro/AdapterLuxon"

muiXTelemetrySettings.disableTelemetry()

createRoot(document.getElementById("root")!).render(
    import.meta.env.DEV
        ? (
            <StrictMode>
                <App />
            </StrictMode>
        )
        : (
            <App />
        ),
)
