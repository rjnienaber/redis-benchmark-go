package main

import (
	"benchmark"
	"bufio"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var wg sync.WaitGroup
	start := time.Now()
	var counter uint64
	for i := 0; i < 50; i++ {
		go func() {
			conn, err := net.Dial("tcp", "localhost:6379")
			if err != nil {
				panic(fmt.Sprintf("Couldn't connect to redis server: %v", err))
			}

			for {
				fmt.Fprintf(conn, "PING\r\n")
				result := benchmark.Parse(bufio.NewReader(conn))

				if result != "+PONG" {
					panic(fmt.Sprintf("Result should have been '+PONG' was '%v'", result))
				}

				if atomic.LoadUint64(&counter) == 100000 {
					wg.Done()
					break
				}
				atomic.AddUint64(&counter, 1)
			}
		}()
		wg.Add(1)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("time taken:", elapsed)
	fmt.Println("counter:", counter)
}

// $ redis-benchmark
// time taken: 704.205418ms
// counter: 100000
