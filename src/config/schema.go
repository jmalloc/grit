package config

// config is the entire Grit configuration.
type config struct {
	Clone struct {
		Path  string
		Order []string
	}
	Providers map[string]provider
}

type provider struct {
	Driver string
	Host   string
}
