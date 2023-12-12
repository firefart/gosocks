package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	socks "github.com/firefart/gosocks"
)

func main() {
	log := &socks.NilLogger{}
	handler := MyCustomHandler{
		Timeout: 1 * time.Second,
		PropA:   "A",
		PropB:   "B",
		Log:     log,
	}
	p := socks.Proxy{
		ServerAddr:   "127.0.0.1:1080",
		Proxyhandler: &handler,
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
	Timeout time.Duration
	PropA   string
	PropB   string
	Log     socks.Logger
}

func (s *MyCustomHandler) Init(request socks.Request) (io.ReadWriteCloser, *socks.Error) {
	target := fmt.Sprintf("%s:%d", request.DestinationAddress, request.DestinationPort)
	s.Log.Infof("Connecting to target %s", target)
	remote, err := net.DialTimeout("tcp", target, s.Timeout)
	if err != nil {
		return nil, socks.NewError(socks.RequestReplyNetworkUnreachable, err)
	}
	return remote, nil
}

func (s *MyCustomHandler) Refresh(ctx context.Context) {
	tick := time.NewTicker(10 * time.Second)
	select {
	case <-ctx.Done():
		return
	case <-tick.C:
		s.Log.Debug("refreshing connection")
	}
}

func (s *MyCustomHandler) ReadFromRemote(ctx context.Context, remote io.ReadCloser, client io.WriteCloser) error {
	i, err := io.Copy(client, remote)
	if err != nil {
		return err
	}
	s.Log.Debugf("wrote %d bytes to client", i)
	return nil
}

func (s *MyCustomHandler) ReadFromClient(ctx context.Context, client io.ReadCloser, remote io.WriteCloser) error {
	i, err := io.Copy(remote, client)
	if err != nil {
		return err
	}
	s.Log.Debugf("wrote %d bytes to remote", i)
	return nil
}

func (s MyCustomHandler) Close() error {
	return nil
}
