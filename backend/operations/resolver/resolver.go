package resolver

import (
	"github.com/rueian/rueidis"
)

type Resolver struct {
	Redis rueidis.Client
}
