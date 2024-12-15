import "swagger-ui-react/swagger-ui.css"
import SwaggerUI from "swagger-ui-react"
import { BaseAPIURL } from "../services/util.ts"

export function APIPlaygroundPage() {
    return (
        <SwaggerUI
            url={`${BaseAPIURL}/openapi.yaml`}
            displayRequestDuration={true}
            requestInterceptor={(req) => {
                req.credentials = "include"
                return req
            }}
            withCredentials={true}
            queryConfigEnabled={true}
            showCommonExtensions={true}
            showExtensions={true}
            showMutatedRequest={true}
            tryItOutEnabled={true}
        />
    )
}
