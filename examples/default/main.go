package main

import (
	"time"

	socks "github.com/firefart/gosocks"
)

func main() {
	handler := socks.DefaultHandler{
		Timeout: 1 * time.Second,
	}
	listen := "127.0.0.1:1080"
	p := socks.Proxy{
		ServerAddr:   listen,
		Proxyhandler: handler,
		Timeout:      1 * time.Second,
		Log:          &socks.NilLogger{},
	}
	p.Log.Infof("starting SOCKS server on %s", listen)
	if err := p.Start(); err != nil {
		panic(err)
	}
	<-p.Done
}
