package socks

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

// DefaultHandler is the default socks5 implementation
type DefaultHandler struct {
	// Timeout defines the connect timeout to the destination
	Timeout time.Duration
	Log     Logger
}

// Init is the default socks5 implementation
func (s DefaultHandler) Init(ctx context.Context, request Request) (io.ReadWriteCloser, *Error) {
	target := request.GetDestinationString()
	if s.Log != nil {
		s.Log.Infof("Connecting to target %s", target)
	}
	remote, err := net.DialTimeout("tcp", target, s.Timeout)
	if err != nil {
		return nil, &Error{Reason: RequestReplyHostUnreachable, Err: fmt.Errorf("error on connecting to server: %w", err)}
	}
	return remote, nil
}

const bufferSize = 10240

// ReadFromClient is the default socks5 implementation
func (s DefaultHandler) ReadFromClient(ctx context.Context, client io.ReadCloser, remote io.WriteCloser) error {
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

// ReadFromRemote is the default socks5 implementation
func (s DefaultHandler) ReadFromRemote(ctx context.Context, remote io.ReadCloser, client io.WriteCloser) error {
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

// Close is the default socks5 implementation
func (s DefaultHandler) Close(ctx context.Context) error {
	return nil
}

// Refresh is the default socks5 implementation
func (s DefaultHandler) Refresh(ctx context.Context) {}
