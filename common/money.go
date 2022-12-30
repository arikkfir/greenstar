package common

import (
	"fmt"
	"strconv"
	"strings"
)

type Currency uint64

func (m Currency) Secondary() uint64 { return uint64(m) % 100 }
func (m Currency) Primary() uint64   { return uint64(m) / 100 }

func NewCurrency(primary, secondary uint64) Currency {
	return Currency((primary * 100) + secondary)
}
func ParseCurrency(s string) (Currency, error) {
	if s == "" {
		return 0, nil
	}

	tokens := strings.Split(s, ".")
	if len(tokens) == 1 {
		tokens = append(tokens, "0")
	} else if len(tokens) > 2 {
		return 0, fmt.Errorf("invalid currency format: %s", s)
	}

	primary, err := strconv.ParseUint(tokens[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse primary currency: %w", err)
	}

	secondary, err := strconv.ParseUint(tokens[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse secondary currency: %w", err)
	}

	return NewCurrency(primary, secondary), nil
}
