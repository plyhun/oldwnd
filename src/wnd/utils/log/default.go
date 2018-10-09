package log

import (
	"fmt"
	"log"
)

type defaultLogger struct {
	
}

func (this *defaultLogger) Init() (err error) {
	return nil
}

func (this *defaultLogger) Tracef(format string, args ...interface{}) {
	log.Printf("[TRACE] %s\n", fmt.Sprintf(format, prettifyArgs(args)...))
}

func (this *defaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] %s\n", fmt.Sprintf(format, prettifyArgs(args)...))
}

func (this *defaultLogger) Debugf(format string, args ...interface{}) {
	log.Printf("[DEBUG] %s\n", fmt.Sprintf(format, prettifyArgs(args)...))
}

func (this *defaultLogger) Warnf(format string, args ...interface{}) {
	log.Printf("[WARN] %s\n", fmt.Sprintf(format, prettifyArgs(args)...))
}

func (this *defaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] %s\n", fmt.Sprintf(format, prettifyArgs(args)...))
}

func (this *defaultLogger) Fatalf(format string, args ...interface{}) {
	log.Printf("[FATAL] %s\n", fmt.Sprintf(format, prettifyArgs(args)...))
}

func (this *defaultLogger) Close() {
	
}