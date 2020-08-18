package service

import (
	"context"
	"fmt"
	"time"

	"github.com/czzle/czzle"
	"github.com/czzle/czzle/pkg/vd"
)

type validating struct {
	next czzle.Service
}

func makeValidating() czzle.Middleware {
	return func(next czzle.Service) czzle.Service {
		return validating{next}
	}
}

func (mw validating) Begin(ctx context.Context, client *czzle.ClientInfo) (*czzle.Puzzle, error) {
	minTime := time.Now().Add(time.Hour*(-24)).UnixNano() / int64(time.Millisecond)
	err := vd.Check(
		vd.Param("client").AND(
			vd.NotNull(client),
			vd.Param("id").NotNull(client.GetID()),
			vd.Param("ip").OR(
				vd.IPv4(client.GetIP()),
				vd.IPv6(client.GetIP()),
			),
			vd.Param("time").Int64Min(client.GetTime(), minTime),
			vd.Param("user_agent").Length(client.GetUserAgent(), 3, 200),
		),
	)
	if err != nil {
		return nil, czzle.ErrInvalidArgument.With(err)
	}
	return mw.next.Begin(ctx, client)
}

func (mw validating) Solve(ctx context.Context, solution *czzle.Solution) (*czzle.Results, error) {
	err := vd.Check(
		vd.Param("solution").AND(
			vd.NotNull(solution),
			vd.Param("actions").AND(
				vd.Length(solution.GetActions(), 1, 1000),
				mw.validateActionList(solution.GetActions()),
			),
		),
	)
	if err != nil {
		return nil, czzle.ErrInvalidArgument.With(err)
	}
	return mw.next.Solve(ctx, solution)
}

func (mw validating) Validate(ctx context.Context, accessToken string) (bool, error) {
	return mw.next.Validate(ctx, accessToken)
}

func (mw validating) validateActionList(list czzle.ActionList) vd.ValidateFunc {
	return func(ctx context.Context, param string) error {
		now := time.Now().UnixNano()
		minTime := now - int64(time.Hour*24)
		for i, a := range list {
			aParam := fmt.Sprintf("%s.action[%d]", param, i)
			vd.Param(aParam).AND(
				vd.NotNull(a),
				vd.Param("type").True(
					czzle.AllowedActions.Has(a.GetType()),
				),
				vd.Param("time").Int64Min(a.GetTime(), minTime),
				vd.Param("data").AND(
					mw.validateActionData(a.GetType(), a.GetData()),
				),
			)(ctx, param)
		}
		return nil
	}
}

func (mw validating) validateActionData(t czzle.ActionType, data czzle.ActionData) vd.ValidateFunc {
	return func(ctx context.Context, param string) error {
		param = fmt.Sprintf("%s.data", param)
		switch d := data.(type) {
		case nil:
			if t != czzle.BeginAction &&
				t == czzle.ConfirmAction {
				return fmt.Errorf("%s: expected null data", param)
			}
			return nil
		case *czzle.FlipActionData:
			if t != czzle.FlipAction {
				return fmt.Errorf("%s: expected flip action data", param)
			}
			return vd.Param("").AND(
				vd.Param("pos").AND(
					vd.NotNull(d.GetPos()),
					vd.Param("x").Int(d.GetPos().GetX(), 0, 2),
					vd.Param("y").Int(d.GetPos().GetY(), 0, 2),
				),
			).Validate(ctx, param)
		default:
			return fmt.Errorf("%s: unknown data", param)
		}
	}
}
