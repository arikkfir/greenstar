import { ReactNode, useMemo } from "react"
import { SvgIcon, SvgIconProps } from "@mui/material"

interface DynamicSvgIconProps extends SvgIconProps {
    svgString: string
}

export function DynamicIcon({ svgString, viewBox: propViewBox, ...rest }: DynamicSvgIconProps): ReactNode {
    const parsedSvgData = useMemo(() => {
        // Handle empty SVG string
        if (!svgString) {
            return { viewBox: propViewBox || "0 0 24 24", innerContent: null }
        }

        // This component relies on DOMParser, so it's primarily for client-side rendering.
        // Ensure the window and DOMParser are available.
        if (typeof window === "undefined" || typeof DOMParser === "undefined") {
            throw new Error("DynamicMuiIcon: DOMParser not available (e.g., during SSR without polyfill).")
        }

        // Create the DOM parser
        const parser  = new DOMParser()
        const doc     = parser.parseFromString(svgString.trim(), "image/svg+xml")
        const svgNode = doc.documentElement

        // Check if parsing resulted in an error document
        const parseError = svgNode.querySelector("parsererror")
        if (parseError) {
            throw new Error("Error parsing SVG string: " + parseError.textContent)
        }

        if (svgNode && svgNode.nodeName.toLowerCase() === "svg") {
            const vb      = svgNode.getAttribute("viewBox")
            const content = svgNode.innerHTML // Gets the inner content (paths, etc.)
            return {
                viewBox: vb || propViewBox || "0 0 24 24", // Prioritize viewBox from string
                innerContent: content,
            }
        } else {
            throw new Error("Given string does not seem like an SVG string: " + svgString)
        }
    }, [svgString, propViewBox])

    return (
        <SvgIcon viewBox={parsedSvgData.viewBox} {...rest}>
            <g dangerouslySetInnerHTML={{ __html: parsedSvgData.innerContent || "" }} />
        </SvgIcon>
    )
}
