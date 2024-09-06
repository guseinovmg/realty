package main

import (
	"log"
	"log/slog"
	"net/http"
	"realty/cache"
	"realty/config"
	"realty/db"
	"realty/router"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	config.Initialize()
	slog.SetLogLoggerLevel(config.GetLogLevel())
	slog.Info("START", "time", time.Now().Format("2006/01/02 15:04:05"))
	db.Initialize()
	cache.Initialize()
	mux := router.Initialize()
	log.Fatal(http.ListenAndServe(config.GetHttpServerPort(), mux))
}
