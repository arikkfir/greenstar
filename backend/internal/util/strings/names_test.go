package strings

import (
	"testing"
)

func TestCamelCaseToHumanReadable(t *testing.T) {
	testCases := map[string]string{
		"id":                  "ID",
		"abc":                 "ABC",
		"Abc":                 "ABC",
		"ABC":                 "ABC",
		"BrownFox":            "Brown Fox",
		"BrownABC":            "Brown ABC",
		"Brown-fox-JumpsOver": "Brown Fox Jumps Over",
	}
	for v, e := range testCases {
		value := v
		expected := e
		t.Run(value, func(t *testing.T) {
			actual := CamelCaseToHumanReadable(value)
			if actual != expected {
				t.Errorf("for %s I got %q, want %q", v, actual, expected)
			}
		})
	}
}
