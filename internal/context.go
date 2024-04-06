package internal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"blue/bsp"
	"blue/common/rand"
	"blue/config"
	"blue/log"
)

const (
	TokenLen = 10
)

// 定义上下文结构体，封装了请求处理所需的基本信息和状态
type Context struct {
	context.Context
	nextExec  Exec          // 下一个执行器
	conn      net.Conn      // 网络连接
	db        uint8         // 数据库索引
	cliToken  string        // 客户端令牌
	request   *bsp.BspProto // 请求协议对象
	response  bsp.Reply     // 响应协议对象
	maxActive time.Duration // 最大活动时间
	isclose   int           // 连接状态标志
}

// 定义连接池，用于复用Context对象
var bconnPool = sync.Pool{
	New: func() interface{} {
		return &Context{
			db:       1,
			isclose:  0,
			cliToken: rand.RandString(TokenLen),
		}
	},
}

// String方法返回上下文的字符串表示
func (c *Context) String() string {
	return fmt.Sprintf("RemoteAddr:%s,cliToken:%s", c.conn.RemoteAddr().String(), c.cliToken)
}

// NewContext创建一个新的上下文对象，如果连接池中有可用对象，则复用，否则新建
func NewContext(ctx context.Context, conn net.Conn) *Context {
	bconn, ok := bconnPool.Get().(*Context)
	if !ok {
		nctx := &Context{
			Context:  ctx,
			conn:     conn,
			db:       1,
			isclose:  0,
			cliToken: rand.RandString(TokenLen),
			maxActive: time.Duration(config.CliCfg.ClientActive) *
				time.Second,
		}
		log.Info(bconn.String())
		return nctx
	}
	bconn.Context = ctx
	bconn.conn = conn
	bconn.db = 1
	bconn.isclose = 0
	bconn.cliToken = rand.RandString(TokenLen)
	bconn.maxActive = time.Duration(config.CliCfg.ClientActive) * time.Second
	log.Info(bconn.String())
	return bconn
}

// GetDB返回当前使用的数据库索引
func (c *Context) GetDB() uint8 {
	return c.db
}

// SetDB设置要使用的数据库索引
func (c *Context) SetDB(index uint8) {
	c.db = index
}

// Reply构造响应并返回给客户端，如果连接已关闭则返回错误
func (c *Context) Reply() (int, error) {
	if c.response == nil {
		return c.conn.Write(bsp.NewErr(bsp.ErrReplication).Bytes())
	}
	log.Info(fmt.Sprintf("reply: %s", c.response.String()))
	if c.isClose() {
		return 0, errors.New("conn is close")
	}
	return c.conn.Write(c.response.Bytes())
}

// isClose检查连接是否已关闭
func (c *Context) isClose() bool {
	return c.isclose == 1
}

// Close关闭连接，清理上下文资源，并将上下文对象返回到连接池
func (c *Context) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.db = 1
	c.Context = nil
	c.nextExec = nil
	c.request = nil
	c.response = nil
	c.maxActive = 0
	c.isclose = 1
	bconnPool.Put(c)
	return
}
