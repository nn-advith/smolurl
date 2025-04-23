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

var GlobalLogger *Logger

type Logger struct {
	filelog     bool
	stdlog      bool
	logfilepath string // no control on this. will be auto-generated based on date
	logfile     *os.File
	logger      *log.Logger
}

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

func (l *Logger) Close() { //must be deferred. file descriptor is kept open for performance
	if l.filelog && l.logfile != nil {
		l.logfile.Close()
	}
}

func (l *Logger) Info(msg ...any) {
	l.logger.Println(append([]any{"[INFO]:"}, msg...)...)
}

func (l *Logger) Error(msg ...any) {
	l.logger.Println(append([]any{"[ERROR]:"}, msg...)...)
}

func (l *Logger) Fatal(msg ...any) {
	l.logger.Fatalln(append([]any{"[FATAL]:"}, msg...)...)
}
