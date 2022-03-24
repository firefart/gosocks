package socks

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
}
