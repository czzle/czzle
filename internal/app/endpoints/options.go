package endpoints

import (
	"github.com/go-kit/kit/log"
)

type options struct {
	logger log.Logger
}

type Option func(*options)

func WithLogger(logger log.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}
