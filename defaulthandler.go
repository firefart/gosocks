package socks

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// DefaultHandler is the default socks5 implementation
type DefaultHandler struct {
	// Timeout defines the connect timeout to the destination
	Timeout time.Duration
}

// PreHandler is the default socks5 implementation
func (s DefaultHandler) PreHandler(request Request) (io.ReadWriteCloser, error) {
	target := fmt.Sprintf("%s:%d", request.DestinationAddress, request.DestinationPort)
	log.Infof("Connecting to target %s", target)
	remote, err := net.DialTimeout("tcp", target, s.Timeout)
	if err != nil {
		return nil, err
	}
	return remote, nil
}

// CopyFromClientToRemote is the default socks5 implementation
func (s DefaultHandler) CopyFromClientToRemote(client, remote io.ReadWriteCloser) error {
	if _, err := io.Copy(client, remote); err != nil {
		return err
	}
	return nil
}

// CopyFromRemoteToClient is the default socks5 implementation
func (s DefaultHandler) CopyFromRemoteToClient(remote, client io.ReadWriteCloser) error {
	if _, err := io.Copy(remote, client); err != nil {
		return err
	}
	return nil
}

// Cleanup is the default socks5 implementation
func (s DefaultHandler) Cleanup() error {
	return nil
}

// Refresh is the default socks5 implementation
func (s DefaultHandler) Refresh(ctx context.Context) {}
