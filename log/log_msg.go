package log

import (
	"fmt"
	"time"
)

type _LogLevel uint8

const (
	InfoLevel  _LogLevel = 0b111
	WarnLevel  _LogLevel = 0b011
	ErrorLevel _LogLevel = 0b001
)

var level = map[string]_LogLevel{
	"error":   ErrorLevel,
	"err":     ErrorLevel,
	"warn":    WarnLevel,
	"warning": WarnLevel,
	"info":    InfoLevel,
}

func logLevel(s string) _LogLevel {
	l, ok := level[s]
	if ok {
		return l
	}
	return InfoLevel
}

func (l _LogLevel) String() string {
	switch l {
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return "unknown"
	}
}

type Msg struct {
	data     string
	dataTime time.Time
	dataType _LogLevel
}

func (m Msg) String() string {
	return fmt.Sprintf("%s %s \"%s\"\n", m.dataTime.Format("2006:01:02-15:04:05"), m.dataType, m.data)
}
