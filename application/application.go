package application

import (
	"sync/atomic"
	"time"
)

var instanceStartTime time.Time
var maxUnSavedChangesQueueCount atomic.Int64
var dbErrCount atomic.Int64
var recoveredPanicsCount atomic.Int64
var gracefullyStop atomic.Bool

func init() {
	instanceStartTime = time.Now()
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

// GetMaxUnSavedChangesQueueCount returns the maximum count of unsaved changes in the queue
func GetMaxUnSavedChangesQueueCount() int64 {
	return maxUnSavedChangesQueueCount.Load()
}

func SetMaxUnSavedChangesQueueCount(count int64) {
	maxUnSavedChangesQueueCount.Store(count)
}
