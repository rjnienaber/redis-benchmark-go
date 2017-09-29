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

	It("handles errors", func() {
		Expect(benchmark.Parse(bytes.NewBufferString("$6\r\nfoobar\r\n"))).To(Equal("foobar"))
	})

	It("handles empty string", func() {
		Expect(benchmark.Parse(bytes.NewBufferString("$0\r\n\r\n"))).To(Equal(""))
	})

	It("handles nil string", func() {
		Expect(benchmark.Parse(bytes.NewBufferString("$-1\r\n"))).To(BeNil())
	})

	It("handles integers", func() {
		Expect(benchmark.Parse(bytes.NewBufferString(":1000\r\n"))).To(Equal(1000))
	})

	It("handles empty array", func() {
		var expected []interface{}
		Expect(benchmark.Parse(bytes.NewBufferString("*0\r\n"))).To(Equal(expected))
	})

	It("handles string array", func() {
		expected := []interface{}{"foo", "bar"}
		Expect(benchmark.Parse(bytes.NewBufferString("*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"))).To(Equal(expected))
	})

	It("handles multiple empty arrays", func() {
		actual := benchmark.Parse(bytes.NewBufferString("*2\r\n*0\r\n*0\r\n")).([]interface{})
		var empty []interface{}
		Expect(actual).To(BeAssignableToTypeOf(empty))
		Expect(len(actual)).To(Equal(2))
		Expect(actual[0]).To(Equal(empty))
		Expect(actual[1]).To(Equal(empty))
	})

	It("handles arrays of nil strings", func() {
		expected := []interface{}{nil, nil}
		Expect(benchmark.Parse(bytes.NewBufferString("*2\r\n$-1\r\n$-1\r\n"))).To(Equal(expected))
	})

	It("handles arrays of empty strings", func() {
		expected := []interface{}{"", ""}
		Expect(benchmark.Parse(bytes.NewBufferString("*2\r\n$0\r\n$0\r\n"))).To(Equal(expected))
	})

	It("handles arrays of integers", func() {
		expected := []interface{}{1, 2, 3}
		Expect(benchmark.Parse(bytes.NewBufferString("*3\r\n:1\r\n:2\r\n:3\r\n"))).To(Equal(expected))
	})

	It("handles mixed types", func() {
		expected := []interface{}{1, 2, 3, 4, "foobar"}
		Expect(benchmark.Parse(bytes.NewBufferString("*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n"))).To(Equal(expected))
	})

	It("handles nested arrays", func() {
		expected := []interface{}{[]interface{}{1, 2, 3}, []interface{}{"Foo"}}
		Expect(benchmark.Parse(bytes.NewBufferString("*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*1\r\n+Foo\r\n"))).To(Equal(expected))
	})

	// It("handles nested arrays with error", func() {
	// 	expected := []interface{}{[]interface{}{1, 2, 3}, []interface{}{"Foo"}}
	// 	Expect(benchmark.Parse(bytes.NewBufferString("*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Foo\r\n-Bar\r\n"))).To(Equal(expected))
	// })

	It("handles null arrays", func() {
		Expect(benchmark.Parse(bytes.NewBufferString("*-1\r\n"))).To(BeNil())
	})

	It("handles nulls in arrays", func() {
		expected := []interface{}{"foo", nil, "bar"}
		Expect(benchmark.Parse(bytes.NewBufferString("*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n"))).To(Equal(expected))
	})
})
