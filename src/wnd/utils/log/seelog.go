package log

import (
	"fmt"

	"github.com/cihub/seelog"
)

type seelogLogger struct {
	logger seelog.LoggerInterface
}

func (this *seelogLogger) Init() (err error) {
	logConfig := "<seelog type=\"sync\" minlevel=\"trace\"><outputs><console/></outputs></seelog>"
	this.logger, err = seelog.LoggerFromConfigAsString(logConfig)
	if err == nil {
		seelog.ReplaceLogger(this.logger)
	} else {
		fmt.Errorf("Cannot init logger: %#v", err)
	}

	return err
}

func (this *seelogLogger) Tracef(format string, args ...interface{}) {
	seelog.Tracef(format, prettifyArgs(args)...)
	seelog.Flush()
}

func (this *seelogLogger) Infof(format string, args ...interface{}) {
	seelog.Infof(format, prettifyArgs(args)...)
	seelog.Flush()
}

func (this *seelogLogger) Debugf(format string, args ...interface{}) {
	seelog.Debugf(format, prettifyArgs(args)...)
	seelog.Flush()
}

func (this *seelogLogger) Warnf(format string, args ...interface{}) {
	seelog.Warnf(format, prettifyArgs(args)...)
	seelog.Flush()
}

func (this *seelogLogger) Errorf(format string, args ...interface{}) {
	seelog.Errorf(format, prettifyArgs(args)...)
	seelog.Flush()
}

func (this *seelogLogger) Fatalf(format string, args ...interface{}) {
	seelog.Criticalf(format, prettifyArgs(args)...)
	seelog.Flush()
}

func (this *seelogLogger) Close() {
	this.logger.Close()
}
