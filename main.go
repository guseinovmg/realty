package main

import (
	"log"
	"log/slog"
	"net/http"
	"realty/cache"
	"realty/config"
	"realty/db"
	"realty/router"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	config.Initialize()
	slog.SetLogLoggerLevel(config.GetLogLevel())
	db.Initialize()
	cache.Initialize()
	mux := router.Initialize()
	log.Fatal(http.ListenAndServe(config.GetHttpServerPort(), mux))
}
