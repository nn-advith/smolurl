// Package logger handles the global logger and it's related functions.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Singular, global instance of type [logger.Logger], usable from every module.
var GlobalLogger *Logger

// A Logger writes logs to stdout and/or configured log files, at varying log levels. Wrapped over [log.Logger]
type Logger struct {
	filelog     bool
	stdlog      bool
	logfilepath string
	logfile     *os.File
	logger      *log.Logger
}

// InitLogger initialises the global logger [logger.GlobalLogger] with a combination of terminal and file logs. The required files/directories are created and a new instance of [logger.Logger] is initialised.
func InitLogger(enableStd bool, enableFile bool) error {

	var writers []io.Writer

	//log file name gen and directory checks
	var DATADIR string
	if runtime.GOOS == "windows" {
		DATADIR = filepath.Join(os.Getenv("APPDATA"), "smolurl")
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not get user home directory: %v", err)
		}
		DATADIR = filepath.Join(home, ".config", "smolurl")
	}
	if f, err := os.Stat(DATADIR); os.IsNotExist(err) {
		//directory not present; create
		err := os.MkdirAll(DATADIR, 0755)
		if err != nil {
			return fmt.Errorf("unable to create data directory")
		}
	} else if err != nil {
		return fmt.Errorf("error during logger init - %v", err)
	} else if !f.IsDir() {
		return fmt.Errorf("%v is not a directory", DATADIR)
	}
	var logfilepath string
	//create logs directory
	fmt.Println(DATADIR)
	err := os.MkdirAll(filepath.Join(DATADIR, "logs"), 0755)
	if err != nil {
		return fmt.Errorf("unable to create logs directory")
	}
	if enableFile {
		logfilepath = filepath.Join(DATADIR, "logs/log-"+strings.Split(time.Now().Format(time.RFC3339), "T")[0]+".log")
	}

	newlogger := &Logger{
		filelog:     enableFile,
		stdlog:      enableStd,
		logfilepath: logfilepath,
	}

	if enableFile {
		file, err := os.OpenFile(newlogger.logfilepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("error creating log file - %v", err)
		}
		newlogger.logfile = file
		writers = append(writers, file)
	}

	if enableStd {
		writers = append(writers, os.Stdout)
	}

	multi := io.MultiWriter(writers...)
	newlogger.logger = log.New(multi, "", log.LstdFlags|log.Lshortfile)
	GlobalLogger = newlogger
	return nil
}

// Close closes the logfile that is used for file logging. Must be closed manually in cleanup or atleast deferred post initialisation.
func (l *Logger) Close() { //must be deferred. file descriptor is kept open for performance
	if l.filelog && l.logfile != nil {
		l.logfile.Close()
	}
}

// Info logs messages at info log level.
func (l *Logger) Info(msg ...any) {
	l.logger.Println(append([]any{"[INFO]:"}, msg...)...)
}

// Error logs messages at error log level.
func (l *Logger) Error(msg ...any) {
	l.logger.Println(append([]any{"[ERROR]:"}, msg...)...)
}

// Fatal logs messages at fatal log level and exits.
func (l *Logger) Fatal(msg ...any) {
	l.logger.Fatalln(append([]any{"[FATAL]:"}, msg...)...)
}
