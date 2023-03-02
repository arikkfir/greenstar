package main

import (
	"fmt"
	"github.com/arikkfir/greenstar/backend/util/httputil"
	"time"
)

type Config struct {
	LogLevel string     `env:"LOG_LEVEL" value-name:"LEVEL" long:"log-level" description:"Log level" default:"info" enum:"trace,debug,info,warn,error,fatal,panic"`
	DevMode  bool       `env:"DEV_MODE" long:"dev-mode" description:"Development mode"`
	HTTP     HTTPConfig `group:"http" namespace:"http" env-namespace:"HTTP"`
	Auth     AuthConfig `group:"auth" namespace:"auth" env-namespace:"AUTH"`
}

func (c *Config) IsDevMode() bool {
	return c.DevMode
}

func (c *Config) GetLogLevel() string {
	return c.LogLevel
}

type HTTPConfig struct {
	Port             int        `env:"PORT" value-name:"PORT" long:"port" description:"Port to listen on" default:"8000"`
	StateCookieName  string     `env:"STATE_COOKIE_NAME" value-name:"NAME" long:"state-cookie-name" description:"Name of the cookie used to store the OAuth state" default:"oauthstate"`
	ClaimsCookieName string     `env:"CLAIMS_COOKIE_NAME" value-name:"NAME" long:"claims-cookie-name" description:"Name of the cookie used to store the OAuth claims" default:"claims"`
	AppURL           string     `env:"APP_URL" value-name:"URL" long:"app-url" description:"URL to redirect to after login or logout" required:"yes"`
	CORS             CORSConfig `group:"cors" namespace:"cors" env-namespace:"CORS"`
}

type CORSConfig struct {
	AllowedOrigins     []string      `env:"ALLOWED_ORIGINS" value-name:"ORIGIN" long:"allowed-origins" description:"List of origins a cross-domain request can be executed from (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin)" required:"yes"`
	AllowMethods       []string      `env:"ALLOWED_METHODS" value-name:"METHOD" long:"allowed-methods" description:"List of HTTP methods a cross-domain request can use (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods)" default:"GET"`
	AllowHeaders       []string      `env:"ALLOWED_HEADERS" value-name:"NAME" long:"allowed-headers" description:"List HTTP headers a cross-domain request can use (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers)" default:"*"`
	DisableCredentials bool          `env:"DISABLE_CREDENTIALS" long:"disable-credentials" description:"Disable access to credentials for JavaScript client code (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials)"`
	ExposeHeaders      []string      `env:"EXPOSE_HEADERS" long:"expose-headers" description:"List of HTTP headers to be made available to JavaScript browser code (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers)" default:"*"`
	MaxAge             time.Duration `env:"MAX_AGE" value-name:"DURATION" long:"max-age" description:"How long results of preflights response can be cached (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age)" default:"5s"`
}

func (c *HTTPConfig) ShouldUseSecureCookies() bool {
	secure, err := httputil.IsSecure(c.AppURL)
	if err != nil {
		panic(fmt.Errorf("failed to check if URL '%s' is secure: %w", c.AppURL, err))
	}
	return secure
}

type AuthConfig struct {
	Google AuthGoogleConfig `group:"google" namespace:"google" env-namespace:"GOOGLE"`
}

type AuthGoogleConfig struct {
	ClientID     string `env:"CLIENT_ID" value-name:"ID" long:"client-id" description:"Google OAuth client ID" required:"yes"`
	ClientSecret string `env:"CLIENT_SECRET" value-name:"SECRET" long:"client-secret" description:"Google OAuth client secret" required:"yes"`
	CallbackURL  string `env:"CALLBACK_URL" value-name:"URL" long:"callback-url" description:"Google OAuth callback URL" required:"yes"`
}
