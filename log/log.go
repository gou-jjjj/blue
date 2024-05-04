package log

import (
	"blue/common/filename"
	pri "blue/common/print"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// 定义常量和全局变量
const (
	bufSize = 10000 // 缓冲区大小
)

var (
	blog *BlueLog // 日志实例

	stdio = Stdio{
		blog:   blog,
		lLevel: InfoLevel,
		null:   true,
	}

	// 信息、警告和错误日志记录函数
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

// BIter 接口定义了日志记录的方法
type BIter interface {
	info(msg string)
	err(msg string)
	warn(msg string)
}

// Stdio 结构体封装了日志记录的标准输入输出设置
type Stdio struct {
	blog   BIter
	lLevel _LogLevel
	null   bool
}

// info、warn和err方法分别用于记录信息、警告和错误日志
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

// Init 初始化日志系统
// output: 日志输出方式
// level: 日志级别
// outPath: 日志文件或目录路径
func Init(output string, level string, outPath string) {
	i := logLevel(level)
	blog = NewBlueLog(output, outPath)

	stdio = Stdio{
		blog:   blog,
		lLevel: i,
		null:   blog == nil,
	}

	pri.LogInitSuccess()
}

// BlueLog 结构体定义了文件日志的功能实现
type BlueLog struct {
	of      *os.File    // 日志文件句柄
	logCh   chan string // 日志消息通道
	bufSize int         // 缓冲区大小
}

// NewBlueLog 创建一个新的日志实例
// output: 日志输出目标（当前仅支持文件）
// outPath: 日志文件或目录路径
func NewBlueLog(output string, outPath string) *BlueLog {
	if output != "file" {
		return nil
	}

	b := &BlueLog{
		bufSize: bufSize,
		logCh:   make(chan string, bufSize),
	}

	// 根据路径打开或创建日志文件
	var err error
	if filepath.Ext(outPath) == ".log" {
		b.of, err = os.Open(outPath)
		if err != nil {
			panic(err)
		}

	} else {
		// 如果目录不存在，则创建目录
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

	go b.sync() // 启动日志同步线程

	return b
}

// sync 方法用于将日志消息写入文件
func (l *BlueLog) sync() {
	var msg string
	for msg = range l.logCh {
		_, err := l.of.WriteString(msg)
		if err != nil {
			return
		}
	}
}

// info、warn和err方法用于向日志通道发送消息
func (l *BlueLog) info(msg string) {
	l.logCh <- msg
}

func (l *BlueLog) warn(msg string) {
	l.logCh <- msg
}

func (l *BlueLog) err(msg string) {
	l.logCh <- msg
}

// _LogMsg 创建一个新的日志消息
func _LogMsg(msg string, level _LogLevel) *Msg {
	m := getMsg()
	m.data = msg
	m.dataType = level
	m.dataTime = time.Now()
	return m
}
