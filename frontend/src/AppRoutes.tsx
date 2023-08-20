import {Route, Routes} from "react-router-dom";
import {HomePage} from "./pages/HomePage.tsx";
import {APIExplorer} from "./pages/APIExplorer.tsx";
import {NotFoundPage} from "./pages/NotFoundPage.tsx";
import {Tenants} from "./pages/Tenants.tsx";
import {Accounts} from "./pages/Accounts.tsx";

interface AppRoutesProps {
    tenant: string
}

export function AppRoutes({tenant}: AppRoutesProps) {
    return (
        <Routes>
            <Route path="/" element={<HomePage/>}/>
            <Route path="/accounts" element={<Accounts tenantID={tenant}/>}/>
            <Route path="/tenants" element={<Tenants/>}/>
            <Route path="/api" element={<APIExplorer/>}/>
            <Route path="*" element={<NotFoundPage/>}/>
        </Routes>
    )
}
