package natsutil

type NATSConfig struct {
	URL string `env:"URL" value-name:"URL" long:"url" description:"NATS URL" default:"nats://nats:4222"`
}
