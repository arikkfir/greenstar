package e2e

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func assertNotZero(t *testing.T, expected any) {
	t.Helper()
	switch v := expected.(type) {
	case string:
		if expected == "" {
			t.Fatalf("Unexpected empty value")
		}
	case time.Time:
		if v.IsZero() {
			t.Fatalf("Unexpected zero value")
		}
	default:
		t.Fatalf("Unexpected type for 'assertNotZero': %T", v)
	}
}

func assertEquals(t *testing.T, expected, actual any, opts ...cmp.Option) {
	t.Helper()
	if !cmp.Equal(expected, actual, opts...) {
		t.Fatalf("Unexpected value encounted:\n%s", cmp.Diff(expected, actual, opts...))
	}
}

func assertNotEquals(t *testing.T, expected, actual any, opts ...cmp.Option) {
	t.Helper()
	if cmp.Equal(expected, actual, opts...) {
		t.Fatalf("Unexpected equal values encounted:\n%+v\n%+v", expected, actual)
	}
}

func assertIsNil(t *testing.T, actual any) {
	t.Helper()
	if actual != nil {
		t.Fatalf("Unexpected non-nil value: %+v", actual)
	}
}

func assertIsNotNil(t *testing.T, actual any) {
	t.Helper()
	if actual == nil {
		t.Fatalf("Unexpected nil value")
	}
}

func assertNotBefore(t *testing.T, expectedLater time.Time, expectedOlder time.Time) {
	if expectedOlder.Before(expectedOlder) {
		t.Errorf("Expected '%s' to be after '%s'", expectedLater, expectedOlder)
	}
}
