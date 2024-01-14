package api

type Stats struct {

	// Target programming language for statistics
	Language string `json:"language"`

	// Total number of repositories using the language
	TotalUsage uint `json:"total_usage"`

	// Total code size in bytes across all repositories using the language
	TotalCodeSize uint `json:"total_code_size"`

	// Total number of repositories in the dataset
	TotalRepositories uint `json:"total_repositories"`

	// Average code size in bytes per repository using the languag
	AverageCodeSize uint `json:"avg_code_size"`
}
