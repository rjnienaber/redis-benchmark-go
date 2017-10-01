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

type BenchmarkRun struct {
	Start time.Time
	End   time.Time
}

func main() {
	var options = benchmark.ProcessArguments()
	if options.ShowHelp {
		fmt.Println(options.HelpText)
		os.Exit(1)
	}
	var wg sync.WaitGroup
	var benchmarkWg sync.WaitGroup
	start := time.Now()
	var counter uint64
	var limit uint64 = uint64(options.Requests)
	results := make(chan BenchmarkRun)
	benchmarkResults := benchmark.Results{
		Command:     "PING",
		Requests:    options.Requests,
		Connections: options.Connections,
	}
	benchmarkResults.ResponseTimes = make(map[int]int)

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

				runStart := time.Now()
				fmt.Fprintf(conn, "PING\r\n")
				result, err := benchmark.Parse(bufio.NewReader(conn))
				run := BenchmarkRun{Start: runStart, End: time.Now()}
				if err != nil {
					panic(err)
				}

				if result != "PONG" {
					panic(fmt.Sprintf("Result should have been 'PONG' was '%v'", result))
				}

				results <- run
			}
		}()
	}

	wg.Add(1)
	go func() {
		for run := range results {
			ms := float64(run.End.Sub(run.Start).Nanoseconds()) / 1000000.0
			if ms < 1 {
				benchmarkResults.ResponseTimes[1] += 1
			} else {
				benchmarkResults.ResponseTimes[int(ms)] += 1
			}
		}
		wg.Done()
	}()

	stopThroughput := benchmark.PrintThroughput(start, &counter)

	benchmarkWg.Wait()
	benchmarkResults.Elapsed = time.Since(start)
	close(results)

	stopThroughput <- true
	wg.Wait()

	benchmark.PrintResults(benchmarkResults)
}

// $ redis-benchmark
// time taken: 704.205418ms
// counter: 100000
