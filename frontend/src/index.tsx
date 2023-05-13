import {Auth0Provider} from "@auth0/auth0-react";
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import {Typography} from "@mui/material";
import React from 'react';
import ReactDOM from 'react-dom/client';
import {App} from "./app/App";
import reportWebVitals from './reportWebVitals';

let auth0Domain = process.env.REACT_APP_AUTH0_DOMAIN ?? ""
if (!auth0Domain) {
    throw new Error("REACT_APP_AUTH0_DOMAIN is not set")
}

let auth0ClientID = process.env.REACT_APP_AUTH0_CLIENT_ID ?? ""
if (!auth0ClientID) {
    throw new Error("REACT_APP_AUTH0_CLIENT_ID is not set")
}

let tenant = window.location.host.split('.')[0]
if (tenant === "localhost") {
    if (!process.env.REACT_APP_AUTH0_ORG_NAME) {
        throw new Error("REACT_APP_AUTH0_ORG_NAME must be set when running in localhost")
    }
    tenant = process.env.REACT_APP_AUTH0_ORG_NAME
}

window.fetch(`/api/util/v1/organizations/${tenant}`)
    .then((response) => response.json())
    .then(json => {
        ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
            <React.StrictMode>
                <Auth0Provider
                    domain={auth0Domain}
                    clientId={auth0ClientID}
                    authorizationParams={{redirect_uri: window.location.origin}}
                >
                    <App
                        environment={process.env.NODE_ENV ?? "development"}
                        version={process.env.REACT_APP_VERSION ?? "local"}
                        organization={{id: json.id, name: json.name, displayName: json.displayName}}
                        adminAPIURL={process.env.REACT_APP_ADMIN_API_URL ?? "http://localhost/api/admin/playground"}
                        operationsAPIURL={process.env.REACT_APP_OPERATIONS_API_URL ?? "http://localhost/api/operations/playground"}
                        publicAPIURL={process.env.REACT_APP_PUBLIC_API_URL ?? "http://localhost/api/public/playground"}
                    />
                </Auth0Provider>
            </React.StrictMode>
        )
    })
    .catch((reason) => {
        console.error(`Failed to fetch tenant: `, reason)
        ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
            <React.StrictMode>
                <Typography>Failed to fetch tenant.</Typography>
            </React.StrictMode>
        );
    })

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
