package socks

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

// DefaultHandler is the default socks5 implementation
type DefaultHandler struct {
	// Timeout defines the connect timeout to the destination
	Timeout time.Duration
	log     Logger
}

// Init is the default socks5 implementation
func (s DefaultHandler) Init(request Request) (io.ReadWriteCloser, error) {
	target := fmt.Sprintf("%s:%d", request.DestinationAddress, request.DestinationPort)
	s.log.Infof("Connecting to target %s", target)
	remote, err := net.DialTimeout("tcp", target, s.Timeout)
	if err != nil {
		return nil, err
	}
	return remote, nil
}

// ReadFromClient is the default socks5 implementation
func (s DefaultHandler) ReadFromClient(ctx context.Context, client io.ReadCloser, remote io.WriteCloser) error {
	if _, err := io.Copy(remote, client); err != nil {
		return err
	}
	return nil
}

// ReadFromRemote is the default socks5 implementation
func (s DefaultHandler) ReadFromRemote(ctx context.Context, remote io.ReadCloser, client io.WriteCloser) error {
	if _, err := io.Copy(client, remote); err != nil {
		return err
	}
	return nil
}

// Close is the default socks5 implementation
func (s DefaultHandler) Close() error {
	return nil
}

// Refresh is the default socks5 implementation
func (s DefaultHandler) Refresh(ctx context.Context) {}
