package config

type schema struct {
	Clone struct {
		Path  string
		Order []string
	}
	Index struct {
		Path string
	}
	Providers map[string]providerSchema
}

type providerSchema struct {
	Driver string
	Host   string
}
