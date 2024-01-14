package restservice

import "github.com/MarouaneMan/github-api/api"

// BuildReposStatistics computes statistics for a given language across multiple repositories.
func BuildReposStatistics(repos []*api.Repository, targetLanguage string) *api.Stats {
	var stats api.Stats

	stats.Language = targetLanguage

	for _, repo := range repos {
		if langStats, ok := repo.Languages[targetLanguage]; ok {
			stats.TotalUsage++
			stats.TotalCodeSize += uint(langStats.Bytes)
		}
	}

	stats.TotalRepositories = uint(len(repos))
	if stats.TotalUsage > 0 {
		stats.AverageCodeSize = stats.TotalCodeSize / stats.TotalUsage
	}
	return &stats
}
