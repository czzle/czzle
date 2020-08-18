package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/czzle/czzle/internal/pkg/appd"
	"github.com/czzle/czzle/internal/pkg/appd/deps"
	"github.com/czzle/czzle/internal/app/transport"
)

func loadTransport(srv *server) appd.BeforeStartFunc {
	return func(env *appd.Env) error {
		logger := deps.GetLogger(env)
		if logger == nil {
			return errors.New("logger not found")
		}
		logger.Log("status", "loading transport")
		addr := fmt.Sprintf(":%s", *srv.cfg.port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		srv.listener = listener
		srv.handler = transport.NewHTTPServer(srv.eps)
		return nil
	}
}

func startTransport(srv *server) appd.StartFunc {
	return func(env *appd.Env, ch chan<- error) {
		logger := deps.GetLogger(env)
		if logger == nil {
			ch <- errors.New("logger not found")
			return
		}
		logger.Log("status", "starting transport")
		ch <- http.Serve(srv.listener, srv.handler)
	}
}

func disconnectTransport(srv *server) appd.BeforeStopFunc {
	return func(env *appd.Env) error {
		if srv.listener == nil {
			return nil
		}
		logger := deps.GetLogger(env)
		if logger == nil {
			return errors.New("logger not found")
		}
		logger.Log("status", "disconnecting transport")
		return srv.listener.Close()
	}
}
