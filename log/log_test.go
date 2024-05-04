package log

import (
	"fmt"
	"testing"
	"time"
)

func TestFileLog(t *testing.T) {
	Init("file", "info", "./testlog/info")

	for i := 0; i < 100; i++ {
		Info("1")
		Warn("2")
		Error("3")
	}

	Init("file", "warn", "./testlog/warn")

	for i := 0; i < 100; i++ {
		Info("1")
		Warn("2")
		Error("3")
	}

	Init("file", "err", "./testlog/err")

	for i := 0; i < 100; i++ {
		Info("1")
		Warn("2")
		Error("3")
	}

	time.Sleep(4 * time.Second)
}

func TestStdioLog(t *testing.T) {
	Init("notfile", "info", "./testlog/info")

	for i := 0; i < 2; i++ {
		Info("1")
		Warn("2")
		Error("3")
	}

	fmt.Println()
	Init("notfile", "warn", "./testlog/warn")

	for i := 0; i < 2; i++ {
		Info("1")
		Warn("2")
		Error("3")
	}

	fmt.Println()

	Init("notfile", "error", "./testlog/err")
	for i := 0; i < 2; i++ {
		Info("1")
		Warn("2")
		Error("3")
	}

	fmt.Println()
}
