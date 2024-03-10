package log

import (
	"blue/common/rand"
	"blue/config"
	"context"
	"fmt"
	lo "log"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
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

func InitSyncLog() {
	Blog = newSyncLog()
}

func newSyncLog() *BlueLog {
	dir := config.BC.LogConfig.LogOut

	logx := NewZeroLog(logLevel(), 1, dir)
	lo.Printf("log init success ...")
	return logx
}

type BlueLog struct {
	count    int
	l        []zerolog.Logger
	sw       sync.WaitGroup
	errChan  chan string
	infoChan chan string
	warnChan chan string
	done     context.Context
	can      context.CancelFunc
}

func NewZeroLog(level zerolog.Level, count int, outPath string) *BlueLog {
	if count <= 0 {
		count = 1
	}

	zerolog.TimestampFieldName = "T"
	zerolog.MessageFieldName = "M"
	zerolog.LevelFieldName = "L"

	ctx, cancelFunc := context.WithCancel(context.Background())

	b := &BlueLog{
		errChan:  make(chan string, 10000), // 缓冲大小可以根据需要调整
		warnChan: make(chan string, 10000), // 缓冲大小可以根据需要调整
		infoChan: make(chan string, 10000), // 缓冲大小可以根据需要调整
		done:     ctx,
		can:      cancelFunc,
		count:    count,
		l:        make([]zerolog.Logger, count),
	}

	err := os.MkdirAll(outPath, 0777)
	if err != nil {
		panic(err)
	}

	for i := 0; i < count; i++ {
		r := rand.RandString(8)
		addtime := fmt.Sprintf("%s-%s", time.Now().Format("2006:01:02-15:04:05"), r)
		create, err := os.Create(outPath + "/" + addtime + ".log")
		if err != nil {
			panic(err)
		}
		b.l[i] = zerolog.New(create).Level(level).With().Timestamp().Logger()
	}

	b.startLogging()

	return b
}

func (l *BlueLog) startLogging() {
	for i := range l.l {
		go func(i int) {
			l.sw.Add(1)
			defer l.sw.Done()
			for {
				select {
				case msg, ok := <-l.errChan:
					if ok {
						l.l[i].Error().Msg(msg)
					}
				case msg, ok := <-l.infoChan:
					if ok {
						l.l[i].Info().Msg(msg)
					}
				case msg, ok := <-l.warnChan:
					if ok {
						l.l[i].Warn().Msg(msg)
					}
				case <-l.done.Done():
					return
				}
			}
		}(i)
	}
}

func (l *BlueLog) StopLogging() {
	close(l.errChan)
	close(l.infoChan)
	close(l.warnChan)

	for len(l.warnChan) != 0 || len(l.errChan) != 0 || len(l.infoChan) != 0 {
	}

	l.can()
	l.sw.Wait()
}

func (l *BlueLog) Info(msg string) {
	l.infoChan <- msg
}

func (l *BlueLog) Warn(msg string) {
	l.warnChan <- msg
}

func (l *BlueLog) Err(msg string) {
	l.errChan <- msg
}
