package benchmark

import (
	"time"
)

func ShowThroughput(start time.Time, counter uint64) float64 {
	elapsed := time.Now().Sub(start)
	return float64(counter) / elapsed.Seconds()
}
