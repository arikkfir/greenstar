import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import {AuthProvider} from '@descope/react-sdk'
import {createRoot} from 'react-dom/client';
import {BrowserRouter} from "react-router-dom";
import {App} from "./app";
import reportWebVitals from "./reportWebVitals";
import {StrictMode} from "react";
import { GrowthBook } from "@growthbook/growthbook-react";

// Augmenting mui "slotProps"
declare module '@mui/x-data-grid' {
    // noinspection JSUnusedGlobalSymbols
    interface NoRowsOverlayPropsOverrides {
        message: string
    }
}

declare module 'react' {
    // noinspection JSUnusedGlobalSymbols
    interface CSSProperties {
        '--tree-view-color'?: string;
        '--tree-view-bg-color'?: string;
    }
}

function getSubdomain(): string {
    const hostname = window.location.hostname;
    const parts = hostname.split('.');
    if (parts.length >= 3 ) {
        return parts[0];
    } else {
        throw new Error("could not extract subdomain from: " + window.location.hostname)
    }
}

const growthBook = new GrowthBook({
    apiHost: "https://cdn.growthbook.io",
    clientKey: "sdk-J805ykMvUGUOcdg",
    enableDevMode: true,
});

const container = document.getElementById("root") as HTMLElement;
const root = createRoot(container);
root.render(
    <StrictMode>
        <BrowserRouter>
            <AuthProvider projectId='P2R7Kif8Zt1IJMwhGvgnR7EX9Mnu'>
                <App growthBook={growthBook} tenant={getSubdomain()}/>
            </AuthProvider>
        </BrowserRouter>
    </StrictMode>
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
