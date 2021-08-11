package config_test

import (
	. "github.com/jmalloc/grit/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func Load()", func() {
	It("can load the default configuration", func() {
		_, err := Load("testdata/default.toml")
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("returns an error if the configuration contains unrecognized keys", func() {
		_, err := Load("testdata/unrecognized.toml")
		Expect(err).To(MatchError("unrecognized keys: unrecognized"))
	})

	It("returns an error if the configuration is malformed", func() {
		_, err := Load("testdata/malformed.txt")
		Expect(err).Should(HaveOccurred())
	})
})
