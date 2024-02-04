package main

import (
	"context"
	"flag"
	"time"

	socks "github.com/firefart/gosocks"
	"github.com/sirupsen/logrus"
)

func main() {
	debugMode := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	log := logrus.New()

	if *debugMode {
		log.SetLevel(logrus.DebugLevel)
		log.Debug("debug mode enabled")
	}

	handler := socks.DefaultHandler{
		Timeout: 1 * time.Second,
		Log:     log,
	}
	listen := "127.0.0.1:1080"
	p := socks.Proxy{
		ServerAddr:   listen,
		Proxyhandler: handler,
		Timeout:      1 * time.Second,
		Log:          log,
	}
	p.Log.Infof("starting SOCKS server on %s", listen)
	if err := p.Start(context.Background()); err != nil {
		panic(err)
	}
	<-p.Done
}
