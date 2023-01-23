package tests

import (
	"testing"
)

func TestCreateAccount(t *testing.T) {
	url := createTenant(t)

	// language=GraphQL
	const mutation = `
		mutation {
  			root: createAccount(
    			input: {
    				id:"root",
    				displayName: "Root", 
    				labels: [ { key:"k1", value:"v1" } ]
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

	response, err := queryGraphQLEndpoint(url, mutation)
	if err != nil {
		t.Fatal("GraphQL request failed!", err)
	}

	const expectedResponse = `{"data":{"root":{"id":"root","displayName":"Root","labels":[{"key":"k1","value":"v1"}]},"child":{"id":"child","displayName":"Child","labels":[{"key":"k2","value":"v2"}],"parent":{"id":"root"}}}}`
	if response != expectedResponse {
		t.Errorf("Incorrect response received! Expected:\n%s\nGot:\n%s", expectedResponse, response)
	}
}
