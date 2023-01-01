package internal

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rueian/rueidis"
)

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	Redis rueidis.Client
	Neo4j neo4j.DriverWithContext
}
