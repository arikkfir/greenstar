package util

import (
	"github.com/rs/zerolog"
)

type Neo4jZerologBoltLogger struct {
	Logger *zerolog.Logger
}

func (nl *Neo4jZerologBoltLogger) LogClientMessage(id, msg string, args ...interface{}) {
	nl.Logger.Trace().
		Str("type", "client").
		Str("id", id).
		Msgf(msg, args...)
}

func (nl *Neo4jZerologBoltLogger) LogServerMessage(id, msg string, args ...interface{}) {
	nl.Logger.Trace().
		Str("type", "server").
		Str("id", id).
		Msgf(msg, args...)
}
