package main

import (
	"net"
	"net/http"

	"github.com/czzle/czzle"
	"github.com/czzle/czzle/internal/app/endpoints"
	"github.com/czzle/czzle/internal/pkg/appd"
	"github.com/czzle/czzle/internal/pkg/appd/deps"
)

type config struct {
	port   *string
	secret *string
}

type server struct {
	cfg      config
	eps      endpoints.Set
	svc      czzle.Service
	listener net.Listener
	handler  http.Handler
}

func main() {
	instance := "czzle"
	srv := new(server)
	srv.cfg = config{
		port:   deps.ConfigParam("port", "8080", "http port"),
		secret: deps.ConfigParam("secret", "secret", "captcha generation secret"),
	}
	app := appd.New(
		deps.Config(instance),
		deps.Logger(instance),
		appd.BeforeStart(
			loadService(srv),
			loadEndpoints(srv),
			loadTransport(srv),
		),
		appd.Start(
			startTransport(srv),
		),
		appd.BeforeStop(
			disconnectTransport(srv),
		),
	)
	app.Run()
}
