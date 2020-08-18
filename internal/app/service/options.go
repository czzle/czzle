package service

import (
	"github.com/czzle/czzle/internal/pkg/cache"
	"github.com/czzle/czzle/internal/pkg/cache/lrucache"

	"github.com/go-kit/kit/log"
)

type options struct {
	secret string
	cache  cache.Cache
	logger log.Logger
}

var defaultOptions = options{
	secret: "secret",
	cache:  lrucache.New(100 * lrucache.MB),
}

type Option func(*options)

func WithSecret(secret string) Option {
	return func(opts *options) {
		opts.secret = secret
	}
}

func WithCache(cache cache.Cache) Option {
	return func(opts *options) {
		opts.cache = cache
	}
}

func WithLogger(logger log.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}
