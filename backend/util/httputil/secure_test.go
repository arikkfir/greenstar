package httputil_test

import (
	"github.com/arikkfir/greenstar/backend/util/httputil"
	"testing"
)

func TestIsSecure(t *testing.T) {
	var secure bool
	var err error

	secure, err = httputil.IsSecure("https://example.com")
	if err != nil {
		t.Fatal(err)
	} else if !secure {
		t.Errorf("Expected https://example.com to be secure")
	}

	secure, err = httputil.IsSecure("http://example.com")
	if err != nil {
		t.Fatal(err)
	} else if secure {
		t.Errorf("Expected http://example.com to be secure")
	}

	secure, err = httputil.IsSecure("HTTPS://EXAMPLE.COM")
	if err != nil {
		t.Fatal(err)
	} else if !secure {
		t.Errorf("Expected https://example.com to be secure")
	}

	secure, err = httputil.IsSecure("HTTP://EXAMPLE.COM")
	if err != nil {
		t.Fatal(err)
	} else if secure {
		t.Errorf("Expected http://example.com to be secure")
	}
}
