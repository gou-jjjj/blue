package blue

import (
	"blue/bsp"
	"blue/common/strbytes"
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
	if !strbytes.CheckInt(num) {
		return "", ErrArgu("num")
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

func (c *Client) Ping() (string, error) {
	build := bsp.NewRequestBuilder(bsp.PING).Build()
	return c.exec(build)
}

func (c *Client) Sadd(k, v string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.SADD).
		WithKey(k).
		WithValueStr(v).
		Build()

	return c.exec(build)
}

func (c *Client) Spop(k string) (string, error) {

	build := bsp.NewRequestBuilder(bsp.SPOP).
		WithKey(k).
		Build()

	return c.exec(build)
}

func (c *Client) Sin(k, v string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.SIN).
		WithKey(k).
		WithValueStr(v).
		Build()

	return c.exec(build)
}

func (c *Client) Sdel(k, v string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.SDEL).
		WithKey(k).
		WithValueStr(v).
		Build()

	return c.exec(build)
}

func (c *Client) Sget(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.SGET).
		WithKey(k).
		Build()

	return c.exec(build)
}

func (c *Client) Help(k string) (string, error) {
	build := bsp.NewRequestBuilder(bsp.HELP).WithKey(k).Build()

	return c.exec(build)
}

func (c *Client) Exit() (s string, err error) {
	build := bsp.NewRequestBuilder(bsp.EXIT).Build()

	c.exit(build)
	return
}

func (c *Client) Dbsize() (s string, err error) {
	build := bsp.NewRequestBuilder(bsp.DBSIZE).Build()
	return c.exec(build)
}

func (c *Client) Type(k string) (s string, err error) {
	build := bsp.NewRequestBuilder(bsp.TYPE).WithKey(k).Build()
	return c.exec(build)
}

func (c *Client) Auth(k ...string) (s string, err error) {
	build := []byte{}
	if len(k) == 0 {
		build = bsp.NewRequestBuilder(bsp.AUTH).Build()
	} else {
		build = bsp.NewRequestBuilder(bsp.AUTH).WithKey(k[0]).Build()
	}
	return c.exec(build)
}
