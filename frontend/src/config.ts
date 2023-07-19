export const globalTenant = "global"
const apiHostname = "api." + window.location.hostname.split('.').slice(1).join('.')
const apiPort = window.location.port ? ":" + window.location.port : ""
export const apiURL = `${window.location.protocol}//${apiHostname}${apiPort}`
export const gqlQueryURL = `${apiURL}/query`
export const gqlPlaygroundURL = `${apiURL}/playground`
