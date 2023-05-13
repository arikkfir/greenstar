package auth

type Config struct {
	Auth0Domain     string `env:"AUTH0_DOMAIN" value-name:"DOMAIN" long:"auth0-domain" description:"Auth0 domain" required:"yes"`
	APIClientID     string `env:"AUTH0_API_CLIENT_ID" value-name:"ID" long:"auth0-api-client-id" description:"Client ID of the Auth0 Greenstar API application" required:"yes"`
	APIClientSecret string `env:"AUTH0_API_CLIENT_SECRET" value-name:"SECRET" long:"auth0-api-client-secret" description:"Client secret of the Auth0 Greenstar API application" required:"yes"`
}
