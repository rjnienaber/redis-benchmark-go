package benchmark

import "time"
import "sync"

type Run struct {
	Start time.Time
	End   time.Time
}

type Results struct {
	Command       string
	Requests      int
	Start         time.Time
	Elapsed       time.Duration
	Connections   int
	ResponseTimes map[int]int
	Processor     chan Run
	WaitGroup     sync.WaitGroup
}

func processor(r *Results) {
	r.WaitGroup.Add(1)
	for run := range r.Processor {
		ms := float64(run.End.Sub(run.Start).Nanoseconds()) / 1000000.0
		if ms < 1 {
			r.ResponseTimes[1] += 1
		} else {
			r.ResponseTimes[int(ms)] += 1
		}
	}
	r.WaitGroup.Done()
}

func (r *Results) LogRun(start time.Time) {
	r.Processor <- Run{Start: start, End: time.Now()}
}

func (r *Results) Stop() {
	r.Elapsed = time.Since(r.Start)
	close(r.Processor)
	r.WaitGroup.Wait()
}

func NewResults(options Options) *Results {
	results := Results{
		Command:     "PING",
		Requests:    options.Requests,
		Connections: options.Connections,
	}
	results.ResponseTimes = make(map[int]int)
	results.Processor = make(chan Run)
	go processor(&results)
	results.Start = time.Now()
	return &results
}
