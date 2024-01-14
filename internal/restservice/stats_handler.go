package restservice

import (
	"encoding/json"
	"github.com/MarouaneMan/github-api/api"
	"github.com/MarouaneMan/github-api/kvstore"
	"github.com/Scalingo/go-utils/logger"
	"net/http"
)

// StatsHandler returns statistics for Git repositories, based on the specified language query parameter.
func StatsHandler(storeReader kvstore.Reader) func(http.ResponseWriter, *http.Request, map[string]string) error {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
		log := logger.Get(r.Context())

		queryParams := r.URL.Query()

		repositories, ok := storeReader.Read(r.Context(), "repositories").([]*api.Repository)
		if !ok {
			repositories = []*api.Repository{}
		}

		language := queryParams.Get("language")
		if language == "" {
			language = "ruby" // fallback to ruby
		}

		stats := BuildReposStatistics(repositories, language)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(stats)
		if err != nil {
			log.WithError(err).Error("Fail to encode JSON")
		}
		return nil
	}
}
