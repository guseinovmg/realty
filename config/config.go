package config

import "os"

type conf struct {
	uploadedFilesPath  string
	staticFilesPath    string
	httpServerPort     string
	dbPath             string
	availableCountries []string
	language           string
}

var c conf

func Initialize() {
	c = conf{
		uploadedFilesPath:  "./uploaded/",
		staticFilesPath:    "./static/",
		httpServerPort:     ":8080",
		dbPath:             "/home/murad/haha.db",
		availableCountries: make([]string, 0),
	}
	if v, ok := os.LookupEnv("UPLOADED_FILES_PATH"); ok {
		c.uploadedFilesPath = v
	}
	if v, ok := os.LookupEnv("STATIC_FILES_PATH"); ok {
		c.staticFilesPath = v
	}
	if v, ok := os.LookupEnv("DB_PATH"); ok {
		c.dbPath = v
	}
	if v, ok := os.LookupEnv("HTTP_SERVER_PORT"); ok {
		c.httpServerPort = v
	}
}

func GetUploadedFilesPath() string {
	return c.uploadedFilesPath
}

func GetStaticFilesPath() string {
	return c.staticFilesPath
}

func GetHttpServerPort() string {
	return c.httpServerPort
}

func GetDbPath() string {
	return c.dbPath
}

func GetAvailableCountries() []string {
	return c.availableCountries
}
