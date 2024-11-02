package strings

import (
	"fmt"
	"math/rand"

	"github.com/lucasepe/codename"
)

var (
	rng *rand.Rand
)

func init() {
	if defaultRNG, err := codename.DefaultRNG(); err != nil {
		panic(fmt.Errorf("failed to initialize codename RNG: %w", err))
	} else {
		rng = defaultRNG
	}
}

func Name() string {
	return codename.Generate(rng, 0)
}
