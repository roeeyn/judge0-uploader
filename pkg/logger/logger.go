package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func LogError(verbose bool, err error) {
	if verbose {
		ErrorLogger.Println(err.Error())
	}
}

func LogFatal(err error) {
	ErrorLogger.Println(err.Error())
	os.Exit(1)
}

func LogInfo(verbose bool, message string) {
	if verbose {
		InfoLogger.Println(message)
	}
}
