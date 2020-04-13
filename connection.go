package socks

import (
	"fmt"
	"io"
)

func connectionRead(conn io.ReadWriteCloser) ([]byte, error) {
	var ret []byte
	bufLen := 1024
	for {
		buf := make([]byte, bufLen)
		i, err := conn.Read(buf)
		if err != nil {
			return nil, fmt.Errorf("could not read from connection: %w", err)
		}
		ret = append(ret, buf[:i]...)
		if i < bufLen {
			break
		}
	}
	return ret, nil
}

func connectionWrite(conn io.ReadWriteCloser, data []byte) error {
	toWriteLeft := len(data)
	for {
		written, err := conn.Write(data)
		if err != nil {
			return err
		}
		if written == toWriteLeft {
			break
		}
		toWriteLeft -= written
	}
	return nil
}
