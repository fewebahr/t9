package logger

func (ll *levelLogger) print(l Level, args ...interface{}) {
	if ll.level > l {
		return
	}
	ll.loggers[l].Print(args...)
}

func (ll *levelLogger) printf(l Level, format string, args ...interface{}) {
	if ll.level > l {
		return
	}
	ll.loggers[l].Printf(format, args...)
}

func (ll *levelLogger) println(l Level, args ...interface{}) {
	if ll.level > l {
		return
	}
	ll.loggers[l].Println(args...)
}
