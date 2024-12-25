package application

import (
	"realty/dto"
	"sync/atomic"
	"time"
)

type RequestHit struct {
	DurationNs int64
	StatusCode int
	Pattern    string
}

var instanceStartTime time.Time = time.Now()
var dbErrCount atomic.Int64
var recoveredPanicsCount atomic.Int64
var gracefullyStop atomic.Bool
var gracefullyStopTime time.Time
var hitsChan = make(chan RequestHit, 5000)
var hitsMap = make(map[string]map[int]dto.RequestMetric)

func init() {
	go func() {
		for hit := range hitsChan {
			if hitsMap[hit.Pattern] == nil {
				hitsMap[hit.Pattern] = make(map[int]dto.RequestMetric)
			}
			prevMetric := hitsMap[hit.Pattern][hit.StatusCode]
			newCount := prevMetric.Count + 1
			newDuration := prevMetric.DurationSumNs + hit.DurationNs
			hitsMap[hit.Pattern][hit.StatusCode] = dto.RequestMetric{
				Count:         newCount,
				DurationSumNs: newDuration,
				AvgNs:         newDuration / newCount,
			}
		}
	}()
}

func GracefullyStopAndExitApp() {
	if !gracefullyStop.Load() {
		gracefullyStopTime = time.Now()
		gracefullyStop.Store(true)
	}
}

func IsGracefullyStopped() bool {
	return gracefullyStop.Load()
}

func GetGracefullyStopTime() *string {
	if gracefullyStopTime.IsZero() {
		return nil
	}
	s := gracefullyStopTime.Format("2006/01/02 15:04:05")
	return &s
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

func Hit(pattern string, status int, durationNs int64) {
	hitsChan <- RequestHit{
		DurationNs: durationNs,
		Pattern:    pattern,
		StatusCode: status,
	}
}

func GetHitsMap() map[string]map[int]dto.RequestMetric {
	return hitsMap
}
