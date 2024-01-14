package api

type Repository struct {

	// Full name of the repository
	FullName string `json:"full_name"`

	// Owner of the repository
	Owner string `json:"owner"`

	// Name of the repository
	Repository string `json:"repository"`

	// Dictionary of used languages
	Languages map[string]Language `json:"languages"`
}

type Language struct {

	// Size of the language in bytes
	Bytes uint64 `json:"bytes"`
}
