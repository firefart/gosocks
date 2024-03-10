package socks

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type NilLogger struct {
}

func (l *NilLogger) Debug(args ...interface{})                 {}
func (l *NilLogger) Debugf(format string, args ...interface{}) {}
func (l *NilLogger) Info(args ...interface{})                  {}
func (l *NilLogger) Infof(format string, args ...interface{})  {}
func (l *NilLogger) Warn(args ...interface{})                  {}
func (l *NilLogger) Warnf(format string, args ...interface{})  {}
func (l *NilLogger) Error(args ...interface{})                 {}
func (l *NilLogger) Errorf(format string, args ...interface{}) {}
