package model

import (
	"context"
	"fmt"
	"github.com/secureworks/errors"
	"io"
	"strconv"
	"strings"
)

type Money uint64

//goland:noinspection GoMixedReceiverTypes
func (m Money) Secondary() uint64 { return uint64(m) % 100 }

//goland:noinspection GoMixedReceiverTypes
func (m Money) Primary() uint64 { return uint64(m) / 100 }

//goland:noinspection GoMixedReceiverTypes
func (m Money) String() string { return fmt.Sprintf("%d.%02d", m.Primary(), m.Secondary()) }

//goland:noinspection GoMixedReceiverTypes
func (m *Money) UnmarshalGQLContext(_ context.Context, v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("failed to unmarshal money: money must be a string")
	}
	money, err := ParseMoney(s)
	if err != nil {
		return errors.New("failed to unmarshal money: %w", err)
	}
	*m = money
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (m Money) MarshalGQLContext(ctx context.Context, w io.Writer) error {
	s, err := m.FormatContext(ctx)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(strconv.Quote(s)))
	if err != nil {
		return errors.New("failed to marshal money: %w", err)
	}
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (m Money) FormatContext(_ context.Context) (string, error) {
	return fmt.Sprintf("%d.%02d", m.Primary(), m.Secondary()), nil
}

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
		return 0, errors.New("invalid currency format: %s", s)
	}

	primary, err := strconv.ParseUint(tokens[0], 10, 64)
	if err != nil {
		return 0, errors.New("failed to parse primary currency: %w", err)
	}

	secondary, err := strconv.ParseUint(tokens[1], 10, 64)
	if err != nil {
		return 0, errors.New("failed to parse secondary currency: %w", err)
	}

	return NewMoney(primary, secondary), nil
}

func MustParseMoney(s string) Money {
	m, err := ParseMoney(s)
	if err != nil {
		panic(err)
	}
	return m
}
