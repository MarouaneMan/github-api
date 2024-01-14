package restservice

import "github.com/MarouaneMan/github-api/api"

var reposMock = []*api.Repository{
	{
		Repository: "repo1",
		Owner:      "owner1",
		Languages: map[string]api.Language{
			"golang": {
				Bytes: 1234,
			},
		},
	},
	{
		Repository: "repo2",
		Owner:      "owner2",
		Languages: map[string]api.Language{
			"c++": {
				Bytes: 1234,
			},
		},
	},
	{
		Repository: "repo3",
		Owner:      "owner3",
		Languages: map[string]api.Language{
			"golang": {
				Bytes: 1234,
			},
		},
	},
	{
		Repository: "repo4",
		Owner:      "owner3",
		Languages: map[string]api.Language{
			"golang": {
				Bytes: 1234,
			},
		},
	},
}
