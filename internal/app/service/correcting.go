package service

import (
	"context"
	"sort"
	"strings"

	"github.com/czzle/czzle"
	"github.com/czzle/czzle/pkg/uuid"
)

type correcting struct {
	next czzle.Service
}

func makeCorrecting() czzle.Middleware {
	return func(next czzle.Service) czzle.Service {
		return correcting{next}
	}
}

func (mw correcting) Begin(ctx context.Context, client *czzle.ClientInfo) (puzzle *czzle.Puzzle, err error) {
	if client.GetID().IsNull() {
		client.SetID(uuid.New())
	}
	// fix ip
	ip := client.GetIP()
	ip = strings.TrimSpace(ip)
	ip = strings.ToLower(ip)
	client.SetIP(ip)
	// fix user agent
	agent := client.GetUserAgent()
	agent = strings.TrimSpace(agent)
	client.SetUserAgent(agent)
	return mw.next.Begin(ctx, client)
}

func (mw correcting) Solve(ctx context.Context, solution *czzle.Solution) (results *czzle.Results, err error) {
	// fix token
	token := solution.GetToken()
	token = strings.TrimSpace(token)
	solution.SetToken(token)
	// sort actions
	actions := solution.GetActions()
	sort.Sort(actions)
	solution.SetActions(actions)
	return mw.next.Solve(ctx, solution)
}

func (mw correcting) Validate(ctx context.Context, accessToken string) (ok bool, err error) {
	// fix access token
	accessToken = strings.TrimSpace(accessToken)
	accessToken = strings.ToLower(accessToken)
	return mw.next.Validate(ctx, accessToken)
}
