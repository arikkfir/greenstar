package sample

import (
	"bytes"
	"fmt"
	"github.com/shopspring/decimal"
)

func decimalFromValue(name string, v any) (decimal.Decimal, error) {
	if s, ok := v.(string); ok {
		tmpl, err := newTemplate(name, s)
		if err != nil {
			return decimal.Zero, fmt.Errorf("failed parsing '%s' template: %w", name, err)
		}
		var res bytes.Buffer
		if err := tmpl.Execute(&res, nil); err != nil {
			return decimal.Zero, fmt.Errorf("failed evaluating '%s' template: %w", name, err)
		} else if d, err := decimal.NewFromString(res.String()); err != nil {
			return decimal.Zero, fmt.Errorf("failed creating decimal from value '%s' of template '%s': %w", res.String(), name, err)
		} else {
			return d, nil
		}
	} else if i, ok := v.(int); ok {
		return decimal.NewFromInt32(int32(i)), nil
	} else if i, ok := v.(int8); ok {
		return decimal.NewFromInt32(int32(i)), nil
	} else if i, ok := v.(int16); ok {
		return decimal.NewFromInt32(int32(i)), nil
	} else if i, ok := v.(int32); ok {
		return decimal.NewFromInt32(i), nil
	} else if i, ok := v.(int64); ok {
		return decimal.NewFromInt(i), nil
	} else if i, ok := v.(uint); ok {
		return decimal.NewFromUint64(uint64(i)), nil
	} else if i, ok := v.(uint8); ok {
		return decimal.NewFromUint64(uint64(i)), nil
	} else if i, ok := v.(uint16); ok {
		return decimal.NewFromUint64(uint64(i)), nil
	} else if i, ok := v.(uint32); ok {
		return decimal.NewFromUint64(uint64(i)), nil
	} else if i, ok := v.(uint64); ok {
		return decimal.NewFromUint64(i), nil
	} else if f, ok := v.(float32); ok {
		return decimal.NewFromFloat32(f), nil
	} else if f, ok := v.(float64); ok {
		return decimal.NewFromFloat(f), nil
	} else if d, ok := v.(decimal.Decimal); ok {
		return d, nil
	} else {
		return decimal.Zero, fmt.Errorf("cannot create a decimal value from a '%T': %v", v, v)
	}
}
