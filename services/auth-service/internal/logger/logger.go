package logger

import (
	"log"
	"os"
	"sync"
)

var (
	defaultLogger *log.Logger
	once          sync.Once
)

// L returns a shared logger with a consistent prefix for the auth service.
func L() *log.Logger {
	once.Do(func() {
		defaultLogger = log.New(os.Stdout, "[auth-service] ", log.LstdFlags|log.Lshortfile)
	})
	return defaultLogger
}
