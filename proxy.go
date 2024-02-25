package socks

import (
	"context"
	"io"
	"net"
	"time"
)

// ProxyHandler is the interface for handling the proxy requests
type ProxyHandler interface {
	Init(context.Context, Request) (context.Context, io.ReadWriteCloser, *Error)
	ReadFromClient(context.Context, io.ReadCloser, io.WriteCloser) error
	ReadFromRemote(context.Context, io.ReadCloser, io.WriteCloser) error
	Close(context.Context) error
	Refresh(context.Context)
}

// Proxy is the main struct
type Proxy struct {
	ServerAddr   string
	Done         chan struct{}
	Proxyhandler ProxyHandler
	Timeout      time.Duration
	Log          Logger
}

// Start is the main function to start a proxy
func (p *Proxy) Start(ctx context.Context) error {
	if p.Log == nil {
		p.Log = &NilLogger{} // allow not to set logger
	}

	listener, err := net.Listen("tcp", p.ServerAddr)
	if err != nil {
		return err
	}
	go p.run(ctx, listener)
	return nil
}

func (p *Proxy) run(ctx context.Context, listener net.Listener) {
	for {
		select {
		case <-p.Done:
			return
		default:
			connection, err := listener.Accept()
			if err == nil {
				go p.handle(ctx, connection)
			} else {
				p.Log.Errorf("Error accepting conn: %v", err)
			}
		}
	}
}

// Stop stops the proxy
func (p *Proxy) Stop() {
	p.Log.Warn("Stopping proxy")
	if p.Done == nil {
		return
	}
	close(p.Done)
	p.Done = nil
}
