package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

// Define your custom logger type.
type Logger struct {
	out      io.Writer // destination for output
	level    int
	handlers map[outputType]logHandler // list of log handlers
}

type logHandler struct {
	handlerType outputType
	logfile     string
}

type outputType int

const (
	DEBUG = 1 << iota
	INFO
	WARNING
	ERROR
)

const (
	STDOUT outputType = iota
	FILE
)

var defaultLogger *Logger

func init() {
	var defaultHandler logHandler
	defaultHandler.handlerType = STDOUT

	defaultLogger = new(Logger)
	defaultLogger.level = INFO

	defaultLogger.handlers = make(map[outputType]logHandler)
	defaultLogger.handlers[STDOUT] = defaultHandler
}

func Error(vars ...interface{}) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(defaultLogger.out)
	if defaultLogger.level <= ERROR {
		logMessage(fmt.Sprintf("[ERROR] %v", vars...))
	}
}

func Info(vars ...interface{}) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(defaultLogger.out)
	if defaultLogger.level <= INFO {
		logMessage(fmt.Sprintf("[INFO] %v", vars...))
	}
}

func Debug(vars ...interface{}) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(defaultLogger.out)
	if defaultLogger.level <= DEBUG {
		logMessage(fmt.Sprintf("[DEBUG] %v", vars...))
	}
}

func Warning(vars ...interface{}) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(defaultLogger.out)
	if defaultLogger.level <= WARNING {
		logMessage(fmt.Sprintf("[WARNING] %v", vars...))
	}
}

func SetLevel(level interface{}) {
	if reflect.TypeOf(level).String() == "int" {
		defaultLogger.level = level.(int)
	} else if reflect.TypeOf(level).String() == "string" {
		switch level.(string) {
		case strings.ToUpper("debug"):
			defaultLogger.level = 1
		case strings.ToUpper("info"):
			defaultLogger.level = 2
		case strings.ToUpper("warning"):
			defaultLogger.level = 3
		default:
			defaultLogger.level = 4
		}
	}
}

func SetFileHandler(logfile string) {
	var fileHandler logHandler
	fileHandler.handlerType = FILE

	if len(logfile) > 0 {
		fileHandler.logfile = logfile
	} else {
		fileHandler.logfile = "logs/application.log"
	}

	defaultLogger.handlers[FILE] = fileHandler
	SetLogger(defaultLogger.level)
}

func SetLogger(level interface{}) {
	SetLevel(level)

	var handlers []io.Writer
	for _, v := range defaultLogger.handlers {
		switch v.handlerType {
		case STDOUT:
			handlers = append(handlers, os.Stdout)
		case FILE:
			logfile := createFolderIfNotExists(v.logfile)
			file, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				Error("Error opening file: %v", err)
			}
			handlers = append(handlers, file)
			//Info(logfile, " is attached as a log file")
		}
	}
	defaultLogger.out = io.MultiWriter(handlers...)
	log.SetOutput(defaultLogger.out)

	if reflect.TypeOf(level).String() == "int" {
		defaultLogger.level = level.(int)
	} else if reflect.TypeOf(level).String() == "string" {
		switch level.(string) {
		case strings.ToUpper("debug"):
			defaultLogger.level = 1
			Info("Current logger level set to: DEBUG")
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

func createFolderIfNotExists(path string) string {
	mPath := strings.Split(path, "/")
	wd, err := os.Getwd()

	if len(mPath) > 1 {
		path = wd + "/" + strings.Join(mPath[:len(mPath)-1], "/") + "/"
	} else {
		return wd
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
	}

	return path + mPath[len(mPath)-1]
}

func logMessage(message string) {
	if handler, ok := defaultLogger.handlers[FILE]; ok {
		logFile, err := os.OpenFile(handler.logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			Error("Error opening file: %v", err)
		}
		if err != nil {
			Error("Error writing to file: %v", err)
		}
		defer logFile.Close()
	}
	_ = log.Output(2, message)
}
