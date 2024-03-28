package log

import (
	"blue/common/rand"
	"blue/config"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

var (
	blog *blueLog

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

var level = map[string]zerolog.Level{
	"Error": zerolog.ErrorLevel,
	"warn":  zerolog.WarnLevel,
	"info":  zerolog.InfoLevel,
}

func logLevel() zerolog.Level {
	l, ok := level[config.LogCfg.LogLevel]
	if !ok {
		return zerolog.InfoLevel
	}

	return l
}

func InitLog() {
	blog = newLog()
	config.LogInitSuccess()
}

func newLog() *blueLog {
	dir := config.LogCfg.LogOut

	logx := newZeroLog(logLevel(), dir)
	return logx
}

type blueLog struct {
	l zerolog.Logger
}

func newZeroLog(level zerolog.Level, outPath string) *blueLog {
	zerolog.TimestampFieldName = "T"
	zerolog.MessageFieldName = "M"
	zerolog.LevelFieldName = "L"

	openFile := &os.File{}
	if filepath.Ext(outPath) == ".log" {
		var err error
		openFile, err = os.Open(outPath)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		if _, err = os.Stat(outPath); os.IsNotExist(err) {
			err = os.MkdirAll(outPath, 0777)
			if err != nil {
				config.ErrPanic(err, outPath)
			}
		}
		r := rand.RandString(8)
		addtime := fmt.Sprintf("%s-%s", time.Now().Format("2006:01:02-15:04:05"), r)
		openFile, err = os.Create(fmt.Sprintf("%s/%s.log", outPath, addtime))
		if err != nil {
			panic(err)
		}
	}

	logger := zerolog.New(openFile).Level(level).With().Timestamp().Logger()

	return &blueLog{
		l: logger,
	}
}

func (l *blueLog) info(msg string) {
	l.l.Info().Msg(msg)
}

func (l *blueLog) warn(msg string) {
	l.l.Warn().Msg(msg)
}

func (l *blueLog) err(msg string) {
	l.l.Error().Msg(msg)
}
