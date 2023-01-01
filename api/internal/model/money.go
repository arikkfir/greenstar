package model

import (
	"fmt"
	"strconv"
	"strings"
)

type Money uint64

func (m Money) Secondary() uint64 { return uint64(m) % 100 }

func (m Money) Primary() uint64 { return uint64(m) / 100 }

func NewMoney(primary, secondary uint64) Money {
	return Money((primary * 100) + secondary)
}

func ParseMoney(s string) (Money, error) {
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

	return NewMoney(primary, secondary), nil
}
