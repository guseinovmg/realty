package application

import (
	"realty/dto"
	"sync/atomic"
	"time"
)

type RequestHit struct {
	DurationNs int64
	Pattern    string
}

var instanceStartTime time.Time = time.Now()
var dbErrCount atomic.Int64
var recoveredPanicsCount atomic.Int64
var gracefullyStop atomic.Bool
var hitsChan = make(chan RequestHit, 5000)
var hitsMap = make(map[string]dto.RequestMetric)

func init() {
	go func() {
		for hit := range hitsChan {
			prevMetric := hitsMap[hit.Pattern]
			hitsMap[hit.Pattern] = dto.RequestMetric{
				Count:         prevMetric.Count + 1,
				DurationSumNs: prevMetric.DurationSumNs + hit.DurationNs,
			}
		}
	}()
}

func GracefullyStopAndExitApp() {
	if !gracefullyStop.Load() {
		gracefullyStop.Store(true)
	}
}

func IsGracefullyStopped() bool {
	return gracefullyStop.Load()
}

// GetInstanceStartTime returns the instance start time
func GetInstanceStartTime() time.Time {
	return instanceStartTime
}

func GetDbErrorsCount() int64 {
	return dbErrCount.Load()
}

func IncDbErrorCounter() {
	dbErrCount.Add(1)
}

// GetRecoveredPanicsCount returns the count of recovered panics
func GetRecoveredPanicsCount() int64 {
	return recoveredPanicsCount.Load()
}

func IncPanicCounter() {
	recoveredPanicsCount.Add(1)
}

func Hit(pattern string, durationNs int64) {
	hitsChan <- RequestHit{
		DurationNs: durationNs,
		Pattern:    pattern,
	}
}

func GetHitsMap() map[string]dto.RequestMetric {
	return hitsMap
}
