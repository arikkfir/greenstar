package bootutil

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"math/rand"
	"os"
	"time"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

type Config interface {
	IsDevMode() bool
	GetLogLevel() string
}

func Boot(config Config) {
	// Ensure we have a random seed
	rand.Seed(time.Now().UnixNano())

	// Parse configuration
	parser := flags.NewParser(config, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		fmt.Printf("ERROR: %s\n\n", err)
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	if config.IsDevMode() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		zerolog.DefaultContextLogger = &log.Logger
	}

	if level, err := zerolog.ParseLevel(config.GetLogLevel()); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	} else {
		zerolog.SetGlobalLevel(level)
	}
	zerolog.DefaultContextLogger = &log.Logger
}
