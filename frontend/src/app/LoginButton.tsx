import {useAuth0} from "@auth0/auth0-react";
import React from "react";

export interface LoginButtonProps {
    organizationId: string
}

export function LoginButton({organizationId: organization}: LoginButtonProps) {
    const {loginWithRedirect} = useAuth0();
    return <button onClick={() => loginWithRedirect({
        authorizationParams: {organization}
    })}>Log In</button>;
}
