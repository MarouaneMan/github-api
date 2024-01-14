package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MarouaneMan/github-api/api"
	"github.com/MarouaneMan/github-api/internal/config"
	"github.com/MarouaneMan/github-api/kvstore"
	"github.com/Scalingo/go-utils/logger"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strings"
	"time"
)

type githubOwner struct {
	Login string `json:"login"`
}

type githubRepository struct {
	Name      string      `json:"name"`
	FullName  string      `json:"full_name"`
	Owner     githubOwner `json:"owner"`
	URL       string      `json:"url"`
	Languages map[string]uint64
}

const (
	githubApiUrl     = "https://api.github.com"
	githubApiVersion = "2022-11-28"
)

// Run is the main function to fetch repositories and their languages from GitHub.
// It initializes the required components and orchestrates the fetching and storing process.
func Run(
	ctx context.Context,
	config *config.Config,
	storeWriter kvstore.Writer,
) {
	// Fetch all repositories
	log := logger.Get(ctx)
	log.Info("Fetching repositories...")
	defer func() {
		log.Info("Fetching repositories finished")
	}()

	// http.transport is reused to avoid creating too many connections
	httpTransport := &http.Transport{
		MaxConnsPerHost: 5, // do not overwhelm Github, http/2.0 takes care of concurrency
	}

	// fetch and parse repositories
	githubRepositories, err := fetchRepositories(ctx, config, httpTransport)
	if err != nil {
		log.WithError(err).Errorf("Failed to fetch repositories")
		return
	}

	// fetch repo languages
	{
		eg, egCtx := errgroup.WithContext(ctx)
		for _, repo := range githubRepositories {
			repo := repo // closure capture fixed in go 1.22 ?
			eg.Go(func() error {
				return fetchRepositoryLanguages(egCtx, config, httpTransport, repo)
			})
		}
		err = eg.Wait()
		if err != nil {
			log.WithError(err).Error("Failed to fetch repositories languages")
			return
		}
	}

	// transform and store repositories
	repositories := mapGithubReposToAPIRepos(githubRepositories)
	err = storeWriter.Write(ctx, "repositories", repositories, kvstore.NoExpiration)
	if err != nil {
		log.WithError(err).Error("Failed to write repositories to store")
	}
}

// fetchRepositories fetches a list of repositories from GitHub.
// It makes an HTTP GET request to the GitHub API and returns a slice of githubRepository.
func fetchRepositories(ctx context.Context, config *config.Config, httpTransport *http.Transport) ([]*githubRepository, error) {

	// Create a context with a timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	// Node: the endpoint /repositories does not return the most recent repositories but the oldest ones.
	// As of the time of writing this, there is no direct method using the '/search' endpoint with sort/q/order filters to fetch the latest repositories.
	req, _ := http.NewRequestWithContext(ctxWithTimeout, "GET", fmt.Sprintf("%s/repositories", githubApiUrl), nil)

	// Add the authorization/apiVersion headers to the request
	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.GithubToken))
	req.Header.Set("X-GitHub-Api-Version", githubApiVersion)

	// Send request
	client := &http.Client{
		Transport: httpTransport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute http request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Unexpected http statusCode = %d", resp.StatusCode)
	}

	// Decode JSON response
	var repositories []*githubRepository
	err = json.NewDecoder(resp.Body).Decode(&repositories)
	if err != nil {
		return nil, errors.Wrap(err, "Error while decoding JSON response")
	}
	return repositories, nil
}

// fetchRepositoryLanguages fetches programming languages for a given repository.
// The result is directly updated in the provided repo object.
func fetchRepositoryLanguages(ctx context.Context, config *config.Config, httpTransport *http.Transport, repo *githubRepository) error {

	// Create a context with a timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctxWithTimeout, "GET", fmt.Sprintf("%s/languages", repo.URL), nil)

	// Add the authorization/apiVersion headers to the request
	req.Header.Set("Authorization", fmt.Sprintf("token %s", config.GithubToken))
	req.Header.Set("X-GitHub-Api-Version", githubApiVersion)

	// Send request
	client := &http.Client{
		Transport: httpTransport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Failed to execute http request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("Unexpected http statusCode = %d", resp.StatusCode)
	}

	// Decode JSON response
	repo.Languages = make(map[string]uint64)
	err = json.NewDecoder(resp.Body).Decode(&repo.Languages)
	if err != nil {
		return errors.Wrap(err, "Error while decoding JSON response")
	}
	return nil
}

// mapGithubReposToAPIRepos transforms a slice of githubRepository to a slice of api.Repository.
func mapGithubReposToAPIRepos(githubRepositories []*githubRepository) []*api.Repository {
	result := make([]*api.Repository, 0, len(githubRepositories))
	for _, repo := range githubRepositories {
		repoModel := &api.Repository{
			FullName:   repo.FullName,
			Owner:      repo.Owner.Login,
			Repository: repo.Name,
			Languages:  map[string]api.Language{},
		}
		for lang, bytes := range repo.Languages {
			repoModel.Languages[strings.ToLower(lang)] = api.Language{
				Bytes: bytes,
			}
		}
		result = append(result, repoModel)
	}
	return result
}
