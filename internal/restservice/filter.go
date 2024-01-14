package restservice

import (
	"github.com/MarouaneMan/github-api/api"
	"strconv"
)

type FilterOption func(*FilterConfig)

type FilterConfig struct {
	Language string
	Owner    string
	Limit    int
}

func NewFilterConfig(opts ...FilterOption) *FilterConfig {
	cfg := &FilterConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func WithLanguage(language string) FilterOption {
	return func(cfg *FilterConfig) {
		cfg.Language = language
	}
}

func WithOwner(owner string) FilterOption {
	return func(cfg *FilterConfig) {
		cfg.Owner = owner
	}
}

func WithLimit(limit string) FilterOption {
	return func(cfg *FilterConfig) {
		cfg.Limit, _ = strconv.Atoi(limit)
	}
}

func FilterRepositories(cfg *FilterConfig, repos []*api.Repository) []*api.Repository {
	var filteredRepos = make([]*api.Repository, 0)

	for _, repo := range repos {

		// Filter by language
		if cfg.Language != "" {
			_, exists := repo.Languages[cfg.Language]
			if exists == false {
				continue
			}
		}

		// Filter by owner
		if cfg.Owner != "" {
			if repo.Owner != cfg.Owner {
				continue
			}
		}

		// Limit results
		if cfg.Limit > 0 && len(filteredRepos) >= cfg.Limit {
			break
		}

		filteredRepos = append(filteredRepos, repo)
	}

	return filteredRepos
}
