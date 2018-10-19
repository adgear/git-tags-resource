package utils

// Source structure
type Source struct {
	PrivateKey         string `json:"private_key"`
	RepositoryName     string `json:"repository_name"`
	URI                string `json:"uri"`
	PrivateKeyPassword string `json:"private_password"`
	TagFilter          string `json:"tag_filter"`
	LatestOnly         string `json:"latest_only"`
}

// Check structure
type Check struct {
	Source Source
}

// Input struct
type Input struct {
	Source  Source
	Version map[string]string
}
