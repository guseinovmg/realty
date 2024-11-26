package api

import "sync/atomic"

var gracefullyStop atomic.Bool

func GracefullyStopAndExitApp() {
	if !gracefullyStop.Load() {
		gracefullyStop.Store(true)
	}
}

func IsGracefullyStopped() bool {
	return gracefullyStop.Load()
}
