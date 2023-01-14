package tests

import (
	"fmt"
	"github.com/arikkfir/greenstar/api/util"
	"net/http"
	"os"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	scheme := os.Getenv("GQL_SCHEME")
	if scheme == "" {
		scheme = "http"
	}

	tenant := util.RandomTenant(7)
	endpoint := os.Getenv("GQL_ENDPOINT")
	if endpoint == "" {
		t.Fatal("GQL_ENDPOINT not set")
	} else {
		endpoint = fmt.Sprintf("%s://%s.%s", scheme, tenant, endpoint)
	}
	t.Logf("Endpoint: %s", endpoint)

	t.Cleanup(func() {
		t.Log("Cleaning up")
		request, err := http.NewRequest(http.MethodDelete, endpoint+"/tenant", nil)
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

	// language=GraphQL
	const mutation = `
		mutation {
  			root: createAccount(
    			input: {
    				id:"root",
    				displayName: "Root", 
    				labels: [ { key:"k1", value:"v12" } ]
    			}
  			) { id, displayName, labels { key, value } }
			child: createAccount(
    			input: {
      				id:"child", 
      				displayName: "Child", 
      				parentID:"root", 
      				labels: [ { key:"k2", value:"v2" } ] 
				}
  			) { id, displayName, labels { key, value }, parent { id } }
		}`

	response, err := queryGraphQLEndpoint(endpoint, mutation)
	if err != nil {
		t.Fatal("GraphQL request failed!", err)
	}
	t.Logf("Response: %s", response)

	const expectedResponse = `{"data":{"root":{"id":"root","displayName":"Root","labels":[{"key":"k1","value":"v1"}]},"child":{"id":"child","displayName":"Child","labels":[{"key":"k2","value":"v2"}],"parent":{"id":"root"}}}}`
	if response != expectedResponse {
		t.Errorf("Incorrect response received! Expected:\n%s\nGot:\n%s", expectedResponse, response)
	}
}
