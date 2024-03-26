package blue

import (
	"blue/bsp"
	"strconv"
)

func (c *Client) Version() (string, error) {
	build := bsp.NewRequestBuilder(bsp.VERSION).Build()

	return c.exec(build)
}

func (c *Client) Del(key string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.DEL).WithKey(key).Build()

	return c.exec(build)
}

func (c *Client) Nset(k, num string) (string, error) {
	_, err := strconv.Atoi(num)
	if err != nil {
		return "", ErrDataType(num)
	}

	build := bsp.NewRequestBuilder(bsp.NSET).
		WithKey(k).
		WithValueNum(num).
		Build()

	return c.exec(build)
}

func (c *Client) Get(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.GET).WithKey(k).Build()
	return c.exec(build)
}

func (c *Client) Set(k, v string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.SET).
		WithKey(k).
		WithValueStr(v).
		Build()

	return c.exec(build)
}

func (c *Client) Len(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.LEN).WithKey(k).Build()

	return c.exec(build)
}

func (c *Client) Kvs() (string, error) {
	build := bsp.NewRequestBuilder(bsp.KVS).Build()

	return c.exec(build)
}

func (c *Client) Nget(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.NGET).
		WithKey(k).
		Build()

	return c.exec(build)
}

func (c *Client) Select(num ...string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.SELECT)
	if len(num) != 0 {
		build.WithKey(num[0])
	}
	return c.exec(build.Build())
}

func (c *Client) Llen(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.LLEN).WithKey(k).Build()

	return c.exec(build)
}

func (c *Client) Lget(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.LGET).
		WithKey(k).
		Build()

	return c.exec(build)
}

func (c *Client) Lset(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.LSET).WithKey(k).Build()

	return c.exec(build)
}

func (c *Client) Rpop(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.RPOP).WithKey(k).Build()

	return c.exec(build)
}

func (c *Client) Rpush(k, v string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.RPUSH).
		WithKey(k).
		WithValueStr(v).
		Build()

	return c.exec(build)
}

func (c *Client) Lpop(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.LPOP).WithKey(k).Build()

	return c.exec(build)
}

func (c *Client) Lpush(k, v string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.LPUSH).
		WithKey(k).
		WithValueStr(v).
		Build()

	return c.exec(build)
}

func (c *Client) Expire(k, v string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.EXPIRE).
		WithKey(k).
		WithValueStr(v).
		Build()

	return c.exec(build)
}

func (c *Client) Incr(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.INCR).WithKey(k).Build()

	return c.exec(build)
}
