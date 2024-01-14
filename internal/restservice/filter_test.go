package restservice

import (
	"github.com/MarouaneMan/github-api/api"
	"reflect"
	"testing"
)

func TestFilterRepositories(t *testing.T) {
	var filteredRepositories = FilterRepositories(
		NewFilterConfig(
			WithLanguage("golang"),
			WithOwner("owner3"),
			WithLimit("1"),
		),
		reposMock,
	)

	if len(filteredRepositories) != 1 {
		t.Fatalf("Expected 1 repository, but got %d", len(filteredRepositories))
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
	if !reflect.DeepEqual(filteredRepositories, expected) {
		t.Errorf("Fetched repositories do not match expected: got %+v, want %+v", filteredRepositories, expected)
	}
}
