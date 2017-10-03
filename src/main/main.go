package main

import (
	"benchmark"
	"fmt"
	"os"
	"runtime"
)

func main() {
	var options = benchmark.ProcessArguments()
	if options.ShowHelp {
		fmt.Println(options.HelpText)
		os.Exit(1)
	}

	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))

	runner := benchmark.NewRunner(options)
	results := runner.Execute()
	benchmark.PrintResults(results)
}
