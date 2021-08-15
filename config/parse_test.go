package config_test

import (
	"net/url"

	"github.com/go-git/go-git/v5/plumbing/transport"
	. "github.com/jmalloc/grit/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("func ParseFile()", func() {
	DescribeTable(
		"it returns the expected configuration",
		func(filename string, expect Config) {
			cfg, err := ParseFile(filename)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cfg).To(Equal(expect))
		},
		Entry(
			"default configuration",
			"testdata/valid/default.conf",
			DefaultConfig,
		),
		Entry(
			"empty configuration (should be the same as the default)",
			"testdata/valid/empty.conf",
			DefaultConfig,
		),
		Entry(
			"implicit github source disabled",
			"testdata/valid/github-disabled.conf",
			Config{Dir: "~/grit"},
		),
		Entry(
			"implicit github source overridden",
			"testdata/valid/github-overridden.conf",
			Config{
				Dir: "~/grit",
				Sources: map[string]Source{
					"github": GitHubSource{
						SourceName: "github",
						API: &url.URL{
							Scheme: "https",
							Host:   "github.example.com",
						},
					},
				},
			},
		),
		Entry(
			"implicit github source augmented with token",
			"testdata/valid/github-augmented.conf",
			Config{
				Dir: "~/grit",
				Sources: map[string]Source{
					"github": GitHubSource{
						SourceName: "github",
						API: &url.URL{
							Scheme: "https",
							Host:   "api.github.com",
						},
						Token: "<token>",
					},
				},
			},
		),
		Entry(
			"custom git source defined",
			"testdata/valid/git-custom.conf",
			Config{
				Dir: "~/grit",
				Sources: map[string]Source{
					"github": DefaultConfig.Sources["github"],
					"my-company": GitSource{
						SourceName: "my-company",
						Endpoint: &transport.Endpoint{
							Protocol: "ssh",
							User:     "git",
							Host:     "git.example.com",
							Port:     22,
							Path:     "{repo}.git",
						},
					},
				},
			},
		),
		Entry(
			"custom github source defined",
			"testdata/valid/github-custom.conf",
			Config{
				Dir: "~/grit",
				Sources: map[string]Source{
					"github": DefaultConfig.Sources["github"],
					"my-company": GitHubSource{
						SourceName: "my-company",
						API: &url.URL{
							Scheme: "https",
							Host:   "github.example.com",
						},
						Token: "<token>",
					},
				},
			},
		),
	)

	DescribeTable(
		"it returns an error if there is a problem with the configuration",
		func(filename string, expect string) {
			_, err := ParseFile(filename)
			Expect(err).To(MatchError(expect))
		},
		Entry(
			`syntax error`,
			`testdata/invalid/syntax-error.conf`,
			`testdata/invalid/syntax-error.conf:3:4 unexpected token ";" (expected "=" Value ";")`,
		),
		Entry(
			`unrecognized global parameter`,
			`testdata/invalid/unrecognized-parameter.conf`,
			`testdata/invalid/unrecognized-parameter.conf:3:1 unrecognized "key" parameter`,
		),
		Entry(
			`git source missing "endpoint" key`,
			`testdata/invalid/git-missing-endpoint.conf`,
			`testdata/invalid/git-missing-endpoint.conf:4:1 missing required "endpoint" parameter in "git" source`,
		),
		Entry(
			`git source with unrecognized parameter`,
			`testdata/invalid/git-unrecognized-parameter.conf`,
			`testdata/invalid/git-unrecognized-parameter.conf:6:5 unrecognized "key" parameter in "my-company" source`,
		),
		Entry(
			`github source missing "api" key`,
			`testdata/invalid/github-missing-api.conf`,
			`testdata/invalid/github-missing-api.conf:4:1 missing required "api" parameter in "my-company" source`,
		),
		Entry(
			`github source with unrecognized parameter`,
			`testdata/invalid/github-unrecognized-parameter.conf`,
			`testdata/invalid/github-unrecognized-parameter.conf:6:5 unrecognized "key" parameter in "my-company" source`,
		),
	)
})
