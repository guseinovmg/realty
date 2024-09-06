package metrics

import (
	"realty/cache"
	"sync/atomic"
	"time"
)

var instanceStartTime time.Time
var maxUnSavedChangesQueueCount atomic.Int64
var recoveredPanicsCount atomic.Int64

func init() {
	instanceStartTime = time.Now()
}

// GetInstanceStartTime returns the instance start time
func GetInstanceStartTime() time.Time {
	return instanceStartTime
}

// GetUnSavedChangesQueueCount returns the count of unsaved changes in the queue
func GetUnSavedChangesQueueCount() int64 {
	count := cache.GetToSaveCount()
	if count > maxUnSavedChangesQueueCount.Load() {
		maxUnSavedChangesQueueCount.Store(count)
	}
	return count
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
