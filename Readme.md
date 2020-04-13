# GOSOCKS

Basic golang implementation of a socks5 proxy. This implementation is currently not feature complete and only supports the `CONNECT` command and no authentication.

This implemention also defines some handlers you can use to implement your own protocol behind this proxy server. This can be useful if you come a across a protocol that can be abused for proxy functionality and build a socks5 proxy around it.

## Handler Interface

```golang
type ProxyHandler interface {
	PreHandler(Request) (io.ReadWriteCloser, *Error)
	CopyFromClientToRemote(io.ReadCloser, io.WriteCloser) error
	CopyFromRemoteToClient(io.ReadCloser, io.WriteCloser) error
	Cleanup() error
	Refresh(ctx context.Context)
}
```

### PreHandler

PreHandler is called before the copy operations and it should return a connection to the target that is ready to receive data.

### CopyFromClientToRemote

CopyFromClientToRemote is the method that handles the data copy from the client (you) to the remote connection. You can see the `DefaultHandler` for a sample implementation.

### CopyFromRemoteToClient

CopyFromRemoteToClient is the method that handles the data copy from the remote connection to the client (you). You can see the `DefaultHandler` for a sample implementation.

### Cleanup

Cleanup is called after the request finishes or errors out. It is used to clean up any connections in your custom implementation.

### Refresh

Refresh is called in a seperate goroutine and should loop forever to do refreshes of the connection if needed. The passed in context is cancelled after the request so be sure to check on the Done event.

## Usage

### Default Usage

```golang
package main

import (
  "time",

  "github.com/firefart/gosocks"
)

func main() {
  handler := socks.DefaultHandler{
    Timeout: 1*time.Second,
  }
	p := socks.Proxy{
		ServerAddr:   "127.0.0.1:1080",
		Proxyhandler: handler,
	}
	log.Infof("starting SOCKS server on %s", listen)
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

  "github.com/firefart/gosocks"
)

func main() {
  handler := MyCustomHandler{
    Timeout: 1*time.Second,
    PropA: "A",
    PropB: "B",
  }
	p := socks.Proxy{
		ServerAddr:   "127.0.0.1:1080",
		Proxyhandler: handler,
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
}

func (s *MyCustomHandler) PreHandler(request socks.Request) (io.ReadWriteCloser, *socks.Error) {
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
		log.Debug("refreshing connection")
	}
}

func (s *MyCustomHandler) CopyFromRemoteToClient(remote io.ReadCloser, client io.WriteCloser) error {
	i, err := io.Copy(client, remote)
	if err != nil {
		return err
	}
	log.Debugf("wrote %d bytes to client", i)
	return nil
}

func (s *MyCustomHandler) CopyFromClientToRemote(client io.ReadCloser, remote io.WriteCloser) error {
	i, err := io.Copy(remote, client)
	if err != nil {
		return err
	}
	log.Debugf("wrote %d bytes to remote", i)
	return nil
}

func (s *MyCustomHandler) Cleanup() error {
	return nil
}
```
