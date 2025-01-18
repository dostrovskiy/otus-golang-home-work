package logger

import "fmt"

var (
	debugSet = map[string]struct{}{"DEBUG": {}}
	infoSet  = map[string]struct{}{"INFO": {}}
	warnSet  = map[string]struct{}{"WARN": {}}
	errorSet = map[string]struct{}{"ERROR": {}}
)

type Logger struct {
	level string
}

func init() {
	addSet(debugSet, infoSet)
	addSet(infoSet, warnSet)
	addSet(warnSet, errorSet)
}

func addSet(from map[string]struct{}, to map[string]struct{}) {
	for key := range from {
		to[key] = struct{}{}
	}
}

func New(level string) *Logger {
	return &Logger{level: level}
}

func (l Logger) write(logSet map[string]struct{}, format string, a ...any) {
	if _, ok := logSet[l.level]; ok {
		fmt.Printf(format+"\n", a...)
	}
}

func (l Logger) Error(format string, a ...any) {
	l.write(errorSet, format, a...)
}

func (l Logger) Warn(format string, a ...any) {
	l.write(warnSet, format, a...)
}

func (l Logger) Info(format string, a ...any) {
	l.write(infoSet, format, a...)
}

func (l Logger) Debug(format string, a ...any) {
	l.write(debugSet, format, a...)
}
