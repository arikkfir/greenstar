package strings

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/lucasepe/codename"
	"math/rand"
	"strings"
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

func CamelCaseToHumanReadable(fieldName string) string {
	delimited := strcase.ToDelimited(fieldName, ' ')
	for i := 0; i < len(delimited); i++ {
		if i == 0 || delimited[i-1] == ' ' {
			delimited = delimited[:i] + strings.ToUpper(string(delimited[i])) + delimited[i+1:]
		}
	}

	acronyms := []string{
		"Id",
		"Abc",
	}
	for _, acronym := range acronyms {
		delimited = strings.ReplaceAll(delimited, acronym, strings.ToUpper(acronym))
	}
	return delimited
}
