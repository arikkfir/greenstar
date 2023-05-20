package util

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/secureworks/errors"
)

func NewNeo4jDriver(host string, port int, tls bool) (neo4j.DriverWithContext, error) {
	var url string
	if tls {
		url = fmt.Sprintf("bolt+s://%s:%d", host, port)
	} else {
		url = fmt.Sprintf("bolt://%s:%d", host, port)
	}
	neo4jDriver, err := neo4j.NewDriverWithContext(url, neo4j.NoAuth())
	if err != nil {
		return nil, errors.New("failed to create Neo4j driver: %w", err)
	}
	return neo4jDriver, nil
}
