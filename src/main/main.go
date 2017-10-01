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
	var wg sync.WaitGroup
	start := time.Now()
	var counter uint64
	var limit uint64 = uint64(options.Requests)
	for i := 0; i < options.Connections; i++ {
		go func() {
			conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", options.Host, options.Port))
			if err != nil {
				panic(fmt.Sprintf("Couldn't connect to redis server: %v", err))
			}

			for {
				fmt.Fprintf(conn, "PING\r\n")
				result, err := benchmark.Parse(bufio.NewReader(conn))
				if err != nil {
					panic(err)
				}

				if result != "PONG" {
					panic(fmt.Sprintf("Result should have been 'PONG' was '%v'", result))
				}

				if atomic.LoadUint64(&counter) == limit {
					wg.Done()
					break
				}
				atomic.AddUint64(&counter, 1)
			}
		}()
		wg.Add(1)
	}

	ticker := time.NewTicker(time.Millisecond * 250)
	go func() {
		for _ = range ticker.C {
			fmt.Printf("\rPING: %0.2f", benchmark.ShowThroughput(start, counter))
		}
	}()

	wg.Wait()
	ticker.Stop()
	fmt.Printf("\r%-100v\r", "")
	elapsed := time.Since(start)
	fmt.Printf("time taken: %0.2fs\n", elapsed.Seconds())
	fmt.Println("counter:", counter)
	fmt.Printf("requests/sec: %0.2f\n", benchmark.ShowThroughput(start, counter))
}

// $ redis-benchmark
// time taken: 704.205418ms
// counter: 100000
