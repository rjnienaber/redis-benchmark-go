package main

import (
	"benchmark"
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var options = benchmark.ProcessArguments()
	if options.ShowHelp {
		fmt.Println(options.HelpText)
		os.Exit(1)
	}
	var benchmarkWg sync.WaitGroup
	var counter uint64
	var limit uint64 = uint64(options.Requests)

	results := benchmark.NewResults(options)
	for i := 0; i < options.Connections; i++ {
		benchmarkWg.Add(1)
		go func() {
			conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", options.Host, options.Port))
			if err != nil {
				panic(fmt.Sprintf("Couldn't connect to redis server: %v", err))
			}

			for {
				if atomic.AddUint64(&counter, 1) > limit {
					benchmarkWg.Done()
					break
				}

				start := time.Now()
				fmt.Fprintf(conn, "PING\r\n")
				result, err := benchmark.Parse(bufio.NewReader(conn))
				results.LogRun(start)

				if err != nil {
					panic(err)
				}

				if result != "PONG" {
					panic(fmt.Sprintf("Result should have been 'PONG' was '%v'", result))
				}
			}
		}()
	}

	stopThroughput := benchmark.PrintThroughput(results.Start, &counter)
	benchmarkWg.Wait()
	results.Stop()
	stopThroughput <- true
	benchmark.PrintResults(results)
}

// $ redis-benchmark
// time taken: 704.205418ms
// counter: 100000
