package internal

import "github.com/arikkfir/greenstar/common"

type Config struct {
	General common.GeneralConfig
	Redis   common.RedisConfig `group:"redis" namespace:"redis" env-namespace:"REDIS"`
	Neo4j   common.Neo4jConfig `group:"neo4j" namespace:"neo4j" env-namespace:"NEO4J"`
	HTTP    HTTPConfig         `group:"http" namespace:"http" env-namespace:"HTTP"`
}

type HTTPConfig struct {
	Port int `env:"PORT" value-name:"PORT" long:"port" description:"Port to listen on" default:"8000"`
}

func NewConfig() *Config {
	config := Config{}
	common.ReadConfig(&config)
	config.General.Apply()
	return &config
}
