package shell_test

import (
	"strings"

	. "github.com/jmalloc/grit/shell"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func NewExecutor()", func() {
	It("returns an Executor that writes escaped commands to the writer", func() {
		w := &strings.Builder{}
		exec := NewExecutor(w)

		err := exec("commandA", "arg1", "arg2")
		Expect(err).ShouldNot(HaveOccurred())

		err = exec("commandB", "arg1", "arg2")
		Expect(err).ShouldNot(HaveOccurred())

		expect := `'commandA' 'arg1' 'arg2'` + "\n"
		expect += `'commandB' 'arg1' 'arg2'` + "\n"
		Expect(w.String()).To(Equal(expect))
	})
})
