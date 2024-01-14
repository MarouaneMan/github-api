package config

type Config struct {
	Port               int    `envconfig:"PORT" default:"5000"`
	GithubToken        string `envconfig:"GITHUB_TOKEN" required:"True"`
	FetchIntervalHours int    `envconfig:"FETCH_INTERVAL_HOURS" default:"3"`
}
