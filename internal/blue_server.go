package internal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	// 导入蓝色框架相关的包
	"blue/bsp"
	"blue/cluster"
	add "blue/common/network"
	print2 "blue/common/print"
	"blue/common/timewheel"
	"blue/config"
	"blue/log"
)

// 版本号定义
const version_ = "blue v0.1"

// Exec 接口定义执行链的接口
type Exec interface {
	ExecChain(*Context)
}

// ServerInter 接口定义服务器接口，继承自 Exec 接口
type ServerInter interface {
	Exec
	Handle(context.Context, net.Conn)
	Close()
}

// BlueServer 结构体实现 tcp.Handler 接口，作为 redis 服务器
type BlueServer struct {
	activeConn sync.Map         // 活跃连接的映射
	db         []*DB            // 数据库数组
	closed     atomic.Int32     // 服务器关闭状态标志
	cc         *cluster.Cluster // 集群控制对象
}

// NewBlueServer 创建一个新的 BlueServer 实例
func NewBlueServer(dbs ...*DB) *BlueServer {
	b := &BlueServer{
		db:         make([]*DB, len(dbs)),
		activeConn: sync.Map{},
	}

	// 初始化数据库数组
	for i := 0; i < len(dbs); i++ {
		b.db[i] = dbs[i]
	}

	// 根据配置决定是否开启集群模式
	if config.OpenCluster() {
		b.initClu()
	}

	return b
}

// initClu 初始化集群连接
func (svr *BlueServer) initClu() {
	svr.cc = cluster.NewCluster(
		config.CluCfg.TryTimes,
		config.CluCfg.Ip,
		config.CluCfg.Port,
		config.CluCfg.MyClusterAddr,
		config.SvrCfg.SvrAddr(),
		time.Duration(config.CluCfg.DialTimeout)*time.Second)

	// 解析并初始化集群地址
	if add.ParseAddr(config.CluCfg.ClusterAddr) {
		addrs := svr.cc.GetClusterAddr(config.CluCfg.ClusterAddr)
		svr.cc.InitClusterAddr(addrs...)
		svr.cc.Online(add.CombineAddr(svr.cc.LocalAddr(), config.SvrCfg.SvrAddr()))
	}

	print2.ClusterInitSuccess()
}

// Handle 处理网络连接，接收并执行 redis 命令
func (svr *BlueServer) Handle(ctx context.Context, conn net.Conn) {
	// 检查服务器是否已关闭，若已关闭则关闭连接并返回
	if svr.isClose() {
		_ = conn.Close()
		return
	}

	// 创建新的客户端上下文
	client := NewContext(ctx, conn)
	svr.activeConn.Store(client, struct{}{})
	defer func() {
		svr.closeClient(client)
	}()

	// 初始化 BSP 协议处理
	canCtx, cancelFunc := context.WithCancel(*client)
	bch, errch := bsp.BspProtos(canCtx, conn)
	defer func() {
		cancelFunc()
		close(bch)
		close(errch)
	}()

	// 根据集群配置选择处理逻辑
	if svr.isCluster() {
		svr.clusterHandle(client, bch, errch)
	} else {
		svr.localHandle(client, bch, errch)
	}
}

// localHandle 处理非集群模式下的客户端请求
func (svr *BlueServer) localHandle(ctx *Context, bch chan *bsp.BspProto, errch chan *bsp.ErrResp) {
	for !ctx.isClose() {
		// 设置超时定时器
		timewheel.Delay(ctx.maxActive, ctx.cliToken, func() {
			svr.closeClient(ctx)
		})

		select {
		case <-ctx.Done():
			return
		case req := <-bch:
			log.Info(fmt.Sprintf("local: %s", req.String()))
			ctx.request = req
			ctx.response = bsp.Reply(nil)

			svr.ExecChain(ctx)
			_, _ = ctx.Reply()
			bsp.PutBspProto(req)

		case err := <-errch:
			if !errors.Is(err, bsp.RequestEnd) {
				ctx.response = err
				_, _ = ctx.Reply()
			}
		}
	}
}

// clusterHandle 处理集群模式下的客户端请求
func (svr *BlueServer) clusterHandle(ctx *Context, bch chan *bsp.BspProto, errch chan *bsp.ErrResp) {
	for !ctx.isClose() {
		select {
		case <-ctx.Done():
			return
		case req := <-bch:
			log.Info(fmt.Sprintf("cluster: %s", req.String()))
			ctx.request = req
			ctx.response = bsp.Reply(nil)

			res, ok := svr.cc.Dial(ctx.request)
			if !ok {
				svr.ExecChain(ctx)
			} else {
				ctx.response = bsp.NewClusterReply(res)
			}

			_, _ = ctx.Reply()
			bsp.PutBspProto(req)

		case err := <-errch:
			if !errors.Is(err, bsp.RequestEnd) {
				ctx.response = err
				_, _ = ctx.Reply()
			}
		}
	}
}

// isCluster 判断是否开启了集群模式
func (svr *BlueServer) isCluster() bool {
	return svr.cc != nil
}

// isClose 判断服务器是否已关闭
func (svr *BlueServer) isClose() bool {
	return svr.closed.Load() == 1
}

// closeClient 关闭客户端连接，并从活跃连接映射中移除
func (svr *BlueServer) closeClient(client *Context) {
	if client == nil || client.isClose() {
		return
	}
	client.Close()
	svr.activeConn.Delete(client)
}

// Close 关闭服务器，停止处理请求
func (svr *BlueServer) Close() {
	svr.closed.Add(1)

	// 关闭所有活跃的客户端连接
	svr.activeConn.Range(func(key interface{}, _ interface{}) bool {
		client := key.(*Context)
		client.Close()
		return true
	})
}

// auth 对客户端进行认证
func (svr *BlueServer) auth(ctx *Context) {
	if ctx.request.Key() == "" {
		ctx.response = bsp.NewStr(ctx.cliToken)
		return
	}

	ctx.cliToken = ctx.request.Key()

	ctx.response = bsp.NewInfo(bsp.OK)
	return
}
