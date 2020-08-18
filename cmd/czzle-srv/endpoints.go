package main

import (
	"errors"

	"github.com/czzle/czzle/internal/app/endpoints"

	"github.com/czzle/czzle/internal/pkg/appd"
	"github.com/czzle/czzle/internal/pkg/appd/deps"
)

func loadEndpoints(srv *server) appd.BeforeStartFunc {
	return func(env *appd.Env) error {
		logger := deps.GetLogger(env)
		if logger == nil {
			return errors.New("logger not found")
		}
		logger.Log("status", "creating endpoints")
		srv.eps = endpoints.New(
			srv.svc,
			endpoints.WithLogger(logger),
		)
		return nil
	}
}
