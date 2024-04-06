package log

import (
	"blue/common/filename"
	pri "blue/common/print"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	bufSize = 10000
)

var (
	blog *BlueLog

	stdio = Stdio{
		blog:   blog,
		lLevel: InfoLevel,
		null:   true,
	}

	Info = func(msg string) {
		stdio.info(msg)
	}

	Warn = func(msg string) {
		stdio.warn(msg)
	}

	Error = func(msg string) {
		stdio.err(msg)
	}
)

type BIter interface {
	info(msg string)
	err(msg string)
	warn(msg string)
}

type Stdio struct {
	blog   BIter
	lLevel _LogLevel
	null   bool
}

func (s *Stdio) info(msg string) {
	if s.lLevel == WarnLevel || s.lLevel == ErrorLevel {
		return
	}

	m := _LogMsg(msg, InfoLevel)
	if !s.null {
		s.blog.info(m.String())
	} else {
		fmt.Print(m.String())
	}

	putMsg(m)
}

func (s *Stdio) warn(msg string) {
	if s.lLevel == ErrorLevel {
		return
	}

	m := _LogMsg(msg, WarnLevel)
	if !s.null {
		s.blog.warn(m.String())
	} else {
		fmt.Print(m.String())
	}
	putMsg(m)
}

func (s *Stdio) err(msg string) {
	m := _LogMsg(msg, ErrorLevel)
	if !s.null {
		s.blog.err(m.String())
	} else {
		fmt.Print(m.String())
	}

	putMsg(m)
}

func InitLog(output string, level string, outPath string) {
	i := logLevel(level)
	blog = NewBlueLog(output, outPath)

	stdio = Stdio{
		blog:   blog,
		lLevel: i,
		null:   blog == nil,
	}

	pri.LogInitSuccess()
}

type BlueLog struct {
	of      *os.File
	logCh   chan string
	bufSize int
}

func NewBlueLog(output string, outPath string) *BlueLog {
	if output != "file" {
		return nil
	}

	b := &BlueLog{
		bufSize: bufSize,
		logCh:   make(chan string, bufSize),
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
				pri.ErrPanic(err, outPath)
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
	var msg string
	for msg = range l.logCh {
		_, err := l.of.WriteString(msg)
		if err != nil {
			return
		}
	}
}

func (l *BlueLog) info(msg string) {
	l.logCh <- msg
}

func (l *BlueLog) warn(msg string) {
	l.logCh <- msg
}

func (l *BlueLog) err(msg string) {
	l.logCh <- msg
}

func _LogMsg(msg string, level _LogLevel) *Msg {
	m := getMsg()
	m.data = msg
	m.dataType = level
	m.dataTime = time.Now()
	return m
}
