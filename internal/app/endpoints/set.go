package endpoints

import (
	"github.com/czzle/czzle"

	"github.com/go-kit/kit/endpoint"
)

type Set struct {
	BeginEndpoint    endpoint.Endpoint
	SolveEndpoint    endpoint.Endpoint
	ValidateEndpoint endpoint.Endpoint
}

func New(svc czzle.Service, with ...Option) Set {
	var opts options
	for _, opt := range with {
		opt(&opts)
	}
	return Set{
		BeginEndpoint:    wrap("Begin", makeBeginEndpoint(svc), opts),
		SolveEndpoint:    wrap("Solve", makeSolveEndpoint(svc), opts),
		ValidateEndpoint: wrap("Validate", makeBeginEndpoint(svc), opts),
	}
}
