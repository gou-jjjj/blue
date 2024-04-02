package log

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	log := NewLog(InfoLevel, "./testlog/info")

	for i := 0; i < 100; i++ {
		log.info("1")
		log.warn("2")
		log.err("3")
	}

	log1 := NewLog(WarnLevel, "./testlog/warn")

	for i := 0; i < 100; i++ {
		log1.info("1")
		log1.warn("2")
		log1.err("3")
	}
	log2 := NewLog(ErrorLevel, "./testlog/err")

	for i := 0; i < 100; i++ {
		log2.info("1")
		log2.warn("2")
		log2.err("3")
	}

	time.Sleep(4 * time.Second)
}
