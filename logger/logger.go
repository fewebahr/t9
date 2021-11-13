package logger

import (
	"log"
	"os"
)

// Logger offers various methods for logging information to the console
type Logger interface {
	SetLevel(Level)

	Debug(...interface{})
	Debugf(string, ...interface{})
	Debugln(...interface{})

	Info(...interface{})
	Infof(string, ...interface{})
	Infoln(...interface{})

	Warn(...interface{})
	Warnf(string, ...interface{})
	Warnln(...interface{})

	Error(...interface{})
	Errorf(string, ...interface{})
	Errorln(...interface{})

	GetLogger(Level) *log.Logger
}

// New returns a Logger at the designated Level.
func New(l Level) Logger {

	return &levelLogger{
		level: l,
		loggers: map[Level]*log.Logger{
			DebugLevel: log.New(os.Stdout, DebugLevel.getPrefix(), log.LstdFlags),
			InfoLevel:  log.New(os.Stdout, InfoLevel.getPrefix(), log.LstdFlags),
			WarnLevel:  log.New(os.Stdout, WarnLevel.getPrefix(), log.LstdFlags),
			ErrorLevel: log.New(os.Stdout, ErrorLevel.getPrefix(), log.LstdFlags),
		},
	}
}

type levelLogger struct {
	level   Level
	loggers map[Level]*log.Logger
}

func (ll *levelLogger) SetLevel(l Level) {

	ll.level = l
}

func (ll *levelLogger) Debug(args ...interface{}) {

	ll.print(DebugLevel, args...)
}

func (ll *levelLogger) Debugf(format string, args ...interface{}) {

	ll.printf(DebugLevel, format, args...)
}

func (ll *levelLogger) Debugln(args ...interface{}) {

	ll.println(DebugLevel, args...)
}

func (ll *levelLogger) Info(args ...interface{}) {

	ll.print(InfoLevel, args...)
}

func (ll *levelLogger) Infof(format string, args ...interface{}) {

	ll.printf(InfoLevel, format, args...)
}

func (ll *levelLogger) Infoln(args ...interface{}) {

	ll.println(InfoLevel, args...)
}

func (ll *levelLogger) Warn(args ...interface{}) {

	ll.print(WarnLevel, args...)
}

func (ll *levelLogger) Warnf(format string, args ...interface{}) {

	ll.printf(WarnLevel, format, args...)
}

func (ll *levelLogger) Warnln(args ...interface{}) {

	ll.println(WarnLevel, args...)
}

func (ll *levelLogger) Error(args ...interface{}) {

	ll.print(ErrorLevel, args...)
}

func (ll *levelLogger) Errorf(format string, args ...interface{}) {

	ll.printf(ErrorLevel, format, args...)
}

func (ll *levelLogger) Errorln(args ...interface{}) {

	ll.println(ErrorLevel, args...)
}

func (ll *levelLogger) GetLogger(l Level) *log.Logger {

	return ll.loggers[l]
}
