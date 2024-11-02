package e2e

import (
	"github.com/gavv/httpexpect/v2"
	"io"
	"net/http"
	"testing"
)

const (
	baseURL = "https://api.greenstar.local"
)

func createHTTPExpect(t *testing.T, baseURL string) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		TestName: t.Name(),
		BaseURL:  baseURL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{httpexpect.NewDebugPrinter(t, true)},
	})
}

func createAuthenticatedHTTPExpect(t *testing.T, baseURL, jwt string) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		RequestFactory: &customRequestFactory{
			RequestFactory: httpexpect.DefaultRequestFactory{},
			Header:         map[string][]string{"Authorization": {"Bearer " + jwt}},
		},
		TestName: t.Name(),
		BaseURL:  baseURL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{httpexpect.NewDebugPrinter(t, true)},
	})
}

type customRequestFactory struct {
	httpexpect.RequestFactory
	Header http.Header
}

func (f *customRequestFactory) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	r, err := f.RequestFactory.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for name, values := range f.Header {
		for _, v := range values {
			r.Header.Add(name, v)
		}
	}
	return r, nil
}
