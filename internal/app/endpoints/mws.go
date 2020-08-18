package endpoints

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func makeLogging(name string, logger log.Logger) endpoint.Middleware {
	logger = log.WithPrefix(logger, "method", name)
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (res interface{}, err error) {
			level.Debug(logger).Log("msg", "calling endpoint")
			defer func(start time.Time) {
				logger := log.With(logger,
					"msg", "called endpoint",
					"took", time.Since(start),
				)
				if err != nil {
					level.Error(logger).Log(
						"err", err,
					)
				} else {
					level.Debug(logger).Log()
				}
				logger.Log()
			}(time.Now())
			return next(ctx, request)
		}
	}
}

func wrap(name string, ep endpoint.Endpoint, opts options) endpoint.Endpoint {
	if opts.logger != nil {
		ep = makeLogging(name, opts.logger)(ep)
	}
	return ep
}
