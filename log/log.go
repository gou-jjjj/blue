package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"blue/common/rand"
	"blue/config"

	"github.com/rs/zerolog"
)

var (
	blog *blueLog

	Info = func(msg ...string) {
		blog.info(msg...)
	}

	Warn = func(msg ...string) {
		blog.warn(msg...)
	}

	Error = func(msg ...string) {
		blog.err(msg...)
	}
)

var level = map[string]zerolog.Level{
	"error":   zerolog.ErrorLevel,
	"err":     zerolog.ErrorLevel,
	"warn":    zerolog.WarnLevel,
	"warning": zerolog.WarnLevel,
	"info":    zerolog.InfoLevel,
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

func (l *blueLog) info(msg ...string) {
	if len(msg) == 0 {
	} else if len(msg) == 1 {
		l.l.Info().Msg(msg[0])
	} else {
		l.l.Info().Msgf("%s,%s", msg[0], msg[1])
	}
}

func (l *blueLog) warn(msg ...string) {
	if len(msg) == 0 {
	} else if len(msg) == 1 {
		l.l.Warn().Msg(msg[0])
	} else {
		l.l.Warn().Msgf("%s,%s", msg[0], msg[1])
	}
}

func (l *blueLog) err(msg ...string) {
	if len(msg) == 0 {
	} else if len(msg) == 1 {
		l.l.Error().Msg(msg[0])
	} else {
		l.l.Error().Msgf("%s,%s", msg[0], msg[1])
	}
}
