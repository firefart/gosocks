package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"time"

	socks "github.com/firefart/gosocks"
	"github.com/sirupsen/logrus"
)

func main() {
	debugMode := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	log := logrus.New()

	if *debugMode {
		log.SetLevel(logrus.DebugLevel)
		log.Debug("debug mode enabled")
	}

	handler := MyCustomHandler{
		Server:  "test",
		Timeout: 1 * time.Second,
		PropA:   "A",
		PropB:   "B",
		Log:     log,
	}
	p := socks.Proxy{
		ServerAddr:   "127.0.0.1:1080",
		Proxyhandler: handler,
		Timeout:      1 * time.Second,
		Log:          log,
	}
	log.Infof("starting SOCKS server on %s", p.ServerAddr)
	if err := p.Start(context.Background()); err != nil {
		panic(err)
	}
	<-p.Done
}

type MyCustomHandler struct {
	Server  string
	Timeout time.Duration
	PropA   string
	PropB   string
	Log     socks.Logger
}

func (s MyCustomHandler) Init(ctx context.Context, request socks.Request) (io.ReadWriteCloser, *socks.Error) {
	conn, err := net.DialTimeout("tcp", s.Server, s.Timeout)
	if err != nil {
		return nil, socks.NewError(socks.RequestReplyHostUnreachable, fmt.Errorf("error on connecting to server: %w", err))
	}
	return conn, nil
}

func (s MyCustomHandler) Refresh(ctx context.Context) {
	tick := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			s.Log.Debug("refreshing connection")
		}
	}
}

const bufferSize = 10240

type readDeadline interface {
	SetReadDeadline(time.Time) error
}
type writeDeadline interface {
	SetWriteDeadline(time.Time) error
}

func (s MyCustomHandler) ReadFromRemote(ctx context.Context, remote io.ReadCloser, client io.WriteCloser) error {
	timeOut := time.Now().Add(s.Timeout)

	ctx, cancel := context.WithDeadline(ctx, timeOut)
	defer cancel()

	if c, ok := client.(writeDeadline); ok {
		if err := c.SetWriteDeadline(timeOut); err != nil {
			return fmt.Errorf("could not set write deadline on client: %v", err)
		}
	}

	if c, ok := remote.(readDeadline); ok {
		if err := c.SetReadDeadline(timeOut); err != nil {
			return fmt.Errorf("could not set read deadline on remote: %v", err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			i, err := io.CopyN(client, remote, bufferSize)
			if errors.Is(err, io.EOF) {
				return nil
			} else if err != nil {
				return fmt.Errorf("ReadFromRemote: %w", err)
			}
			s.Log.Debugf("[socks] wrote %d bytes to client", i)
		}
	}
}

func (s MyCustomHandler) ReadFromClient(ctx context.Context, client io.ReadCloser, remote io.WriteCloser) error {
	timeOut := time.Now().Add(s.Timeout)

	ctx, cancel := context.WithDeadline(ctx, timeOut)
	defer cancel()

	if c, ok := remote.(writeDeadline); ok {
		if err := c.SetWriteDeadline(timeOut); err != nil {
			return fmt.Errorf("could not set write deadline on remote: %v", err)
		}
	}

	if c, ok := client.(readDeadline); ok {
		if err := c.SetReadDeadline(timeOut); err != nil {
			return fmt.Errorf("could not set read deadline on client: %v", err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			i, err := io.CopyN(remote, client, bufferSize)
			if errors.Is(err, io.EOF) {
				return nil
			} else if err != nil {
				return fmt.Errorf("ReadFromClient: %w", err)
			}
			s.Log.Debugf("[socks] wrote %d bytes to remote", i)
		}
	}
}

func (s MyCustomHandler) Close(ctx context.Context) error {
	return nil
}
