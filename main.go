package main

import (
	"log"
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
	router.Initialize()
}
