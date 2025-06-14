import "./Layout.scss"
import { Paper } from "@mui/material"
import { ScrapersTable } from "./ScrapersTable.tsx"

export function ScrapersPage() {
    return (
        <main className="scrapers-page">
            <title>Scrapers - GreenSTAR</title>
            
            <Paper className="scrapers-container" elevation={3}>
                <ScrapersTable />
            </Paper>
        </main>
    )
}
