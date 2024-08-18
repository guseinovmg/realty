package main

import (
	"log"
	"net/http"
	"realty/cache"
	"realty/config"
	"realty/db"
	"realty/router"
)

func main() {
	log.SetFlags(log.Lshortfile)
	config.Initialize()
	db.Initialize()
	cache.Initialize()
	mux := router.Initialize()
	log.Fatal(http.ListenAndServe(config.GetHttpServerPort(), mux))
}
