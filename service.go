package czzle

import (
	"context"
)

type Middleware func(Service) Service

type Service interface {
	Begin(ctx context.Context, client *ClientInfo) (*Puzzle, error)
	Solve(ctx context.Context, solution *Solution) (*Results, error)
	Validate(ctx context.Context, accessToken string) (bool, error)
}
