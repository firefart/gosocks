package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	socks "github.com/firefart/gosocks"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
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
	if err := p.Start(); err != nil {
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
		return nil, &socks.Error{Reason: socks.RequestReplyHostUnreachable, Err: fmt.Errorf("error on connecting to server: %w", err)}
	}
	return conn, nil
}

func (s MyCustomHandler) Refresh(ctx context.Context) {
	tick := time.NewTicker(10 * time.Second)
	select {
	case <-ctx.Done():
		return
	case <-tick.C:
		s.Log.Debug("refreshing connection")
	}
}

const bufferSize = 10240

func (s MyCustomHandler) ReadFromRemote(ctx context.Context, remote io.ReadCloser, client io.WriteCloser) error {
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
