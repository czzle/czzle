package appd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/czzle/czzle/pkg/multierr"
)

type BeforeStartFunc func(env *Env) error
type BeforeStopFunc func(env *Env) error
type StartFunc func(env *Env, ch chan<- error)

type Logger interface {
	Log(kv ...interface{}) error
}

type AppOpt func(*AppD)

func WithLogger(logger Logger) AppOpt {
	return func(a *AppD) {
		a.logger = logger
	}
}
func Start(funcs ...StartFunc) AppOpt {
	return func(a *AppD) {
		a.startFuncs = append(a.startFuncs, funcs...)
	}
}

func BeforeStop(funcs ...BeforeStopFunc) AppOpt {
	return func(a *AppD) {
		a.prestopFuncs = append(a.prestopFuncs, funcs...)
	}
}

func BeforeStart(funcs ...BeforeStartFunc) AppOpt {
	return func(a *AppD) {
		a.prestartFuncs = append(a.prestartFuncs, funcs...)
	}
}

func ProbePort(port string) AppOpt {
	return func(a *AppD) {
		a.probePort = port
	}
}

func ProbeEnabled(enabled bool) AppOpt {
	return func(a *AppD) {
		a.probeEnabled = enabled
	}
}

// AppD - application
type AppD struct {
	ready         bool
	probeEnabled  bool
	probePort     string
	stopch        chan bool
	startFuncs    []StartFunc
	prestartFuncs []BeforeStartFunc
	prestopFuncs  []BeforeStopFunc
	logger        Logger
	env           Env
	sync.Mutex
}

func New(opts ...AppOpt) *AppD {
	app := &AppD{
		probePort:    "8888",
		probeEnabled: true,
		env:          make(Env),
	}
	for _, opt := range opts {
		opt(app)
	}
	return app
}

// Run - starts application
func (a *AppD) Run() {
	a.log("status", "initiating")
	// create stop channel
	a.stopch = make(chan bool)
	// create error channel
	errch := make(chan error)
	probe := new(probe)
	if a.probeEnabled {
		probe.SetPort(a.probePort)
		go func(ch chan<- error) {
			probe.SetHealthy(true)
			a.log("status", "initiating probe")
			errch <- probe.Run()
		}(errch)
	}
	err := a.beforeStart()
	if err != nil {
		merr := multierr.From(err)
		merr = merr.With(a.beforeStop())
		if a.probeEnabled {
			merr = merr.With(probe.Close())
		}
		fmt.Println(err)
		os.Exit(1)
		return
	}
	a.log("status", "finished before start processes")
	// run app
	go func(ch chan<- error) {
		// run app
		a.log("status", "running")
		ch <- a.start()
	}(errch)
	var sigch = make(chan os.Signal)
	signal.Notify(sigch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	go func() {
		sig := <-sigch
		a.log("status", "caught sig", "sig", sig)
		a.Stop()
	}()
	if a.probeEnabled {
		probe.SetReady(true)
	}
	code := 0
	// wait for stop or error
	for {
		select {
		// stop
		case <-a.stopch:
			if a.probeEnabled {
				probe.SetReady(false)
			}
			merr := multierr.From(err)
			merr = merr.With(a.beforeStop())
			if a.probeEnabled {
				a.log("status", "closing probe")
				merr = merr.With(probe.Close())
			}
			a.stopch <- true
			a.stopch = nil
			a.log("status", "stopped")
			os.Exit(code)
			return
		// error
		case err = <-errch:
			if err != nil {
				a.log("err", err)
				code = 1
			}
			go a.Stop()
		}
	}
}

// Stop - stops application
func (a *AppD) Stop() {
	a.log("status", "stopping")
	if a.stopch == nil {
		return
	}
	a.stopch <- true
	<-a.stopch
}

// run - runs application
func (a *AppD) start() (err error) {
	a.log("status", "starting")
	ch := make(chan error)
	for _, f := range a.startFuncs {
		go f(&a.env, ch)
	}
	return <-ch
}

func (a *AppD) log(kv ...interface{}) {
	a.Lock()
	defer a.Unlock()
	if a.logger == nil {
		return
	}
	a.logger.Log(kv...)
}

func (a *AppD) beforeStop() error {
	a.log("status", "running before stop processes")
	merr := multierr.New("got errors while stopping")
	gotErr := false
	for _, f := range a.prestopFuncs {
		for f == nil {
			continue
		}
		if err := f(&a.env); err != nil {
			a.log("err", err)
			merr = merr.With(err)
			gotErr = true
		}
	}
	if gotErr {
		return merr
	}
	return nil
}

func (a *AppD) beforeStart() error {
	a.log("status", "running before start processes")
	for _, f := range a.prestartFuncs {
		for f == nil {
			continue
		}
		if err := f(&a.env); err != nil {
			a.log("err", err)
			return err
		}
	}
	return nil
}
