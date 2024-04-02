package log

import (
	"blue/common/filename"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"blue/config"
)

const (
	bufSize  = 10000
	flushDef = EveySec
)

var (
	blog *BlueLog

	Info = func(msg string) {
		blog.info(msg)
	}

	Warn = func(msg string) {
		blog.warn(msg)
	}

	Error = func(msg string) {
		blog.err(msg)
	}
)

func InitLog(level string, outPath string) {
	i := logLevel(level)
	blog = NewLog(i, outPath)
	config.LogInitSuccess()
}

type BlueLog struct {
	of        *os.File
	lLevel    _LogLevel
	logCh     chan *Msg
	bufSize   int
	flushType FlushType
	tk        *time.Ticker
	nowWrite  chan struct{}
}

func NewLog(level _LogLevel, outPath string) *BlueLog {
	b := &BlueLog{
		flushType: flushDef,
		bufSize:   bufSize,
		logCh:     make(chan *Msg, bufSize),
		tk:        time.NewTicker(time.Second),
		nowWrite:  make(chan struct{}),
		lLevel:    level,
	}

	var err error
	if filepath.Ext(outPath) == ".log" {
		b.of, err = os.Open(outPath)
		if err != nil {
			panic(err)
		}

	} else {
		if _, err = os.Stat(outPath); os.IsNotExist(err) {
			err = os.MkdirAll(outPath, os.ModePerm)
			if err != nil {
				config.ErrPanic(err, outPath)
			}
		}

		b.of, err = os.Create(fmt.Sprintf("%s/%s.log", outPath, filename.LogName()))
		if err != nil {
			panic(err)
		}
	}

	go b.sync()

	return b
}

func (l *BlueLog) sync() {
	var msg *Msg
	for msg = range l.logCh {
		_, err := l.of.WriteString(msg.String())
		if err != nil {
			return
		}

		putMsg(msg)
	}
}

func (l *BlueLog) info(msg string) {
	if l.lLevel == WarnLevel || l.lLevel == ErrorLevel {
		return
	}

	m := l.getMsg()
	m.data = msg
	m.dataType = InfoLevel
	l.logCh <- m
}

func (l *BlueLog) warn(msg string) {
	if l.lLevel == ErrorLevel {
		return
	}

	m := l.getMsg()
	m.data = msg
	m.dataType = WarnLevel
	l.logCh <- m
}

func (l *BlueLog) err(msg string) {
	m := l.getMsg()
	m.data = msg
	m.dataType = ErrorLevel
	l.logCh <- m
}

func (l *BlueLog) getMsg() *Msg {
	m := getMsg()
	m.dataTime = time.Now()
	return m
}
