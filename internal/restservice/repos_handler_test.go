package restservice

import (
	"context"
	"encoding/json"
	"github.com/MarouaneMan/github-api/api"
	"github.com/MarouaneMan/github-api/kvstore"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestReposHandler(t *testing.T) {

	store := kvstore.NewInMemoryStore(kvstore.DefaultExpiration, kvstore.DefaultExpiration)

	_ = store.Write(context.Background(), "repositories", reposMock, kvstore.NoExpiration)

	// Create a request with the desired query parameters
	req, err := http.NewRequest("GET", "/repositories?language=golang&owner=owner3&limit=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := ReposHandler(store)
	err = handler(rr, req, nil)
	if err != nil {
		t.Fatalf("Failed to handle reposMock request: %v", err)
		return
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("unexpected status code: got %v, want %v", status, http.StatusOK)
	}

	// Parse the response body and check its content
	var responseRepositories []*api.Repository
	err = json.Unmarshal(rr.Body.Bytes(), &responseRepositories)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON response: %v", err)
	}

	expected := []*api.Repository{
		{
			Repository: "repo3",
			Owner:      "owner3",
			Languages: map[string]api.Language{
				"golang": {
					Bytes: 1234,
				},
			},
		},
	}

	if !reflect.DeepEqual(responseRepositories, expected) {
		t.Errorf("Fetched repositories do not match expected: got %+v, want %+v", responseRepositories, expected)
	}
}
