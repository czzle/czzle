package deps

import (
	"flag"
	"os"
	"strings"

	"github.com/czzle/czzle/internal/pkg/appd"

	"github.com/peterbourgon/ff"
)

type cfgEntry struct {
	key          string
	value        *string
	defaultValue string
	usage        string
}

var cfgm map[string]*cfgEntry

func ConfigParam(key, defaultValue, usage string) *string {
	entry, ok := cfgm[key]
	if !ok {
		var value string
		entry = &cfgEntry{
			key:          key,
			value:        &value,
			defaultValue: defaultValue,
			usage:        usage,
		}
		cfgm[key] = entry
	}
	return entry.value
}

func Config(prefix string) appd.AppOpt {
	return func(d *appd.AppD) {
		appd.BeforeStart(
			loadConfig(prefix),
		)(d)
	}
}

func loadConfig(prefix string) appd.BeforeStartFunc {
	fs := flag.NewFlagSet(prefix, flag.ExitOnError)
	return func(env *appd.Env) error {
		for _, e := range cfgm {
			fs.StringVar(e.value, e.key, e.defaultValue, e.usage)
		}
		return ff.Parse(fs, os.Args[1:],
			ff.WithEnvVarPrefix(strings.ToUpper(prefix)),
		)
	}
}

func init() {
	cfgm = make(map[string]*cfgEntry)
}
