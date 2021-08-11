package config

// Config holds Grit configuration.
type Config struct {
	BaseDir string `toml:"base_dir"`

	Sources struct {
		GitHub map[string]GitHubSource `toml:"github"`
	} `toml:"sources"`
}

// GitHubSource contains the configuration for a GitHub source.
type GitHubSource struct {
	Enabled  bool   `toml:"enabled"`
	CloneDir string `toml:"clone_dir"`
	APIURL   string `toml:"api_url"`
}
