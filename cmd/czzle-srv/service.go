package main

import (
	"errors"

	"github.com/czzle/czzle/internal/app/service"

	"github.com/czzle/czzle/internal/pkg/appd"
	"github.com/czzle/czzle/internal/pkg/appd/deps"
)

func loadService(srv *server) appd.BeforeStartFunc {
	return func(env *appd.Env) error {
		logger := deps.GetLogger(env)
		if logger == nil {
			return errors.New("logger not found")
		}
		logger.Log("status", "creating service")
		srv.svc = service.New(
			service.WithSecret(*srv.cfg.secret),
			service.WithLogger(logger),
		)
		return nil
	}
}
