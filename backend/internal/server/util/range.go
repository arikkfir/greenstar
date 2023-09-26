package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseRowsRange(r string) (int, int, error) {
	if !strings.HasPrefix(r, "rows=") {
		return 0, 0, fmt.Errorf("unsupported range '%s' given: %w", r, ErrRequestedRangeNotSatisfiable)
	}

	tokens := strings.SplitN(strings.TrimPrefix(r, "rows="), "-", 2)
	if len(tokens) == 0 {
		return 0, 0, fmt.Errorf("range missing start & stop indices: %w", ErrRequestedRangeNotSatisfiable)
	} else if len(tokens) != 2 {
		return 0, 0, fmt.Errorf("range must specify both start & stop indices: %w", ErrRequestedRangeNotSatisfiable)
	} else if start, err := strconv.ParseInt(tokens[0], 0, 0); err != nil {
		return 0, 0, fmt.Errorf("invalid range start index: %w", ErrRequestedRangeNotSatisfiable)
	} else if end, err := strconv.ParseInt(tokens[1], 0, 0); err != nil {
		return 0, 0, fmt.Errorf("invalid range end index: %w", ErrRequestedRangeNotSatisfiable)
	} else if start <= 0 || end <= 0 {
		return 0, 0, fmt.Errorf("range indices must be positibe: %w", ErrRequestedRangeNotSatisfiable)
	} else if end < start {
		return 0, 0, fmt.Errorf("range start cannot be greater than range end: %w", ErrRequestedRangeNotSatisfiable)
	} else {
		return int(start), int(end), nil
	}
}
