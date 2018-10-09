package log

import (
	"testing"
	"fmt"
)

type testLogger struct {
	t *testing.T
}

func (this *testLogger) Init() (err error) {
	return nil
}

func (this *testLogger) Tracef(format string, args ...interface{}) {
	this.t.Logf("[TRACE] %s\n", fmt.Sprintf(format, prettifyArgs(args)...))
}

func (this *testLogger) Infof(format string, args ...interface{}) {
	this.t.Logf("[INFO] %s\n", fmt.Sprintf(format, prettifyArgs(args)...))
}

func (this *testLogger) Debugf(format string, args ...interface{}) {
	this.t.Logf("[DEBUG] %s\n", fmt.Sprintf(format, prettifyArgs(args)...))
}

func (this *testLogger) Warnf(format string, args ...interface{}) {
	this.t.Logf("[WARN] %s\n", fmt.Sprintf(format, prettifyArgs(args)...))
}

func (this *testLogger) Errorf(format string, args ...interface{}) {
	this.t.Errorf("[ERROR] %s\n", fmt.Sprintf(format, prettifyArgs(args)...))
}

func (this *testLogger) Fatalf(format string, args ...interface{}) {
	this.t.Fatalf("[FATAL] %s\n", fmt.Sprintf(format, prettifyArgs(args)...))
}

func (this *testLogger) Close() {
	
}

func NewTestLogger(t *testing.T) {
	logger = &testLogger{t: t}
}