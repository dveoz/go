package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
)

// Define your custom logger type.
type Logger struct {
	mu     sync.Mutex // ensures atomic writes; protects the following fields
	prefix string     // prefix to write at beginning of each line
	flag   int        // properties
	out    io.Writer  // destination for output
	buf    []byte     // for accumulating text to write
	level  int        // One of DEBUG, ERROR, INFO
}

const (
	DEBUG = 1 << iota
	INFO
	WARNING
	ERROR
)

var defaultLogger *Logger

func init() {
	defaultLogger = new(Logger)
}

func Error(vars ...interface{}) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(defaultLogger.out)
	if defaultLogger.level <= ERROR {
		message := fmt.Sprintf("[ERROR  ] %v", vars...)
		_ = log.Output(2, message)
	}
}

func Info(format string, vars ...interface{}) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(defaultLogger.out)
	if defaultLogger.level <= INFO {
		message := fmt.Sprintf("[INFO   ] "+format, vars...)
		_ = log.Output(2, message)
	}
}

func Debug(vars ...interface{}) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(defaultLogger.out)
	if defaultLogger.level <= DEBUG {
		message := fmt.Sprintf("[DEBUG  ] %v", vars...)
		_ = log.Output(2, message)
	}
}

func Warning(vars ...interface{}) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(defaultLogger.out)
	if defaultLogger.level <= WARNING {
		message := fmt.Sprintf("[WARNING] %v", vars...)
		_ = log.Output(2, message)
	}
}

func SetLogger(level interface{}, logfile string, logpath string) {
	// open log file for writing in case filename set up in config
	if len(logfile) > 0 {
		logpath = checkLogsFolder(logpath)
		logFile, err := os.OpenFile(logpath+logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			Error("Error opening file: %v", err)
		}
		defer logFile.Close()
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		defaultLogger.out = multiWriter
		log.SetOutput(defaultLogger.out)
	}

	if reflect.TypeOf(level).String() == "int" {
		defaultLogger.level = level.(int)
	} else if reflect.TypeOf(level).String() == "string" {
		switch level.(string) {
		case strings.ToUpper("debug"):
			defaultLogger.level = 1
			Debug("Current logger level set to: DEBUG")
		case strings.ToUpper("info"):
			defaultLogger.level = 2
			Info("Current logger level set to: INFO")
		case strings.ToUpper("warning"):
			defaultLogger.level = 3
			Info("Current logger level set to: WARNING")
		default:
			defaultLogger.level = 4
			Info("Current logger level set to: ERROR")
		}
	}
}

func checkLogsFolder(path string) string {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
	}
	if path[len(path)-1:len(path)] != "/" {
		path += "/"
	}
	return path
}
