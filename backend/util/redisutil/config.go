package redisutil

type RedisConfig struct {
	Host     string `env:"HOST" value-name:"HOST" long:"host" description:"Redis host name" default:"greenstar-redis"`
	Port     int    `env:"PORT" value-name:"PORT" long:"port" description:"Redis port" default:"6379"`
	TLS      bool   `env:"TLS" value-name:"TLS" long:"tls" description:"Whether to use TLS to connect to Redis"`
	PoolSize int    `env:"POOL_SIZE" value-name:"POOL_SIZE" long:"pool-size" description:"Redis connection pool size" default:"3"`
}
