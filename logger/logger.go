package logger

type Logger interface {
	Error(msg string)
	Errorf(format string, args ...interface{})
	Info(msg string)
	Infof(format string, args ...interface{})
}
