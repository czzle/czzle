package deps

import (
	"os"
	"github.com/czzle/czzle/internal/pkg/appd"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func Logger(instance string) appd.AppOpt {

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = level.NewFilter(logger, level.AllowAll())
	logger = level.NewInjector(logger, level.InfoValue())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "instance", instance)
	logger = log.With(logger, "caller", log.DefaultCaller)

	return func(d *appd.AppD) {
		appd.WithLogger(logger)(d)
		appd.BeforeStart(
			appendLogger(logger),
		)(d)
	}
}

func GetLogger(env *appd.Env) log.Logger {
	res := env.Get("deps-logger")
	if res == nil {
		return nil
	}
	return res.(log.Logger)
}

func appendLogger(logger log.Logger) appd.BeforeStartFunc {
	return func(env *appd.Env) error {
		env.Set("deps-logger", logger)
		return nil
	}
}
