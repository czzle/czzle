package appd

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/czzle/czzle/pkg/multierr"
)

type probe struct {
	ready   bool
	healthy bool
	port    string
	server  *http.Server
	stopch  chan bool
	sync.Mutex
}

func (p *probe) SetReady(ready bool) *probe {
	p.Lock()
	defer p.Unlock()
	p.ready = ready
	return p
}

func (p *probe) SetHealthy(healthy bool) *probe {
	p.Lock()
	defer p.Unlock()
	p.healthy = healthy
	return p
}

func (p *probe) SetPort(port string) *probe {
	p.Lock()
	defer p.Unlock()
	p.port = port
	return p
}

func (p *probe) IsReady() bool {
	p.Lock()
	defer p.Unlock()
	return p.ready
}

func (p *probe) IsHealthy() bool {
	p.Lock()
	defer p.Unlock()
	return p.healthy
}

func (p *probe) GetPort() string {
	p.Lock()
	defer p.Unlock()
	return p.port
}

func (p *probe) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/liveliness", p.handleLiveliness())
	mux.HandleFunc("/readiness", p.handleReadiness())
	addr := fmt.Sprintf(":%s", p.GetPort())
	server := &http.Server{Addr: addr, Handler: mux}
	p.server = server
	return server.ListenAndServe()
}

func (p *probe) Close() error {
	err := p.server.Close()
	if err != nil {
		return multierr.New("got err while closing probe").With(err)
	}
	return nil
}

func (p *probe) handleLiveliness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !p.IsHealthy() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (p *probe) handleReadiness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !p.IsReady() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
