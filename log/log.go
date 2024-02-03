package log

import (
	"blue/common/rand"
	"blue/config"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
)

var (
	Blog *BlueLog
)

func InitLog() {
	Blog = newLog()
}

func newLog() *BlueLog {
	filePath := config.BC.LogConfig.LogOut

	// 获取文件所在的目录路径
	dir := filepath.Dir(filePath)

	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err = os.MkdirAll(dir, 0777) // 使用0755权限创建目录
		if err != nil {
			panic(fmt.Sprintf("Failed to create directory: %v", err))
		}
	}

	open, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	log := NewZeroLog(zerolog.InfoLevel, 5, open)
	log.Info("log init success ...")
	return log
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
		create, err := os.Create(outPath + "/" + r + ".log")
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
