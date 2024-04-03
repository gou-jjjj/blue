package internal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"blue/log"
)

const (
	NetWork        = "tcp"
	MaxClientLimit = 1<<16 - 1
	TimeOut        = 10 * time.Second
)

type ConfigFunc func(*Config)

type Config struct {
	Ip          string
	Timeout     time.Duration
	Port        int
	ClientLimit int
	HandlerFunc ServerInter
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Ip, c.Port)
}

var defaultConfig = Config{
	Ip:          "127.0.0.1",
	Port:        8080,
	ClientLimit: 10,
	Timeout:     10 * time.Second,
}

type Server struct {
	c Config

	listen net.Listener

	cliOnlineTime sync.Map
	waitClient    sync.WaitGroup
	currentClient int32
	isClo         bool
	errClo        chan error
}

func checkConfig(c *Config) error {
	ip := net.ParseIP(c.Ip)
	if ip == nil {
		return errors.New("ip is invalid")
	}
	c.Ip = ip.String()

	if c.ClientLimit == 0 {
		c.ClientLimit = MaxClientLimit
	}

	if c.Timeout == 0 {
		c.Timeout = TimeOut
	}

	if c.HandlerFunc == nil {
		return errors.New("handler func is nil")
	}

	return nil
}

func NewServer(fs ...ConfigFunc) *Server {
	c := defaultConfig
	for _, f := range fs {
		f(&c)
	}
	err := checkConfig(&c)
	if err != nil {
		panic(err)
	}
	return &Server{
		c:             c,
		isClo:         false,
		errClo:        make(chan error, 1),
		currentClient: 0,
		cliOnlineTime: sync.Map{},
	}
}

func (s *Server) close() {
	_ = s.listen.Close()
	s.isClo = true
	close(s.errClo)
	log.Info("server start closing ...")
}

func (s *Server) Start() {
	var err error
	s.listen, err = net.Listen(NetWork, s.c.Addr())
	if err != nil {
		panic(err)
	}

	sigclo := make(chan os.Signal)
	signal.Notify(sigclo, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		select {
		case <-sigclo:
			log.Info("receive close signal...")
		case err = <-s.errClo:
			log.Error(fmt.Sprintf("server error: %v", err))
		}

		close(sigclo)
		s.close()
	}()

	log.Info(fmt.Sprintf("server listen on %v ...", s.c.Addr()))
	s.server()
}

func (s *Server) limit(conn net.Conn) bool {
	if s.currentClient == int32(s.c.ClientLimit) {
		_, _ = conn.Write([]byte("too many client ,try again later ..."))
		_ = conn.Close()
		return true
	}
	return false
}

func (s *Server) server() {
	for !s.isClo {
		conn, err := s.listen.Accept()
		if err != nil {
			s.errClo <- err
			return
		}

		if s.limit(conn) {
			continue
		}

		s.currentClient++
		s.waitClient.Add(1)
		s.cliOnlineTime.Store(conn.RemoteAddr().String(), time.Now())

		log.Info(fmt.Sprintf("new conn: %s ,%d ", conn.RemoteAddr().String(), s.currentClient))

		go func(conn net.Conn) {
			defer func() {
				atomic.AddInt32(&s.currentClient, -1)

				subtime, _ := s.cliOnlineTime.Load(conn.RemoteAddr().String())
				s.cliOnlineTime.Delete(conn.RemoteAddr().String())
				log.Info(fmt.Sprintf("close conn: %s, %v;", conn.RemoteAddr().String(),
					time.Now().Sub(subtime.(time.Time))))
				s.waitClient.Done()
				_ = conn.Close()
			}()

			s.c.HandlerFunc.Handle(context.Background(), conn)
		}(conn)
	}

	s.waitClient.Wait()
	log.Error("server closed ...")
}
