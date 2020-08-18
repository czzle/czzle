package deps

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/czzle/czzle/internal/pkg/appd"
	"github.com/czzle/czzle/internal/pkg/version"

	sd "github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
)

type ConsulRegistrarConfig struct {
	ServiceName string
	ServicePort *string
}

func ConsulClient() appd.AppOpt {
	addr := ConfigParam("consul-addr", "consul:8500", "consul address")
	return func(d *appd.AppD) {
		appd.BeforeStart(
			connectConsulClient(addr),
		)(d)
	}
}

func GetConsulClient(env *appd.Env) *consulapi.Client {
	res := env.Get("deps-consul-client")
	if res == nil {
		return nil
	}
	return res.(*consulapi.Client)
}

func ConsulRegistrar(serviceName string, servicePort *string) appd.AppOpt {
	return func(d *appd.AppD) {
		appd.BeforeStart(
			initiateRegistrar(serviceName, servicePort),
			registerConsul(),
		)(d)
		appd.BeforeStop(
			deregisterConsul(),
		)(d)
	}
}

func initiateRegistrar(serviceName string, servicePort *string) appd.BeforeStartFunc {
	return func(env *appd.Env) error {
		logger := GetLogger(env)
		if logger == nil {
			return errors.New("logger not found")
		}
		logger.Log("status", "initiating registrar")
		client := GetConsulClient(env)
		if client == nil {
			return errors.New("consul client not found")
		}
		id := make([]byte, 8)
		rand.Read(id)
		port, err := strconv.Atoi(*servicePort)
		if err != nil {
			return err
		}
		addr, err := os.Hostname()
		if err != nil {
			return err
		}
		r := sd.NewRegistrar(sd.NewClient(client), &consulapi.AgentServiceRegistration{
			ID:      fmt.Sprintf("%s-%s", serviceName, hex.EncodeToString(id)),
			Name:    serviceName,
			Address: addr,
			Port:    port,
			Tags: []string{
				fmt.Sprintf("czzle-%s", version.Version),
			},
			Check: &consulapi.AgentServiceCheck{
				HTTP:     fmt.Sprintf("http://%s:%d/readiness", addr, 8888),
				Interval: "10s",
				Timeout:  "1s",
				Notes:    "Ready health check",
			},
		}, logger)
		env.Set("deps-consul-registrar", r)
		return nil
	}
}

func connectConsulClient(addr *string) appd.BeforeStartFunc {
	return func(env *appd.Env) error {
		logger := GetLogger(env)
		if logger != nil {
			logger.Log(
				"status", "connecting to consul",
				"addr", *addr,
			)
		}
		cc := consulapi.DefaultConfig()
		cc.Address = *addr
		client, err := consulapi.NewClient(cc)
		if err != nil {
			return err
		}
		env.Set("deps-consul-client", client)
		return nil
	}
}

func registerConsul() appd.BeforeStartFunc {
	return func(env *appd.Env) error {
		logger := GetLogger(env)
		if logger == nil {
			return errors.New("logger not found")
		}
		logger.Log("status", "registering service to consul")
		r := env.Get("deps-consul-registrar")
		if r == nil {
			return errors.New("consul registrar not found")
		}
		reg, ok := r.(*sd.Registrar)
		if !ok {
			return errors.New("consul registrar type assertion")
		}
		reg.Register()
		return nil
	}
}

func deregisterConsul() appd.BeforeStopFunc {
	return func(env *appd.Env) error {
		logger := GetLogger(env)
		if logger == nil {
			return errors.New("logger not found")
		}
		logger.Log("status", "deregistering service from consul")
		r := env.Get("deps-consul-registrar")
		if r == nil {
			return nil
		}
		reg, ok := r.(*sd.Registrar)
		if !ok {
			return errors.New("consul registrar type assertion")
		}
		reg.Deregister()
		return nil
	}
}
