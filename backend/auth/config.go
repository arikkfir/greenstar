package auth

type Config struct {
	Auth0Domain           string `env:"AUTH0_DOMAIN" value-name:"DOMAIN" long:"auth0-domain" description:"Auth0 domain" required:"yes"`
	Auth0LocalhostOrgName string `env:"AUTH0_LOCALHOST_ORG_NAME" value-name:"NAME" long:"auth0-localhost-org-name" description:"Organization name to use when accessed via localhost"`
	APIClientID           string `env:"AUTH0_API_CLIENT_ID" value-name:"ID" long:"auth0-api-client-id" description:"Client ID of the Auth0 Greenstar API application" required:"yes"`
	APIClientSecret       string `env:"AUTH0_API_CLIENT_SECRET" value-name:"SECRET" long:"auth0-api-client-secret" description:"Client secret of the Auth0 Greenstar API application" required:"yes"`
	APPClientID           string `env:"AUTH0_APP_CLIENT_ID" value-name:"ID" long:"auth0-app-client-id" description:"Client ID of the Auth0 Greenstar frontend application" required:"yes"`
}
