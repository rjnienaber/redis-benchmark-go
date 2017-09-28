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

	// it('handles multiple empty arrays', () => {
	// const buffer = new Buffer('*2\r\n*0\r\n*0\r\n');
	// const values = parser.parse(buffer);
	// assert.deepEqual(values, [[],[]]);
	// });

	// it('handles arrays of nil strings', () => {
	// const buffer = new Buffer('*2\r\n$-1\r\n$-1\r\n');
	// const values = parser.parse(buffer);
	// assert.deepEqual(values, [undefined, undefined]);
	// });

	// it('handles arrays of empty strings', () => {
	// const buffer = new Buffer('*2\r\n$0\r\n$0\r\n');
	// assert.deepEqual(parser.parse(buffer), ['', '']);
	// });

	// it('handles arrays of integers', () => {
	// const buffer = new Buffer('*3\r\n:1\r\n:2\r\n:3\r\n');
	// assert.deepEqual(parser.parse(buffer), [1, 2, 3]);
	// });

	// it('handles mixed types', () => {
	// const buffer = new Buffer('*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n');
	// assert.deepEqual(parser.parse(buffer), [1, 2, 3, 4, 'foobar']);
	// });

	// it('handles nested arrays', () => {
	// const buffer = new Buffer('*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*1\r\n+Foo\r\n');
	// assert.deepEqual(parser.parse(buffer), [[1, 2, 3], ['Foo']]);
	// });

	// it('handles nested arrays with error', () => {
	// const buffer = new Buffer('*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Foo\r\n-Bar\r\n');
	// expect(() => parser.parse(buffer)).to.throw('Bar');
	// });

	// it('handles null arrays', () => {
	// const buffer = new Buffer('*-1\r\n');
	// assert.equal(parser.parse(buffer), undefined);
	// });

	// it('handles nulls in arrays', () => {
	// const buffer = new Buffer('*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n');
	// assert.deepEqual(parser.parse(buffer), ['foo', undefined, 'bar']);
	// });
})
