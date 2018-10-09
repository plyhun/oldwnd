package log

import (
	"github.com/golang/glog"
)

type glogLogger struct {
}

func (this *glogLogger) Init() error {
	return nil
}

func (this *glogLogger) Tracef(format string, args ...interface{}) {
	if glog.V(5) {
		ilog(format+"\n", args...)
	}
}
func (this *glogLogger) Infof(format string, args ...interface{}) {
	if glog.V(4) {
		ilog(format+"\n", args...)
	}
}
func (this *glogLogger) Debugf(format string, args ...interface{}) {
	if glog.V(3) {
		ilog(format+"\n", args...)
	}
}
func (this *glogLogger) Warnf(format string, args ...interface{}) {
	if glog.V(2) {
		ilog(format+"\n", args...)
	}
}
func (this *glogLogger) Errorf(format string, args ...interface{}) {
	if glog.V(1) {
		ilog(format+"\n", args...)
	}
}
func (this *glogLogger) Fatalf(format string, args ...interface{}) {
	if glog.V(0) {
		ilog(format+"\n", args...)
	}
}

func ilog(format string, args ...interface{}) {
	glog.Infof(format+"\n", args...)
}

func (this *glogLogger) Close() {

}
