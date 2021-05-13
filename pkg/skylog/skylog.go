package skylog

import (
	"io"
	"log"
	"os"
)

//These two loggers are written so one can pass other writer opbjects to them
//for testing (to write to a buffer) and also for logging to file at some point.
//these loggers will have to be moved to a package file.

//GetInfoLogger is centeralized info logger to be used across the board.
func GetInfoLogger(out io.Writer) func() *log.Logger {
	infoLog := log.New(out, "INFO\t", log.Ldate|log.Ltime)
	return func() *log.Logger {
		return infoLog
	}
}

//GetErrorLogger is centeralized error logger to be used across the board.
func GetErrorLogger(out io.Writer) func() *log.Logger {
	errorLog := log.New(out, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return func() *log.Logger {
		return errorLog
	}
}

//InfoLog for general info logging
var InfoLog = GetInfoLogger(os.Stdout)()

//ErrorLog for general error logging
var ErrorLog = GetErrorLogger(os.Stdout)()
