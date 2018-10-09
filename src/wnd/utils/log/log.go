package log

import (
	"runtime"
	"strings"
	"strconv"
	
	"github.com/kr/pretty"
)	

type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
	LevelOff 
)

type wndlogger interface {
	Init() error

	Tracef(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Close()
}

var logger wndlogger
var level Level

var filter string

func Filter(f string) {
	filter = f
}

func LoggerLevel(l Level) {
	level = l
}

func prettifyArgs(args []interface{}) []interface{} {
	for i, a := range args {
		args[i] = pretty.Formatter(a)
	}

	return args
}

func Tracef(format string, args ...interface{}) {
	if level > LevelTrace {
		return
	}
	
	pc, _, _, _ := runtime.Caller(1)
	
	if filter != "" && strings.Index(runtime.FuncForPC(pc).Name(), filter) < 0{
		return
	}
	
	logger.Tracef("[" + runtime.FuncForPC(pc).Name() + "] " + format, prettifyArgs(args)...)
}

func Infof(format string, args ...interface{}) {
	if level > LevelInfo {
		return
	}
	
	pc, _, _, _ := runtime.Caller(1)
	
	if filter != "" && strings.Index(runtime.FuncForPC(pc).Name(), filter) < 0{
		return
	}
	
	logger.Infof("[" + runtime.FuncForPC(pc).Name() + "] " + format, prettifyArgs(args)...)
}

func Debugf(format string, args ...interface{}) {
	if level > LevelDebug {
		return
	}
	
	pc, _, _, _ := runtime.Caller(1)
	
	if filter != "" && strings.Index(runtime.FuncForPC(pc).Name(), filter) < 0{
		return
	}
	
	logger.Debugf("[" + runtime.FuncForPC(pc).Name() + "] " + format, prettifyArgs(args)...)
}

func Warnf(format string, args ...interface{}) {
	if level > LevelWarning {
		return
	}
	
	pc, f, l, _ := runtime.Caller(1)
	
	if filter != "" && strings.Index(runtime.FuncForPC(pc).Name(), filter) < 0{
		return
	}
	
	logger.Warnf("[" + runtime.FuncForPC(pc).Name() + "] [" + f + " at " + strconv.Itoa(l) + "] " + format, prettifyArgs(args)...)
}

func Errorf(format string, args ...interface{}) {
	if level > LevelError {
		return
	}
	
	pc, f, l, _ := runtime.Caller(1)
	
	if filter != "" && strings.Index(runtime.FuncForPC(pc).Name(), filter) < 0{
		return
	}
	
	logger.Errorf("[" + runtime.FuncForPC(pc).Name() + "] [" + f + " at " + strconv.Itoa(l) + "] " + format, prettifyArgs(args)...)
}

func Fatalf(format string, args ...interface{}) {
	if level > LevelFatal {
		return
	}
	
	pc, f, l, _ := runtime.Caller(1)
	
	if filter != "" && strings.Index(runtime.FuncForPC(pc).Name(), filter) < 0{
		return
	}
	
	logger.Fatalf("[" + runtime.FuncForPC(pc).Name() + "] [" + f + " at " + strconv.Itoa(l) + "] " + format, prettifyArgs(args)...)
}

func init() {
	logger = new(seelogLogger)
	//logger = new(defaultLogger)
	//logger = new(glogLogger)
	
	logger.Init()
}
