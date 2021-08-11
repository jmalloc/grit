package config_test

import (
	. "github.com/jmalloc/grit/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Load()", func() {
	It("treats an empty configuration the same as the defaults", func() {
		def, err := Load("testdata/default.toml")
		Expect(err).ShouldNot(HaveOccurred())

		empty, err := Load("testdata/empty.toml")
		Expect(err).ShouldNot(HaveOccurred())

		Expect(empty).To(Equal(def))
	})

	It("returns an error if the configuration contains unrecognized keys", func() {
		_, err := Load("testdata/unrecognized.toml")
		Expect(err).To(MatchError("unrecognized keys: unrecognized"))
	})

	It("returns an error if the configuration is malformed", func() {
		_, err := Load("testdata/malformed.txt")
		Expect(err).Should(HaveOccurred())
	})

	When("the is a custom github source defined", func() {
		It("returns an error if the clone_dir setting is empty", func() {
			_, err := Load("testdata/github.empty-clone-dir.toml")
			Expect(err).To(MatchError("sources.github.ghe.clone_dir is empty"))
		})

		It("returns an error if the api_url setting is empty", func() {
			_, err := Load("testdata/github.empty-api-url.toml")
			Expect(err).To(MatchError("sources.github.ghe.api_url is empty"))
		})
	})
})
