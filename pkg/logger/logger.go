package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	IsVerbose   bool
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime)
}

func LogError(err error) {
	if IsVerbose {
		ErrorLogger.Println(err.Error())
	}
}

func LogFatal(err error) {
	ErrorLogger.Println(err.Error())
	os.Exit(1)
}

func LogInfo(message string) {
	if IsVerbose {
		InfoLogger.Println(message)
	}
}
