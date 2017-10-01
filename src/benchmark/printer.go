package benchmark

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type Results struct {
	Command       string
	Requests      int
	Elapsed       time.Duration
	Connections   int
	ResponseTimes map[int]int
}

// PrintThroughput prints the realtime throughput of the benchmark
func PrintThroughput(start time.Time, counter *uint64) chan bool {
	ticker := time.NewTicker(time.Millisecond * 250)
	stop := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-ticker.C:
				throughput := float64(*counter) / time.Since(start).Seconds()
				fmt.Printf("\rPING: %0.2f", throughput)
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
	return stop
}

// PrintResults prints the final results for benchmark
func PrintResults(results Results) {
	fmt.Printf("\r%-100v\r", "") // clear last line
	fmt.Printf("====== %s ======\n", results.Command)
	fmt.Printf("  %d requests completed in %0.2f seconds\n", results.Requests, results.Elapsed.Seconds())
	fmt.Printf("  %v concurrent clients\n", results.Connections)
	fmt.Println("  3 bytes payload")
	fmt.Println("  keep alive: 1")
	fmt.Println()

	times := make([]int, 0)
	for k := range results.ResponseTimes {
		times = append(times, k)
	}
	sort.Ints(times)
	runningTotal := 0
	for _, k := range times {
		runningTotal += results.ResponseTimes[k]
		percentage := math.Floor((float64(runningTotal)/float64(results.Requests))*10000.0) / 100.0
		fmt.Printf("%0.2f%% <= %d milliseconds\n", percentage, k)
	}

	throughput := float64(results.Requests) / results.Elapsed.Seconds()
	fmt.Printf("%0.2f requests per second\n", throughput)
	fmt.Println()
}
