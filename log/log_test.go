package log

import (
	"fmt"
	"testing"
	"time"
)

func TestFileLog(t *testing.T) {
	InitLog("file", "info", "./testlog/info")

	for i := 0; i < 100; i++ {
		Info("1")
		Warn("2")
		Error("3")
	}

	InitLog("file", "warn", "./testlog/warn")

	for i := 0; i < 100; i++ {
		Info("1")
		Warn("2")
		Error("3")
	}

	InitLog("file", "err", "./testlog/err")

	for i := 0; i < 100; i++ {
		Info("1")
		Warn("2")
		Error("3")
	}

	time.Sleep(4 * time.Second)
}

func TestStdioLog(t *testing.T) {
	InitLog("notfile", "info", "./testlog/info")

	for i := 0; i < 2; i++ {
		Info("1")
		Warn("2")
		Error("3")
	}

	fmt.Println()
	InitLog("notfile", "warn", "./testlog/warn")

	for i := 0; i < 2; i++ {
		Info("1")
		Warn("2")
		Error("3")
	}

	fmt.Println()

	InitLog("notfile", "error", "./testlog/err")
	for i := 0; i < 2; i++ {
		Info("1")
		Warn("2")
		Error("3")
	}

	fmt.Println()
}
