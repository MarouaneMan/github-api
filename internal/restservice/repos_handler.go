package restservice

import (
	"encoding/json"
	"github.com/MarouaneMan/github-api/api"
	"github.com/MarouaneMan/github-api/kvstore"
	"github.com/Scalingo/go-utils/logger"
	"net/http"
)

// ReposHandler returns filtered Git repositories as JSON.
func ReposHandler(storeReader kvstore.Reader) func(http.ResponseWriter, *http.Request, map[string]string) error {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
		log := logger.Get(r.Context())

		queryParams := r.URL.Query()

		repositories, ok := storeReader.Read(r.Context(), "repositories").([]*api.Repository)
		if !ok {
			repositories = []*api.Repository{}
		}

		var filteredRepositories = FilterRepositories(
			NewFilterConfig(
				WithLanguage(queryParams.Get("language")),
				WithOwner(queryParams.Get("owner")),
				WithLimit(queryParams.Get("limit")),
				// WithLicense(queryParams.Get("license")), // not provided in /repositories
			),
			repositories,
		)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(filteredRepositories)
		if err != nil {
			log.WithError(err).Error("Fail to encode JSON")
		}
		return nil
	}
}
