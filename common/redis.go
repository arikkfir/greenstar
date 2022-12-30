package common

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisConfig struct {
	Host               string        `env:"HOST" value-name:"HOST" long:"host" description:"Redis host name" default:"localhost"`
	Port               int           `env:"PORT" value-name:"PORT" long:"port" description:"Redis port" default:"6379"`
	PoolSize           int           `env:"POOL_SIZE" value-name:"POOL_SIZE" long:"pool-size" description:"Redis connection pool size" default:"3"`
	MinIdleConnections int           `env:"MIN_IDLE_CONNECTIONS" value-name:"MIN_IDLE_CONNECTIONS" long:"min-idle-connections" description:"Redis minimum idle connections" default:"1"`
	MaxConnectionAge   time.Duration `env:"MAX_CONNECTION_AGE" value-name:"MAX_CONNECTION_AGE" long:"max-connection-age" description:"Redis maximum connection age" default:"1h"`
	PoolTimeout        time.Duration `env:"POOL_TIMEOUT" value-name:"POOL_TIMEOUT" long:"pool-timeout" description:"Redis connection pool timeout" default:"10s"`
}

func NewRedisClient(cfg RedisConfig) *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			PoolFIFO:     true,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConnections,
			MaxConnAge:   cfg.MaxConnectionAge,
			PoolTimeout:  cfg.PoolTimeout,
		},
	)
}
