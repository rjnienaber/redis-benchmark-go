package benchmark

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Options represents command line arguments for benchmark
type Options struct {
	ShowHelp    bool
	Host        string
	Port        int
	Requests    int
	Connections int
	Tests       []string
	HelpText    string
}

var helpText = `Usage: redis-benchmark [options]

  -h, --help          displays help
  -H, --host String   Server hostname - default: 127.0.0.1
  -n, --requests Int  Total number of requests - default: 100000
  -c, --clients Int   Number of parallel connections - default: 50
  -t, --tests Array   Only run the comma separated list of tests. The test names are the same as the ones produced as output. - default: PING
  -p, --port Int      Server port - default: 6379
  
Version 1.0.0`

func buildHelp(message string) Options {
	return Options{ShowHelp: true, HelpText: message + `

` + helpText}
}

var emptyOptions = Options{}

func parseNumber(args []string, index *int) (int, Options) {
	*index++
	if *index >= len(args) {
		return -1, buildHelp("Error: Incorrect parameters specified")
	}
	number, err := strconv.Atoi(args[*index])
	if err != nil {
		return -1, buildHelp("Error: Invalid type for parameter")
	}

	if number < 0 {
		return -1, buildHelp("Error: Parameter value must be non-negative")
	}

	return number, emptyOptions
}

var defaultOptions = Options{
	Host:        "127.0.0.1",
	Port:        6379,
	Requests:    100000,
	Connections: 50,
	Tests:       []string{"PING"},
}

// ParseArguments parses a string array and returns a populated Options struct
func ParseArguments(arguments []string) Options {
	options := defaultOptions

	args := arguments[1:]
	var errOptions Options
	for i := 0; i < len(args); i++ {
		if args[i] == "--help" || args[i] == "-h" {
			return Options{ShowHelp: true, HelpText: helpText}
		} else if args[i] == "--host" || args[i] == "-H" {
			i++
			if i >= len(args) {
				return buildHelp("Error: Incorrect parameters specified")
			}
			options.Host = args[i]
		} else if args[i] == "--requests" || args[i] == "-n" {
			options.Requests, errOptions = parseNumber(args, &i)
			if errOptions.ShowHelp {
				return errOptions
			}
		} else if args[i] == "--clients" || args[i] == "-c" {
			options.Connections, errOptions = parseNumber(args, &i)
			if errOptions.ShowHelp {
				return errOptions
			}
		} else if args[i] == "--tests" || args[i] == "-t" {
			i++
			if i >= len(args) {
				return buildHelp("Error: Incorrect parameters specified")
			}
			options.Tests = strings.Split(args[i], ",")
			for i := range options.Tests {
				options.Tests[i] = strings.ToUpper(options.Tests[i])
			}
		} else if args[i] == "--port" || args[i] == "-p" {
			options.Port, errOptions = parseNumber(args, &i)
			if errOptions.ShowHelp {
				return errOptions
			}
		} else {
			return buildHelp(fmt.Sprintf("Error: Invalid parameter: %v", args[i]))
		}
	}

	return options
}

// ProcessArguments gets arguments from the command line and parses them
func ProcessArguments() Options {
	return ParseArguments(os.Args)
}
