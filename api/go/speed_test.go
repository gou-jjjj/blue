package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"net"
	"testing"
)

func Redis() {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	c.Set(context.Background(), "a", 1, 0)
}

var c = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

var dial, _ = net.Dial("tcp", "localhost:8080")

// BenchmarkRedis-24           2731            421280 ns/op
// BenchmarkRedis-24           3144            386315 ns/op
// BenchmarkRedis-24           3334            375299 ns/op
// BenchmarkRedis-24             10         118124530 ns/op
// BenchmarkRedis-24             10         118174020 ns/op
func BenchmarkRedis(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.Set(context.Background(), "a", 1, 0)
	}
}

func Blue() {

	//  [1000001 10 1100001 1010 1 1010]
	_, _ = dial.Write([]byte{
		0x41, 0x02, 0x61, 0x0a, 0x01, 0x0a,
	})
}

// BenchmarkBlue-24            3145            426977 ns/op
func BenchmarkBlue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Blue()
	}
}
