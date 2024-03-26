package log

import (
	"blue/common/rand"
	"blue/config"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	lo "log"
)

var (
	Blog *BlueLog
)

var LogLevel = map[string]zerolog.Level{
	"Error": zerolog.ErrorLevel,
	"Warn":  zerolog.WarnLevel,
	"Info":  zerolog.InfoLevel,
}

func logLevel() zerolog.Level {
	l, ok := LogLevel[config.BC.LogConfig.LogLevel]
	if !ok {
		return zerolog.InfoLevel
	}

	return l
}

func InitLog() {
	Blog = newSyncLog()
}

func newSyncLog() *BlueLog {
	dir := config.BC.LogConfig.LogOut

	logx := NewZeroLog(logLevel(), dir)
	lo.Printf("log init success ...")
	return logx
}

type BlueLog struct {
	l zerolog.Logger
}

func NewZeroLog(level zerolog.Level, outPath string) *BlueLog {
	zerolog.TimestampFieldName = "T"
	zerolog.MessageFieldName = "M"
	zerolog.LevelFieldName = "L"

	err := os.MkdirAll(outPath, 0777)
	if err != nil {
		panic(err)
	}

	r := rand.RandString(8)
	addtime := fmt.Sprintf("%s-%s", time.Now().Format("2006:01:02-15:04:05"), r)
	create, err := os.Create(fmt.Sprintf("%s/%s.log", outPath, addtime))
	if err != nil {
		panic(err)
	}
	logger := zerolog.New(create).Level(level).With().Timestamp().Logger()

	return &BlueLog{
		l: logger,
	}
}

func (l *BlueLog) Info(msg string) {
	l.l.Info().Msg(msg)
}

func (l *BlueLog) Warn(msg string) {
	l.l.Warn().Msg(msg)
}

func (l *BlueLog) Err(msg string) {
	l.l.Error().Msg(msg)
}
