package test

import (
	"benchmark"
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parser", func() {
	It("handles a simple string", func() {
		Expect(benchmark.Parse(bytes.NewBufferString("+OK\r\n"))).To(Equal("OK"))
	})

	It("handles errors", func() {
		Expect(benchmark.Parse(bytes.NewBufferString("-Error message\r\n"))).To(Equal("Error message"))
	})
})
