package logger

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// Level expresses a logging level
type Level uint

// Available log levels
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// ParseLevel turns the designated string into a log level or returns an error.
func ParseLevel(levelStr string) (Level, error) {
	switch levelStr {
	case DebugLevel.String():
		return DebugLevel, nil
	case InfoLevel.String():
		return InfoLevel, nil
	case WarnLevel.String():
		return WarnLevel, nil
	case ErrorLevel.String():
		return ErrorLevel, nil
	default:
		return InfoLevel, errors.Errorf(`must be %s|%s|%s|%s (received '%s')`,
			DebugLevel, InfoLevel, WarnLevel, ErrorLevel, levelStr)
	}
}

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return `debug`
	case InfoLevel:
		return `info`
	case WarnLevel:
		return `warn`
	case ErrorLevel:
		return `error`
	default:
		panic(fmt.Sprintf(`unknown log level: %d`, l))
	}
}

func (l Level) getColor() *color.Color {
	switch l {
	case DebugLevel:
		return color.New(color.Reset)
	case InfoLevel:
		return color.New(color.FgCyan)
	case WarnLevel:
		return color.New(color.FgYellow)
	case ErrorLevel:
		return color.New(color.FgRed)
	default:
		panic(fmt.Sprintf(`got unexpected log level: '%s'`, l))
	}
}

func (l Level) getPrefix() string {
	str := l.String()
	paddedString := formatPadRight(str, 6)
	coloredString := l.getColor().SprintFunc()(paddedString)
	prefix := coloredString

	return prefix
}

func formatPadRight(in string, length int) string {
	numSpaces := length - len(in)
	return in + strings.Repeat(` `, numSpaces)
}
