package redisutil

type RedisConfig struct {
	Host     string `env:"HOST" value-name:"HOST" long:"host" description:"Redis host name" default:"localhost"`
	Port     int    `env:"PORT" value-name:"PORT" long:"port" description:"Redis port" default:"6379"`
	PoolSize int    `env:"POOL_SIZE" value-name:"POOL_SIZE" long:"pool-size" description:"Redis connection pool size" default:"3"`
}
