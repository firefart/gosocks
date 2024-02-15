# GOSOCKS

Basic golang implementation of a socks5 proxy. This implementation is currently not feature complete and only supports the `CONNECT` command and no authentication.

This implemention also defines some handlers you can use to implement your own protocol behind this proxy server. This can be useful if you come a across a protocol that can be abused for proxy functionality and build a socks5 proxy around it.

The SOCKS protocol is defined in [rfc1928](https://tools.ietf.org/html/rfc1928)

## Documentation

[https://pkg.go.dev/github.com/firefart/gosocks](https://pkg.go.dev/github.com/firefart/gosocks)

## Handler Interface

```golang
type ProxyHandler interface {
	Init(context.Context, Request) (io.ReadWriteCloser, *Error)
	ReadFromClient(context.Context, io.ReadCloser, io.WriteCloser) error
	ReadFromRemote(context.Context, io.ReadCloser, io.WriteCloser) error
	Close(context.Context) error
	Refresh(context.Context)
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

## Examples

Please see the `examples` directory.
