package benchmark

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Runner struct {
	Connections []net.Conn
	Options     Options
}

func (r *Runner) Execute() *Results {
	var wg sync.WaitGroup
	var counter uint64
	var limit = uint64(r.Options.Requests)

	results := NewResults(r.Options)
	for _, conn := range r.Connections {
		wg.Add(1)
		go func(c net.Conn) {
			command := fmt.Sprintf("%v\r\n", r.Options.Tests[0])
			for {
				if atomic.AddUint64(&counter, 1) > limit {
					wg.Done()
					break
				}

				start := time.Now()
				fmt.Fprint(c, command)
				_, err := Parse(bufio.NewReader(c))
				if err != nil {
					fmt.Println(err)
				}
				results.LogRun(start)
			}
		}(conn)
	}

	stopThroughput := PrintThroughput(results.Start, &counter)
	wg.Wait()
	results.Stop()
	for _, conn := range r.Connections {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing redis connection: ", err)
		}
	}

	stopThroughput <- true
	return results
}

func NewRunner(options Options) *Runner {
	runner := Runner{Options: options}
	runner.Connections = make([]net.Conn, options.Connections)
	for i := 0; i < options.Connections; i++ {
		conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", options.Host, options.Port))
		if err != nil {
			panic(fmt.Sprintf("Couldn't connect to redis server: %v", err))
		}
		fmt.Fprintf(conn, "PING\r\n")
		result, err := Parse(bufio.NewReader(conn))

		if err != nil {
			panic(err)
		}

		if result != "PONG" {
			panic(fmt.Sprintf("Result should have been 'PONG' was '%v'", result))
		}
		runner.Connections[i] = conn
	}
	return &runner
}
