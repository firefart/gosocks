package socks

// https://tools.ietf.org/html/rfc1928

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
)

func (p *Proxy) handle(conn net.Conn) {
	defer conn.Close()
	defer func() {
		p.Log.Debug("client connection closed")
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p.Log.Debugf("got connection from %s", conn.RemoteAddr().String())

	if err := p.socks(ctx, conn); err != nil {
		// send error reply
		p.Log.Errorf("socks error: %v", err.Err)
		if err := p.socksErrorReply(ctx, conn, err.Reason); err != nil {
			p.Log.Error(err)
			return
		}
	}
}

func (p *Proxy) socks(ctx context.Context, conn net.Conn) *Error {
	defer func() {
		if err := p.Proxyhandler.Close(ctx); err != nil {
			p.Log.Errorf("error on close: %v", err)
		}
	}()

	if err := p.handleConnect(ctx, conn); err != nil {
		return err
	}

	request, err := p.handleRequest(ctx, conn)
	if err != nil {
		return err
	}

	p.Log.Infof("Connecting to %s", request.GetDestinationString())

	// Should we assume connection succeed here?
	remote, err := p.Proxyhandler.Init(ctx, *request)
	if err != nil {
		return err
	}
	defer remote.Close()
	p.Log.Infof("Connection established %s - %s", conn.RemoteAddr().String(), request.GetDestinationString())

	err = p.handleRequestReply(ctx, conn, request)
	if err != nil {
		return err
	}

	p.Log.Debug("beginning of data copy")

	wg := &sync.WaitGroup{}
	errChannel1 := make(chan error, 1)
	errChannel2 := make(chan error, 1)
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	wg.Add(2)

	go p.copyClientToRemote(ctx2, conn, remote, wg, errChannel1)
	go p.copyRemoteToClient(ctx2, remote, conn, wg, errChannel2)
	go p.Proxyhandler.Refresh(ctx2)

	p.Log.Debug("waiting for copy to finish")
	wg.Wait()
	// stop refreshing the connection
	cancel()
	if err := <-errChannel1; err != nil {
		return NewError(RequestReplyHostUnreachable, err)
	}
	if err := <-errChannel2; err != nil {
		return NewError(RequestReplyHostUnreachable, err)
	}
	p.Log.Debug("end of connection handling")

	return nil
}

func (p *Proxy) copyClientToRemote(ctx context.Context, client io.ReadCloser, remote io.WriteCloser, wg *sync.WaitGroup, errChannel chan<- error) {
	defer wg.Done()
	defer close(errChannel)

	select {
	case <-p.Done:
		errChannel <- nil
		return
	default:
		if err := p.Proxyhandler.ReadFromClient(ctx, client, remote); err != nil {
			errChannel <- fmt.Errorf("error on copy from Client to Remote: %v", err)
			return
		}
		errChannel <- nil
		return
	}
}

func (p *Proxy) copyRemoteToClient(ctx context.Context, remote io.ReadCloser, client io.WriteCloser, wg *sync.WaitGroup, errChannel chan<- error) {
	defer wg.Done()
	defer close(errChannel)

	select {
	case <-p.Done:
		errChannel <- nil
		return
	default:
		if err := p.Proxyhandler.ReadFromRemote(ctx, remote, client); err != nil {
			errChannel <- fmt.Errorf("error on copy from Remote to Client: %v", err)
			return
		}
		errChannel <- nil
		return
	}
}

func (p *Proxy) socksErrorReply(ctx context.Context, conn io.ReadWriteCloser, reason RequestReplyReason) error {
	// send error reply
	repl, err := requestReply(nil, reason)
	if err != nil {
		return err
	}
	err = connectionWrite(ctx, conn, repl, p.Timeout)
	if err != nil {
		return err
	}

	return nil
}

func (p *Proxy) handleConnect(ctx context.Context, conn io.ReadWriteCloser) *Error {
	buf, err := connectionRead(ctx, conn, p.Timeout)
	if err != nil {
		return NewError(RequestReplyConnectionRefused, err)
	}
	header, err := parseHeader(buf)
	if err != nil {
		return NewError(RequestReplyConnectionRefused, err)
	}
	switch header.Version {
	case Version4:
		return NewError(RequestReplyCommandNotSupported, fmt.Errorf("socks4 not yet implemented"))
	case Version5:
	default:
		return NewError(RequestReplyCommandNotSupported, fmt.Errorf("version %#x not yet implemented", byte(header.Version)))
	}

	methodSupported := false
	for _, x := range header.Methods {
		if x == MethodNoAuthRequired {
			methodSupported = true
			break
		}
	}
	if !methodSupported {
		return NewError(RequestReplyMethodNotSupported, fmt.Errorf("we currently only support no authentication"))
	}
	reply := make([]byte, 2)
	reply[0] = byte(Version5)
	reply[1] = byte(MethodNoAuthRequired)
	err = connectionWrite(ctx, conn, reply, p.Timeout)
	if err != nil {
		return NewError(RequestReplyGeneralFailure, fmt.Errorf("could not send connect reply: %w", err))
	}
	return nil
}

func (p *Proxy) handleRequest(ctx context.Context, conn io.ReadWriteCloser) (*Request, *Error) {
	buf, err := connectionRead(ctx, conn, p.Timeout)
	if err != nil {
		return nil, NewError(RequestReplyGeneralFailure, fmt.Errorf("error on ConnectionRead: %w", err))
	}
	request, err2 := parseRequest(buf)
	if err2 != nil {
		return nil, err2
	}
	return request, nil
}

func (p *Proxy) handleRequestReply(ctx context.Context, conn io.ReadWriteCloser, request *Request) *Error {
	repl, err := requestReply(request, RequestReplySucceeded)
	if err != nil {
		return NewError(RequestReplyGeneralFailure, fmt.Errorf("error on requestReply: %w", err))
	}
	err = connectionWrite(ctx, conn, repl, p.Timeout)
	if err != nil {
		return NewError(RequestReplyGeneralFailure, fmt.Errorf("error on RequestResponse: %w", err))
	}

	return nil
}
