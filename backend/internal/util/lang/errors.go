package lang

import (
	"errors"
)

func IgnoreErrorOfType(err error, ignored ...error) error {
	for _, ignore := range ignored {
		if errors.Is(err, ignore) {
			return nil
		}
	}
	return err
}
