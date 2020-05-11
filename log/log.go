package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"sync"

	"github.com/fatih/color"
)

type Handler func(format string, args ...interface{})

const (
	bold = color.Bold
	dim  = color.Faint

	fgRed   = color.FgRed
	fgWhite = color.FgWhite
	fgBlack = color.FgBlack

	bgDgray  = color.BgHiBlack
	bgRed    = color.BgRed
	bgGreen  = color.BgGreen
	bgYellow = color.BgYellow

	reset = color.Reset
)

const (
	DEBUG = iota
	STAT
	INFO
	WARNING
	ERROR
	FATAL
)

var (
	WithColors = true
	Output     = color.Output
	DateFormat = "Mon Jan _2 15:04:05 2006"
	MinLevel   = DEBUG

	bg = false

	mutex  = &sync.Mutex{}
	labels = map[int]string{
		DEBUG:   "DBG",
		INFO:    "INF",
		WARNING: "WAR",
		ERROR:   "ERR",
		STAT:    "STAT",
		FATAL:   "!!!",
	}

	colors = map[int]*color.Color{
		DEBUG:   color.New(dim, fgBlack, bgDgray),
		STAT:    color.New(dim, fgBlack, bgDgray),
		INFO:    color.New(fgWhite, bgGreen),
		WARNING: color.New(fgWhite, bgYellow),
		ERROR:   color.New(fgWhite, bgRed),
		FATAL:   color.New(fgWhite, bgRed, bold),
	}

	Options LoggingOptions
	file    *os.File = nil
)

type LoggingOptions struct {
	LogRequestPath string
}

func Init(debug bool, background bool, file string) {
	bg = background

	if debug == true {
		MinLevel = INFO
	} else {
		MinLevel = ERROR
	}

	Options = LoggingOptions{
		LogRequestPath: file,
	}
}

func Wrap(s string, effect *color.Color) string {
	if WithColors == true {
		s = effect.Sprint(s)
	}
	return s
}

func Dim(s string) string {
	return Wrap(s, color.New(dim))
}

func Log(level int, format string, args ...interface{}) {
	if level >= MinLevel && !bg {
		label := labels[level]
		colorL := colors[level]
		when := time.Now().UTC().Format(DateFormat)

		what := fmt.Sprintf(format, args...)
		if strings.HasSuffix(what, "\n") == false {
			what += "\n"
		}

		l := Dim("[%s]")
		r := Wrap(" %s ", colorL) + " %s"

		//save console write from chaos
		mutex.Lock()
		fmt.Fprintf(Output, l+" "+r, when, label, what)
		mutex.Unlock()
	}
}

func Debugf(format string, args ...interface{}) {
	Log(DEBUG, format, args...)

}

func Infof(format string, args ...interface{}) {
	Log(INFO, format, args...)
}

func Warningf(format string, args ...interface{}) {
	Log(WARNING, format, args...)
}

func Errorf(format string, args ...interface{}) {
	Log(ERROR, format, args...)
}

func Statf(format string, args ...interface{}) {
	Log(STAT, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Log(FATAL, format, args...)
	os.Exit(1)
}

func Fatal(err error) {
	Log(FATAL, "%s", err)
	os.Exit(1)
}

func LogRequestFile(data string) {
	if Options.LogRequestPath != "" {
		if file == nil {
			file, _ = os.OpenFile(Options.LogRequestPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		}

		//add time
		data = time.Now().UTC().Format(DateFormat) + " " + data + "\n"
		if _, err := file.Write([]byte(data)); err != nil {
			Debugf("[Log File] %s", err)
		}
	}
}
