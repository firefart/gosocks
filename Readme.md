# GOSOCKS

Basic golang implementation of a socks5 proxy. This implementation is currently not feature complete and only supports the `CONNECT` command and no authentication.

This implemention also defines some handlers you can use to implement your own protocol behind this proxy server. This can be useful if you come a across a protocol that can be abused for proxy functionality and build a socks5 proxy around it.

The SOCKS protocol is defined in [rfc1928](https://tools.ietf.org/html/rfc1928)

## Documentation

[https://pkg.go.dev/github.com/firefart/gosocks](https://pkg.go.dev/github.com/firefart/gosocks)

## Handler Interface

```golang
type ProxyHandler interface {
	Init(Request) (io.ReadWriteCloser, *Error)
	ReadFromClient(context.Context, io.ReadCloser, io.WriteCloser) error
	ReadFromRemote(context.Context, io.ReadCloser, io.WriteCloser) error
	Close() error
	Refresh(ctx context.Context)
}
```

### Init

Init is called before the copy operations and it should return a connection to the target that is ready to receive data.

### ReadFromClient

ReadFromClient is the method that handles the data copy from the client (you) to the remote connection. You can see the `DefaultHandler` for a sample implementation.

### ReadFromRemote

ReadFromRemote is the method that handles the data copy from the remote connection to the client (you). You can see the `DefaultHandler` for a sample implementation.

### Close

Close is called after the request finishes or errors out. It is used to clean up any connections in your custom implementation.

### Refresh

Refresh is called in a seperate goroutine and should loop forever to do refreshes of the connection if needed. The passed in context is cancelled after the request so be sure to check on the Done event.

## Usage

### Default Usage

```golang
package main

import (
	"time",

	socks "github.com/firefart/gosocks"
	"github.com/sirupsen/logrus"
)

func main() {
	handler := socks.DefaultHandler{
		Timeout: 1*time.Second,
	}
	listen := "127.0.0.1:1080"
	p := socks.Proxy{
		ServerAddr:   listen,
		Proxyhandler: handler,
		Timeout:      1*time.Second,
		Log:          logrus.New(),
	}
	p.Log.Infof("starting SOCKS server on %s", listen)
	if err := p.Start(); err != nil {
		panic(err)
	}
	<-p.Done
}
```

### Usage with custom handlers

```golang
package main

import (
	"time"
	"io"
	"fmt"
	"net"
	"context"

	socks "github.com/firefart/gosocks"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	handler := MyCustomHandler{
		Timeout: 1*time.Second,
		PropA:  "A",
		PropB:  "B",
		Log:    log,
	}
	p := socks.Proxy{
		ServerAddr:   "127.0.0.1:1080",
		Proxyhandler: handler,
		Timeout:      1*time.Second,
		Log:          log,
	}
	log.Infof("starting SOCKS server on %s", listen)
	if err := p.Start(); err != nil {
		panic(err)
	}
	<-p.Done
}

type MyCustomHandler struct {
	Timeout time.Duration,
	PropA   string,
	PropB   string,
	Log     Logger,
}

func (s *MyCustomHandler) Init(request socks.Request) (io.ReadWriteCloser, *socks.Error) {
	conn, err := net.DialTimeout("tcp", s.Server, s.Timeout)
	if err != nil {
		return nil, &socks.SocksError{Reason: socks.RequestReplyHostUnreachable, Err: fmt.Errorf("error on connecting to server: %w", err)}
	}
	return conn, nil
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

func (s *MyCustomHandler) Close() error {
	return nil
}
```
