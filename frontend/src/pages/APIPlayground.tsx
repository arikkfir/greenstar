import "swagger-ui-react/swagger-ui.css"
import SwaggerUI from "swagger-ui-react";
import {BaseAPIURL} from "../client/common.ts";
import {useSession} from "@descope/react-sdk";

export function APIPlaygroundPage() {
    const {sessionToken} = useSession()
    return (
        <SwaggerUI url={`${BaseAPIURL}/openapi.yaml`}
                   displayRequestDuration={true}
                   requestInterceptor={(req) => {
                       req.headers['Authorization'] = `Bearer ${sessionToken}`
                       req.credentials = "include"
                       return req
                   }}
                   withCredentials={true}
                   onComplete={(system) => system.authActions.authorize({
                       Bearer: {value: `Bearer ${sessionToken}`}
                   })}
                   queryConfigEnabled={true}
                   showCommonExtensions={true}
                   showExtensions={true}
                   showMutatedRequest={true}
                   tryItOutEnabled={true}/>
    )
}
