package czzle

import (
	"context"
	"errors"
)

var ErrClientServiceMissing = errors.New("client service is missing. use czzle.NewClient(...) to create a new client")

type Client struct {
	svc Service
}

func (c Client) Begin(client *ClientInfo) (*Puzzle, error) {
	return c.BeginCtx(context.Background(), client)
}

func (c Client) Solve(solution *Solution) (*Results, error) {
	return c.SolveCtx(context.Background(), solution)
}

func (c Client) Validate(ctx context.Context, accessToken string) (bool, error) {
	return c.ValidateCtx(context.Background(), accessToken)
}

func (c Client) BeginCtx(ctx context.Context, client *ClientInfo) (*Puzzle, error) {
	if c.svc == nil {
		return nil, ErrClientServiceMissing
	}
	return c.svc.Begin(ctx, client)
}

func (c Client) SolveCtx(ctx context.Context, solution *Solution) (*Results, error) {
	return c.svc.Solve(ctx, solution)
}

func (c Client) ValidateCtx(ctx context.Context, accessToken string) (bool, error) {
	return c.svc.Validate(ctx, accessToken)
}
