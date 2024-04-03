package benchmark

import (
	blue "blue/api/go"
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"testing"
)

var (
	bluetest1 = func(b *testing.B) {
		c, err := blue.NewClient(blue.WithDefaultOpt(), func(c *blue.Config) {
			c.Addr = "39.101.169.250:7894"
		})
		if err != nil {
			b.Fatal(err)
		}
		defer c.Close()

		for i := 0; i < 1e2; i++ {
			_, err = c.Set("1", "1")
			if err != nil {
				b.Fatal(err)
			}
		}
		for i := 0; i < 1e2; i++ {
			res, err := c.Get("1")
			if err != nil || res != "1" {
				b.Fatal(err)
			}
		}

		for i := 0; i < 1e2; i++ {
			_, err = c.Del("1")
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	redistest1 = func(b *testing.B) {
		rdb := redis.NewClient(&redis.Options{
			Addr: "39.101.169.250:7893",
		})
		defer rdb.Close()
		ctx := context.Background()
		for i := 0; i < 1e2; i++ {
			err := rdb.Set(ctx, "1", "1", 0).Err()
			if err != nil {
				b.Fatal(err)
			}
		}

		for i := 0; i < 1e2; i++ {
			err := rdb.Get(ctx, "1").Err()
			if err != nil {
				b.Fatal(err)
			}
		}
		for i := 0; i < 1e2; i++ {
			err := rdb.Del(ctx, "1").Err()
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	bluetestn = func(b *testing.B) {
		c, err := blue.NewClient(blue.WithDefaultOpt(), func(c *blue.Config) {
			c.Addr = "39.101.169.250:7894"
		})
		if err != nil {
			b.Fatal(err)
		}
		defer c.Close()
		for i := 0; i < 1e2; i++ {
			_, err = c.Set(strconv.Itoa(i), strconv.Itoa(i))
			if err != nil {
				b.Fatal(err)
			}
		}

		for i := 0; i < 1e2; i++ {
			res, err := c.Get(strconv.Itoa(i))
			if err != nil || res != strconv.Itoa(i) {
				b.Fatal(err)
			}
		}

		for i := 0; i < 1e2; i++ {
			_, err = c.Del(strconv.Itoa(i))
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	redistestn = func(b *testing.B) {
		rdb := redis.NewClient(&redis.Options{
			Addr: "39.101.169.250:7893",
		})
		defer rdb.Close()
		ctx := context.Background()
		for i := 0; i < 1e2; i++ {
			err := rdb.Set(ctx, strconv.Itoa(i), strconv.Itoa(i), 0).Err()
			if err != nil {
				b.Fatal(err)
			}
		}
		for i := 0; i < 1e2; i++ {
			err := rdb.Get(ctx, strconv.Itoa(i)).Err()
			if err != nil {
				b.Fatal(err)
			}
		}

		for i := 0; i < 1e2; i++ {
			err := rdb.Del(ctx, strconv.Itoa(i)).Err()
			if err != nil {
				b.Fatal(err)
			}
		}

	}

	bluetest1e5 = func(b *testing.B) {
		c, err := blue.NewClient(blue.WithDefaultOpt(), func(c *blue.Config) {
			c.Addr = "39.101.169.250:7894"
		})
		if err != nil {
			b.Fatal(err)
		}
		defer c.Close()

		for i := 0; i < 1e1; i++ {
			vv := strconv.Itoa(i)
			_, err = c.Set(vv, vv)
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	redistest1e5 = func(b *testing.B) {
		rdb := redis.NewClient(&redis.Options{
			Addr: "39.101.169.250:7893",
		})
		defer rdb.Close()
		ctx := context.Background()
		for i := 0; i < 1e1; i++ {
			vv := strconv.Itoa(i)
			err := rdb.Set(ctx, vv, vv, 0).Err()
			if err != nil {
				b.Fatal(err)
			}
		}

	}
)

func BenchmarkSingleThreadSingleKey(b *testing.B) {
	b.Run("BenchmarkRedis", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			redistest1(b)
		}
	})

	b.Run("BenchmarkBlue", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bluetest1(b)
		}
	})
}

func BenchmarkMultiThreadSingleKey(b *testing.B) {
	bluetest := func() {
		c, err := blue.NewClient(blue.WithDefaultOpt(), func(c *blue.Config) {
			c.Addr = "39.101.169.250:7894"
		})
		if err != nil {
			b.Fatal(err)
		}
		defer c.Close()

		for i := 0; i < 1e2; i++ {
			_, err = c.Set("1", "1")
			if err != nil {
				b.Fatal(err)
			}
		}

		for i := 0; i < 1e2; i++ {
			_, err = c.Del("1")
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	redistest := func() {
		rdb := redis.NewClient(&redis.Options{
			Addr: "39.101.169.250:7893",
		})
		ctx := context.Background()
		for i := 0; i < 1e2; i++ {
			err := rdb.Set(ctx, "1", "1", 0).Err()
			if err != nil {
				b.Fatal(err)
			}
		}

		for i := 0; i < 1e2; i++ {
			err := rdb.Del(ctx, "1").Err()
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	b.Run("BenchmarkRedis", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wg := new(sync.WaitGroup)
			for j := 0; j < 100; j++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					redistest()
				}()
			}
			wg.Wait()
		}
	})

	b.Run("BenchmarkBlue", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wg := new(sync.WaitGroup)
			for j := 0; j < 100; j++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					bluetest()
				}()
			}
			wg.Wait()
		}
	})

}

func BenchmarkSingleThreadMultiKey(b *testing.B) {
	b.Run("BenchmarkRedis", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			redistestn(b)
		}
	})

	b.Run("BenchmarkBlue", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bluetestn(b)
		}
	})
}

func BenchmarkMultiThreadMultiKey(b *testing.B) {
	b.Run("BenchmarkRedis", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wg := new(sync.WaitGroup)
			for j := 0; j < 100; j++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					redistestn(b)
				}()
			}
			wg.Wait()

		}
	})

	b.Run("BenchmarkBlue", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wg := new(sync.WaitGroup)
			for j := 0; j < 100; j++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					bluetestn(b)
				}()
			}
			wg.Wait()
		}
	})
}

func BenchmarkSet1e5(b *testing.B) {
	b.Run("BenchmarkRedis", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			redistest1e5(b)
		}
	})

	b.Run("BenchmarkBlue", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			bluetest1e5(b)
		}
	})
}
