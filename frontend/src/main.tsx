import React from 'react'
import ReactDOM from 'react-dom/client'
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import {App} from "./App.tsx";
import {AuthProvider} from "@descope/react-sdk";

export function render() {
    ReactDOM.createRoot(document.getElementById('root')!).render(
        <React.StrictMode>
            <AuthProvider sessionTokenViaCookie={false} projectId={import.meta.env.VITE_DESCOPE_PROJECT_ID}>
                <App/>
            </AuthProvider>
        </React.StrictMode>,
    )
}

render()
