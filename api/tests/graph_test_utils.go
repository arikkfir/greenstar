package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/arikkfir/greenstar/api/util"
	"io"
	"net/http"
	"strings"
	"testing"
)

func createTenant(t *testing.T) string {
	t.Helper()

	tenant := util.RandomTenant(7)
	url := "http://localhost:8000/tenants?tenant=" + tenant

	t.Logf("Creating tenant: %s", tenant)
	t.Logf("Tenant creation endpoint: %s", url)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		t.Fatal("Failed to create HTTP request to delete tenant", err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal("Failed to send HTTP request to delete tenant", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to delete tenant (expected HTTP status code %d, got %d)", http.StatusOK, resp.StatusCode)
	}

	tenantURL := "http://localhost:8000/tenants/" + tenant
	t.Cleanup(func() {
		t.Log("Deleting tenant")
		request, err := http.NewRequest(http.MethodDelete, tenantURL, nil)
		if err != nil {
			t.Fatal("Failed to create HTTP request to delete tenant", err)
		}

		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatal("Failed to send HTTP request to delete tenant", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to delete tenant (expected HTTP status code %d, got %d)", http.StatusOK, resp.StatusCode)
		}
	})
	return tenantURL + "/query"
}

func queryGraphQLEndpoint(url string, query string) (string, error) {
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return "", fmt.Errorf("failed to marshal query to JSON string: %w\nQuery string: %s", err, query)
	}

	bodyReader := strings.NewReader(fmt.Sprintf(`{"query": %s, "variables": {}}`, string(queryBytes)))
	request, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request to '%s': %w", request.URL, err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		buffer := bytes.Buffer{}
		buffer.WriteString(fmt.Sprintf("URL: %s\n", url))
		buffer.WriteString(fmt.Sprintf("HTTP method: %s\n", request.Method))
		buffer.WriteString(fmt.Sprintf("HTTP status code: %d\n", resp.StatusCode))
		buffer.WriteString(fmt.Sprintf("HTTP response headers:\n"))
		if err := resp.Header.Write(&buffer); err != nil {
			buffer.WriteString(fmt.Sprintf("Failed writing HTTP headers to buffer: %s\n", err))
		}
		buffer.WriteString(fmt.Sprintf("HTTP response body: %s\n", string(b)))
		return "", fmt.Errorf("unexpected status code:\n%s", buffer.String())
	}

	return string(b), nil
}
