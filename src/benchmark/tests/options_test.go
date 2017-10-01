package test

import (
	"benchmark"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var helpString = `Usage: redis-benchmark [options]

  -h, --help          displays help
  -H, --host String   Server hostname - default: 127.0.0.1
  -n, --requests Int  Total number of requests - default: 100000
  -c, --clients Int   Number of parallel connections - default: 50
  -t, --tests Array   Only run the comma separated list of tests. The test names are the same as the ones produced as output. - default: PING
  -p, --port Int      Server port - default: 6379
  
Version 1.0.0`

// TODO: use varadic arguments
func parseArguments(args ...string) benchmark.Options {
	return benchmark.ParseArguments(append([]string{"redis-benchmark"}, args...))
}

var _ = Describe("ProcessOptions", func() {
	It("handles help switch", func() {
		options := parseArguments("--help")
		Expect(options.ShowHelp).To(BeTrue())
		Expect(options.HelpText).To(Equal(helpString))

		options = parseArguments("-h")
		Expect(options.ShowHelp).To(BeTrue())
		Expect(options.HelpText).To(Equal(helpString))
	})

	It("handles host switch", func() {
		options := parseArguments("-H", "172.17.0.1")
		Expect(options.Host).To(Equal("172.17.0.1"))

		options = parseArguments("--host", "10.10.10.1")
		Expect(options.Host).To(Equal("10.10.10.1"))
	})

	It("handles requests switch", func() {
		options := parseArguments("-n", "3")
		Expect(options.Requests).To(Equal(3))

		options = parseArguments("--requests", "5")
		Expect(options.Requests).To(Equal(5))
	})

	It("handles connections switch", func() {
		options := parseArguments("-c", "3")
		Expect(options.Connections).To(Equal(3))

		options = parseArguments("--clients", "5")
		Expect(options.Connections).To(Equal(5))
	})

	It("handles tests switch", func() {
		options := parseArguments("-t", "PiNg")
		Expect(options.Tests).To(Equal([]string{"PING"}))

		options = parseArguments("--tests", "get,SET")
		Expect(options.Tests).To(Equal([]string{"GET", "SET"}))
	})

	It("handles port switch", func() {
		options := parseArguments("-p", "8080")
		Expect(options.Port).To(Equal(8080))

		options = parseArguments("--port", "1337")
		Expect(options.Port).To(Equal(1337))
	})

	It("handles defaults", func() {
		options := parseArguments()
		Expect(options.Host).To(Equal("127.0.0.1"))
		Expect(options.Port).To(Equal(6379))
		Expect(options.Requests).To(Equal(100000))
		Expect(options.Connections).To(Equal(50))
		Expect(options.Tests).To(Equal([]string{"PING"}))
	})

	It("handles parameters without modifiers", func() {
		var expectedHelpText = `Error: Incorrect parameters specified

` + helpString

		var testInvalidArg = func(arg string) {
			options := parseArguments(arg)
			Expect(options.ShowHelp).To(BeTrue())
			Expect(options.HelpText).To(Equal(expectedHelpText))
		}

		testInvalidArg("--host")
		testInvalidArg("--port")
		testInvalidArg("--requests")
		testInvalidArg("-c")
		testInvalidArg("-t")
	})

	It("handles invalid parameter types", func() {
		var expectedHelpText = `Error: Invalid type for parameter

` + helpString

		var testInvalidArg = func(args ...string) {
			options := parseArguments(args...)
			Expect(options.ShowHelp).To(BeTrue())
			Expect(options.HelpText).To(Equal(expectedHelpText))
		}

		testInvalidArg("--port", "abc")
		testInvalidArg("--requests", "ABC")
		testInvalidArg("-c", "two")
	})

	It("handles invalid negative parameter values", func() {
		var expectedHelpText = `Error: Parameter value must be non-negative

` + helpString

		var testInvalidArg = func(args ...string) {
			options := parseArguments(args...)
			Expect(options.ShowHelp).To(BeTrue())
			Expect(options.HelpText).To(Equal(expectedHelpText))
		}

		testInvalidArg("--port", "-1")
		testInvalidArg("--requests", "-2")
		testInvalidArg("-c", "-3")
	})

	It("handles invalid parameters", func() {
		var expectedHelpText = `Error: Invalid parameter: --proxy

` + helpString
		options := parseArguments("--proxy", "8080")
		Expect(options.ShowHelp).To(BeTrue())
		Expect(options.HelpText).To(Equal(expectedHelpText))
	})

	// TODO: parsing of integers must be non-negative
})
