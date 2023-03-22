package auth

import "time"

type Config struct {
	StateCookieName  string        `env:"STATE_COOKIE_NAME" value-name:"NAME" long:"state-cookie-name" description:"Name of the cookie used to store the OAuth state" default:"greenstar-oauth-state"`
	ClaimsCookieName string        `env:"CLAIMS_COOKIE_NAME" value-name:"NAME" long:"claims-cookie-name" description:"Name of the cookie used to store the OAuth claims" default:"greenstar-session"`
	SessionDuration  time.Duration `env:"SESSION_DURATION" value-name:"DURATION" long:"session-duration" description:"Duration of a session" default:"24h"`
	Google           GoogleConfig  `group:"google" namespace:"google" env-namespace:"GOOGLE"`
}

type GoogleConfig struct {
	ClientID     string `env:"CLIENT_ID" value-name:"ID" long:"client-id" description:"Google OAuth client ID" required:"yes"`
	ClientSecret string `env:"CLIENT_SECRET" value-name:"SECRET" long:"client-secret" description:"Google OAuth client secret" required:"yes"`
	CallbackURL  string `env:"CALLBACK_URL" value-name:"URL" long:"callback-url" description:"Google OAuth callback URL" required:"yes"`
}
