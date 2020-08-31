package socks

import (
	"context"
	"fmt"
	"io"
)

// connectionRead reads all data from a connection
func connectionRead(ctx context.Context, conn io.ReadCloser) ([]byte, error) {
	var ret []byte
	bufLen := 1024

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout on reading on connection")
		default:
			buf := make([]byte, bufLen)
			i, err := conn.Read(buf)
			if err != nil {
				return nil, fmt.Errorf("could not read from connection: %w", err)
			}
			ret = append(ret, buf[:i]...)
			if i < bufLen {
				return ret, nil
			}
		}
	}
}

// connectionWrite makes sure to write all data to a connection
func connectionWrite(ctx context.Context, conn io.WriteCloser, data []byte) error {
	toWriteLeft := len(data)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout on writing on connection")
		default:
			written, err := conn.Write(data)
			if err != nil {
				return err
			}
			if written == toWriteLeft {
				return nil
			}
			toWriteLeft -= written
		}
	}
}
