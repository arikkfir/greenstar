package internal

import "github.com/arikkfir/greenstar/common"

type Config struct {
	General common.GeneralConfig
	Redis   common.RedisConfig `group:"redis" namespace:"redis" env-namespace:"REDIS"`
}

func NewConfig() *Config {
	config := Config{}
	common.ReadConfig(&config)
	config.General.Apply()
	return &config
}
