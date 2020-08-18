package endpoints

import (
	"context"

	"github.com/czzle/czzle"

	"github.com/go-kit/kit/endpoint"
)

func makeBeginEndpoint(svc czzle.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*czzle.BeginReq)
		if !ok {
			return nil, ErrTypeAssertion
		}
		puzzle, err := svc.Begin(ctx, req.GetClient())
		if err != nil {
			return nil, err
		}
		return &czzle.BeginRes{
			Puzzle: puzzle,
		}, nil
	}
}

func makeSolveEndpoint(svc czzle.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*czzle.SolveReq)
		if !ok {
			return nil, ErrTypeAssertion
		}
		results, err := svc.Solve(ctx, req.GetSolution())
		if err != nil {
			return nil, err
		}
		return &czzle.SolveRes{
			Results: results,
		}, nil
	}
}

func makeValidateEndpoint(svc czzle.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*czzle.ValidateReq)
		if !ok {
			return nil, ErrTypeAssertion
		}
		ok, err := svc.Validate(ctx, req.GetAccessToken())
		if err != nil {
			return nil, err
		}
		return &czzle.ValidateRes{
			OK: ok,
		}, nil
	}
}
