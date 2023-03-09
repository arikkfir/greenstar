import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import React from 'react';
import ReactDOM from 'react-dom/client';
import {App} from "./app/App";
import reportWebVitals from './reportWebVitals';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
    <React.StrictMode>
        <App environment={process.env.NODE_ENV ?? "development"}
             version={process.env.REACT_APP_VERSION ?? "local"}
             userInfoURL={process.env.REACT_USER_INFO_URL ?? "http://localhost:8000/auth/user"}
             loginURL={process.env.REACT_LOGIN_URL ?? "http://localhost:8000/auth/google/login"}
             adminAPIURL={process.env.REACT_ADMIN_API_URL ?? "http://localhost:8000/admin/playground"}
             operationsAPIURL={process.env.REACT_OPERATIONS_API_URL ?? "http://localhost:8000/operations/playground"}
             publicAPIURL={process.env.REACT_PUBLIC_API_URL ?? "http://localhost:8000/public/playground"}
        />
    </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
