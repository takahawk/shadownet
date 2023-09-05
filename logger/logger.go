package logger

type Logger interface {
	Info(msg string)
	Infof(format string, args ...interface{})
}
