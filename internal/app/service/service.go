package service

import (
	"context"
	"errors"

	"github.com/czzle/czzle"
	"github.com/czzle/czzle/internal/pkg/cache"
	"github.com/czzle/czzle/pkg/puzzler"
)

type service struct {
	puzzler *puzzler.Puzzler
	cache   cache.Cache
}

func New(with ...Option) czzle.Service {
	opts := defaultOptions
	for _, opt := range with {
		opt(&opts)
	}
	var svc czzle.Service
	{
		svc = service{
			cache:   opts.cache,
			puzzler: puzzler.New(opts.secret),
		}
		svc = makeValidating()(svc)
		if opts.logger != nil {
			svc = makeLogging(opts.logger)(svc)
		}
		svc = makeCorrecting()(svc)
	}
	return svc
}

func (svc service) Begin(ctx context.Context, client *czzle.ClientInfo) (*czzle.Puzzle, error) {
	return svc.puzzler.Make(client)
}

func (svc service) Solve(ctx context.Context, solution *czzle.Solution) (*czzle.Results, error) {
	return nil, errors.New("not implemented")
}

func (svc service) Validate(ctx context.Context, accessToken string) (bool, error) {
	return false, errors.New("not implemented")
}
