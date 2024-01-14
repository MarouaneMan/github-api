package fetcher

import (
	"context"
	"github.com/MarouaneMan/github-api/api"
	"github.com/MarouaneMan/github-api/internal/config"
	"github.com/MarouaneMan/github-api/kvstore"
	"github.com/jarcoal/httpmock"
	"reflect"
	"testing"
)

const repositoriesResponseMock = `
[
{
	"name": "foo",
	"full_name": "gopher/foo",
	"owner": {
		"login": "gopher"
	},
	"url": "https://api.github.com/repos/gopher/foo"
},
{
	"name": "bar",
	"full_name": "gopher/bar",
	"owner": {
		"login": "gopher"
	},
	"url": "https://api.github.com/repos/gopher/bar"
}
]
`

const languagesResponseMockFirst = `
{
	"golang": 1234
}
`
const languagesResponseMockSecond = `
{
	"c++": 5678
}
`

func TestFetcher(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	store := kvstore.NewInMemoryStore(kvstore.DefaultExpiration, kvstore.DefaultExpiration)
	ctx := context.Background()

	// mock /repositories
	{
		httpmock.RegisterResponder("GET", `=~^.+/repositories\z`,
			httpmock.NewStringResponder(200, repositoriesResponseMock),
		)
	}

	// mock /repos/owner/repo_name/languages
	{
		httpmock.RegisterResponder("GET", `=~^.+/repos/gopher/foo/languages\z`,
			httpmock.NewStringResponder(200, languagesResponseMockFirst),
		)
		httpmock.RegisterResponder("GET", `=~^.+/repos/gopher/bar/languages\z`,
			httpmock.NewStringResponder(200, languagesResponseMockSecond),
		)
	}

	// run fetcher
	Run(ctx, &config.Config{}, store, httpmock.DefaultTransport)

	repositories, ok := store.Read(ctx, "repositories").([]*api.Repository)
	if !ok {
		t.Error("Failed to retrieve cached data from the store")
	}

	expected := []*api.Repository{
		{
			FullName:   "gopher/foo",
			Repository: "foo",
			Owner:      "gopher",
			Languages: map[string]api.Language{
				"golang": {
					Bytes: 1234,
				},
			},
		},
		{
			FullName:   "gopher/bar",
			Repository: "bar",
			Owner:      "gopher",
			Languages: map[string]api.Language{
				"c++": {
					Bytes: 5678,
				},
			},
		},
	}

	if !reflect.DeepEqual(repositories, expected) {
		t.Errorf("Fetched repositories do not match expected: got %+v, want %+v", repositories, expected)
	}
}
