import { Scraper, ScraperJob, Tenant } from "../../../../graphql/graphql.ts"
import { useTenantID } from "../../../../hooks/tenant.ts"
import { LazyLog } from "@melloware/react-logviewer"
import { useEffect, useMemo, useRef, useState } from "react"
import { Divider, Typography } from "@mui/material"
import "./ScraperJobDrawer.scss"
import { apiHost } from "../../../../util/ApolloClient.ts"
import { useQuery } from "@apollo/client"
import { ScraperJob as ScraperJobQuery } from "../../ScrapersQueries.ts"

export interface ScraperJobDrawerProps {
    scraperID: Scraper["id"]
    scraperJobID: ScraperJob["id"]
}

// No longer need page size since we're using websockets for streaming

export function ScraperJobDrawer({ scraperID, scraperJobID }: ScraperJobDrawerProps) {
    const tenantID: Tenant["id"]  = useTenantID()
    const [ logs, setLogs ]       = useState<string[]>([])
    const [ loading, setLoading ] = useState(true)
    const [ error, setError ]     = useState<string | null>(null)
    const wsRef                   = useRef<WebSocket | null>(null)

    const { data: job, loading: loadingJob } = useQuery(ScraperJobQuery, {
        skip: !scraperID,
        variables: {
            tenantID,
            scraperID: scraperID!,
            scraperJobID,
        },
    })

    // Set up WebSocket connection
    useEffect(() => {

        if (loadingJob || !job) {
            return
        }

        // Clean up previous connection if it exists
        if (wsRef.current) {
            wsRef.current.close()
        }

        setLoading(true)
        setLogs([])

        // Create WebSocket URL with query parameters
        const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"
        const wsUrl    = `${protocol}//${apiHost}/ws/scraper-job-logs?tenantID=${tenantID}&scraperID=${scraperID}&scraperJobID=${scraperJobID}`

        // Create WebSocket connection
        const ws      = new WebSocket(wsUrl)
        wsRef.current = ws

        // Handle WebSocket events
        ws.onopen = () => setLoading(false)

        ws.onmessage = (event) => setLogs(prevLogs => [ ...prevLogs, event.data ])

        ws.onerror = (event) => {
            console.error("WebSocket error:", event)
            setError("Failed to connect to log stream")
            setLoading(false)
        }

        ws.onclose = () => setLoading(false)

        // Clean up WebSocket connection on component unmount
        return () => {
            if (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING) {
                ws.close()
            }
        }
    }, [ tenantID, scraperID, scraperJobID, loadingJob, job ])

    // Combine logs into a single string
    const logText = useMemo(
        () => {
            if (!logs.length) {
                if (loading) {
                    return "Connecting to log stream..."
                } else if (error) {
                    return error
                } else {
                    return "No logs available."
                }
            } else {
                return logs.join("\n") + (loading ? "\n...waiting for more logs..." : "")
            }
        },
        [ logs, loading, error ],
    )

    return (
        <article className="scraper-job-drawer">
            <Typography
                variant="h2">{loadingJob ? "Loading..." : job?.tenant?.scraper?.job ? "Scraper Job details" : "Job not found"}</Typography>
            <Divider />
            {job?.tenant?.scraper?.job && <LazyLog caseInsensitive={true}
                                                   enableGutters={true}
                                                   enableHotKeys={true}
                                                   enableLineNumbers={true}
                                                   enableLinks={true}
                                                   enableMultilineHighlight={true}
                                                   enableSearch
                                                   enableSearchNavigation={true}
                                                   external={false}
                                                   extraLines={1}
                                                   follow
                                                   selectableLines={true}
                                                   text={logText} />}
        </article>
    )
}
