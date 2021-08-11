package config

import (
	"fmt"
)

// normalize normalizes and validates the configuration.
func normalize(cfg *Config) error {
	if cfg.BaseDir == "" {
		cfg.BaseDir = "~/grit"
	}

	if err := normalizePublicGitHubPublic(cfg); err != nil {
		return err
	}

	for name, s := range cfg.Sources.GitHub {
		if err := normalizeGitHubSource(name, s); err != nil {
			return err
		}
	}

	return nil
}

// normalizePublicGitHubPublic ensures that there is a source for github.com
// present in the configuration.
func normalizePublicGitHubPublic(cfg *Config) error {
	s, ok := cfg.Sources.GitHub["public"]

	if !ok {
		if cfg.Sources.GitHub == nil {
			cfg.Sources.GitHub = map[string]*GitHubSource{}
		}

		s = &GitHubSource{}
		cfg.Sources.GitHub["public"] = s
	}

	if s.CloneDir == "" {
		s.CloneDir = "github.com"
	}

	if s.APIURL == "" {
		s.APIURL = "https://api.github.com"
	}

	return nil
}

// normalizeGitHubSource returns an error if s is invalid.
func normalizeGitHubSource(name string, s *GitHubSource) error {
	if s.CloneDir == "" {
		return fmt.Errorf("sources.github.%s.clone_dir is empty", name)
	}

	if s.APIURL == "" {
		return fmt.Errorf("sources.github.%s.api_url is empty", name)
	}

	return nil
}
