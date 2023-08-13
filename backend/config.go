package main

import "time"

type Config struct {
	LogLevel  string      `env:"LOG_LEVEL" value-name:"LEVEL" long:"log-level" description:"Log level" default:"info" enum:"trace,debug,info,warn,error,fatal,panic"`
	DevMode   bool        `env:"DEV_MODE" long:"dev-mode" description:"Development mode"`
	DevModeID string      `env:"DEV_MODE_ID" long:"dev-mode-id" description:"Development mode ID, used for generating developer-specific resources."`
	Redis     RedisConfig `group:"redis" namespace:"redis" env-namespace:"REDIS"`
	Neo4j     Neo4jConfig `group:"neo4j" namespace:"neo4j" env-namespace:"NEO4J"`
	HTTP      HTTPConfig  `group:"http" namespace:"http" env-namespace:"HTTP"`
	Auth      AuthConfig  `group:"auth" namespace:"auth" env-namespace:"AUTH"`
}

type Neo4jConfig struct {
	Host string `env:"HOST" value-name:"HOST" long:"host" description:"Neo4j host name" required:"yes"`
	Port int    `env:"PORT" value-name:"PORT" long:"port" description:"Neo4j port" required:"yes"`
	TLS  bool   `env:"TLS" value-name:"TLS" long:"tls" description:"Whether to use TLS to connect to Neo4j"`
}

type RedisConfig struct {
	Host string `env:"HOST" value-name:"HOST" long:"host" description:"Redis host name" required:"yes"`
	Port int    `env:"PORT" value-name:"PORT" long:"port" description:"Redis port" default:"6379"`
	TLS  bool   `env:"TLS" long:"tls" description:"Whether to use TLS to connect to Redis"`
}

type HTTPConfig struct {
	Port                       int        `env:"PORT" value-name:"PORT" long:"port" description:"Port to listen on" default:"8000"`
	DisableAccessLog           bool       `env:"DISABLE_ACCESS_LOG" long:"disable-access-log" description:"Disable access log"`
	HealthPort                 int        `env:"HEALTH_PORT" value-name:"PORT" long:"health-port" description:"Port to listen on for health checks" default:"9000"`
	CORS                       CORSConfig `group:"cors" namespace:"cors" env-namespace:"CORS"`
	AccessLogExcludedHeaders   []string   `env:"ACCESS_LOG_EXCLUDED_HEADERS" value-name:"PATTERN" long:"access-log-excluded-headers" description:"List of header patterns to exclude from the access log"`
	AccessLogExcludeRemoteAddr bool       `env:"ACCESS_LOG_EXCLUDE_REMOTE_ADDR" long:"access-log-exclude-remote-addr" description:"Exclude the client remote address from the access log"`
}

type CORSConfig struct {
	AllowedOrigins     []string      `env:"ALLOWED_ORIGINS" value-name:"ORIGIN" long:"allowed-origins" description:"List of origins a cross-domain request can be executed from (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin)" required:"yes"`
	AllowMethods       []string      `env:"ALLOWED_METHODS" value-name:"METHOD" long:"allowed-methods" description:"List of HTTP methods a cross-domain request can use (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods)"`
	AllowHeaders       []string      `env:"ALLOWED_HEADERS" value-name:"NAME" long:"allowed-headers" description:"List HTTP headers a cross-domain request can use (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers)" default:"accept" default:"x-greenstar-tenant-id" default:"authorization" default:"content-type"`
	DisableCredentials bool          `env:"DISABLE_CREDENTIALS" long:"disable-credentials" description:"Disable access to credentials for JavaScript client code (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials)"`
	ExposeHeaders      []string      `env:"EXPOSE_HEADERS" long:"expose-headers" description:"List of HTTP headers to be made available to JavaScript browser code (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers)"`
	MaxAge             time.Duration `env:"MAX_AGE" value-name:"DURATION" long:"max-age" description:"How long results of preflights response can be cached (https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age)" default:"60s"`
}

type AuthConfig struct {
	DescopeLogLevel      string `env:"DESCOPE_LOG_LEVEL" value-name:"LEVEL" long:"descope-log-level" description:"Descope log level (none, debug, info)" default:"info"`
	DescopeProjectID     string `env:"DESCOPE_PROJECT_ID" value-name:"ID" long:"descope-project-id" description:"Descope project ID" required:"yes"`
	DescopeManagementKey string `env:"DESCOPE_MANAGEMENT_KEY" value-name:"KEY" long:"descope-management-key" description:"Descope management key" required:"yes"`
}
