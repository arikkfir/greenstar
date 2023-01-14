package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func queryGraphQLEndpoint(endpoint string, query string) (string, error) {
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return "", fmt.Errorf("failed to marshal query to JSON string: %w\nQuery string: %s", err, query)
	}

	bodyReader := strings.NewReader(fmt.Sprintf(`{"query": %s, "variables": {}}`, string(queryBytes)))
	request, err := http.NewRequest(http.MethodPost, endpoint, bodyReader)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
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
		buffer.WriteString(fmt.Sprintf("HTTP status code: %d\n", resp.StatusCode))
		buffer.WriteString(fmt.Sprintf("HTTP response headers:\n"))
		if err := resp.Header.Write(&buffer); err != nil {
			buffer.WriteString(fmt.Sprintf("Failed writing HTTP headers to buffer: %s\n", err))
		}
		buffer.WriteString(fmt.Sprintf("HTTP response body: %s\n", string(b)))
		return "", fmt.Errorf("response error:\n%s", buffer.String())
	}

	return string(b), nil
}
