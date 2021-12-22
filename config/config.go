package config

import (
	_ "github.com/lib/pq" // we don't neeed de pg variable
	log "github.com/sirupsen/logrus"
)

// Configuration ...
type Configuration struct {
	JWTSecret   string `env:"JWT_SECRET,default=asdf"`
	Port        int    `env:"PORT,default=3001"`
	Timeout     int    `env:"TIMEOUT,default=100"`
	Environment string `env:"ENVIRONMENT,default=production"`
}

// AppContext keeps a reference of the common objects for the application
// any http handler will have access to this context.
type AppContext struct {
	Config         Configuration
	CloseGoRoutine chan bool
}

var logger *log.Entry

// SetLogger ...
func SetLogger(newLogger *log.Entry) {
	logger = newLogger
}

// GetLogger ...
func GetLogger() *log.Entry {
	return logger
}
