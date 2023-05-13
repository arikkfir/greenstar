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

window.fetch(`/api/v1/config`)
    .then((response) => response.json())
    .then(json => {
        const orgID = json.organization.id
        const orgName = json.organization.name
        const orgDisplayName = json.organization.display_name
        const auth0Domain = json.auth0.domain
        const auth0ClientID = json.auth0.clientId
        ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
            <React.StrictMode>
                <Auth0Provider
                    domain={auth0Domain}
                    clientId={auth0ClientID}
                    authorizationParams={{redirect_uri: window.location.origin}}
                >
                    <App organization={{id: orgID, name: orgName, displayName: orgDisplayName}}
                         adminAPIURL={process.env.REACT_APP_ADMIN_API_URL ?? "http://localhost/api/admin/playground"}
                         operationsAPIURL={process.env.REACT_APP_OPERATIONS_API_URL ?? "http://localhost/api/operations/playground"}
                         publicAPIURL={process.env.REACT_APP_PUBLIC_API_URL ?? "http://localhost/api/public/playground"}
                    />
                </Auth0Provider>
            </React.StrictMode>
        )
    })
    .catch((reason) => {
        console.error(`Failed to fetch configuration: `, reason)
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
