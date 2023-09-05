import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import App from './App.tsx'
import {GrowthBook} from "@growthbook/growthbook-react";
import {BrowserRouter} from "react-router-dom";
import {AuthProvider} from "@descope/react-sdk";

// TODO: restore module declarations
// declare module 'react' {
//     // noinspection JSUnusedGlobalSymbols
//     interface CSSProperties {
//         '--tree-view-color'?: string;
//         '--tree-view-bg-color'?: string;
//     }
// }

// Augmenting mui "slotProps"
declare module '@mui/x-data-grid' {
    interface NoRowsOverlayPropsOverrides {
        message: string
    }

    interface ToolbarPropsOverrides {
        creationDisabled: boolean
        onCreate: () => void
    }
}

function getSubdomain(): string | null {
    const hostname = window.location.hostname;
    if (hostname == "localhost") {
        return "localhost"
    }
    const parts = hostname.split('.');
    if (parts.length >= 3) {
        return parts[0];
    } else {
        return null;
    }
}

const growthBook = new GrowthBook({
    apiHost: "https://cdn.growthbook.io",
    clientKey: "sdk-J805ykMvUGUOcdg",
    enableDevMode: true,
});

const rootElt = document.getElementById('root')!;
const root = createRoot(rootElt);
const tenant = getSubdomain();
if (tenant == null || tenant == "") {
    root.render(
        <StrictMode>
            <div>Tenant could not be found.</div>
        </StrictMode>,
    )
} else {
    root.render(
        <StrictMode>
            <BrowserRouter>
                <AuthProvider projectId='P2R7Kif8Zt1IJMwhGvgnR7EX9Mnu'>
                    <App growthBook={growthBook} tenant={tenant}/>
                </AuthProvider>
            </BrowserRouter>
        </StrictMode>,
    )
}
