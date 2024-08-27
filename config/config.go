package config

import (
	"log"
	"log/slog"
	"os"
)

type conf struct {
	uploadedFilesPath  string
	staticFilesPath    string
	httpServerPort     string
	dataDir            string
	dataUsersPath      string
	dataAdvsPath       string
	dataWatchesPath    string
	availableCountries []string
	language           string
	domain             string
	adminId            int64
	logLevel           slog.Level
}

var c conf

func Initialize() {
	c = conf{
		uploadedFilesPath:  "/home/murad/GolandProjects/realty/uploaded/",
		staticFilesPath:    "/home/murad/GolandProjects/realty/static/",
		httpServerPort:     ":8080",
		dataDir:            ":memory:",
		availableCountries: make([]string, 0),
		domain:             "localhost",
		adminId:            35456456,
		logLevel:           slog.LevelInfo,
	}
	if v, ok := os.LookupEnv("UPLOADED_FILES_PATH"); ok {
		c.uploadedFilesPath = v
	}
	if v, ok := os.LookupEnv("STATIC_FILES_PATH"); ok {
		c.staticFilesPath = v
	}
	if v, ok := os.LookupEnv("DATA_DIR"); ok {
		c.dataDir = v
	}
	if v, ok := os.LookupEnv("HTTP_SERVER_PORT"); ok {
		c.httpServerPort = v
	}
	if v, ok := os.LookupEnv("DOMAIN"); ok {
		c.domain = v
	}
	if v, ok := os.LookupEnv("LOG_LEVEL"); ok {
		switch v {
		case "debug":
			c.logLevel = slog.LevelDebug
		case "info":
			c.logLevel = slog.LevelInfo
		case "warn":
			c.logLevel = slog.LevelWarn
		case "error":
			c.logLevel = slog.LevelError
		default:
			log.Fatal("unknown log level")
		}
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

func GetDataDir() string {
	return c.dataDir
}

func GetDbUsersPath() string {
	if c.dataDir == ":memory:" {
		return c.dataDir
	}
	return c.dataDir + "/users.sqlite"
}

func GetDbAdvsPath() string {
	if c.dataDir == ":memory:" {
		return c.dataDir
	}
	return c.dataDir + "/advs.sqlite"
}

func GetDbPhotosPath() string {
	if c.dataDir == ":memory:" {
		return c.dataDir
	}
	return c.dataDir + "/photos.sqlite"
}

func GetDbWatchesPath() string {
	if c.dataDir == ":memory:" {
		return c.dataDir
	}
	return c.dataDir + "/watches.sqlite"
}

func GetCurrencyRatesFilepath() string {
	return c.dataDir + "/currency.json"
}

func GetAvailableCountries() []string {
	return c.availableCountries
}

func GetDomain() string {
	return c.domain
}

func GetAdminId() int64 {
	return c.adminId
}

func GetLogLevel() slog.Level {
	return c.logLevel
}
