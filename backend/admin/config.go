package main

type Config struct {
	LogLevel string      `env:"LOG_LEVEL" value-name:"LEVEL" long:"log-level" description:"Log level" default:"info" enum:"trace,debug,info,warn,error,fatal,panic"`
	DevMode  bool        `env:"DEV_MODE" long:"dev-mode" description:"Development mode"`
	NATS     NATSConfig  `group:"nats" namespace:"nats" env-namespace:"NATS"`
	Neo4j    Neo4jConfig `group:"neo4j" namespace:"neo4j" env-namespace:"NEO4J"`
	HTTP     HTTPConfig  `group:"http" namespace:"http" env-namespace:"HTTP"`
	Auth     AuthConfig  `group:"auth" namespace:"auth" env-namespace:"AUTH"`
}

func (c *Config) IsDevMode() bool {
	return c.DevMode
}

func (c *Config) GetLogLevel() string {
	return c.LogLevel
}

type NATSConfig struct {
	URL string `env:"URL" value-name:"URL" long:"url" description:"NATS URL" default:"nats://nats:4222"`
}

type Neo4jConfig struct {
	Host string `env:"HOST" value-name:"HOST" long:"host" description:"Neo4j host name" default:"localhost"`
	Port int    `env:"PORT" value-name:"PORT" long:"port" description:"Neo4j port" default:"7687"`
}

type HTTPConfig struct {
	Port             int    `env:"PORT" value-name:"PORT" long:"port" description:"Port to listen on" default:"8000"`
	ClaimsCookieName string `env:"CLAIMS_COOKIE_NAME" value-name:"NAME" long:"claims-cookie-name" description:"Name of the cookie used to store the OAuth claims" default:"claims"`
}

type AuthConfig struct {
	Google AuthGoogleConfig `group:"google" namespace:"google" env-namespace:"GOOGLE"`
}

type AuthGoogleConfig struct {
	ClientID     string `env:"CLIENT_ID" value-name:"ID" long:"client-id" description:"Google OAuth client ID" required:"yes"`
	ClientSecret string `env:"CLIENT_SECRET" value-name:"SECRET" long:"client-secret" description:"Google OAuth client secret" required:"yes"`
}
