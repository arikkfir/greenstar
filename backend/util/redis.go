package util

import (
	"crypto/tls"
	"github.com/rueian/rueidis"
	"github.com/secureworks/errors"
	"strconv"
)

func CreateRedisClient(host string, port int, clientName string, useTLS bool) (rueidis.Client, error) {
	redisClientOption := rueidis.ClientOption{
		InitAddress: []string{host + ":" + strconv.Itoa(port)},
		ClientName:  clientName,
	}
	if useTLS {
		redisClientOption.TLSConfig = &tls.Config{ServerName: host}
	}
	redisClient, err := rueidis.NewClient(redisClientOption)
	if err != nil {
		return nil, errors.New("failed creating Redis client: %w", err)
	}
	return redisClient, nil
}
