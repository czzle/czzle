package service

import (
	"context"
	"time"

	"github.com/czzle/czzle"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type logging struct {
	logger log.Logger
	next   czzle.Service
}

func makeLogging(logger log.Logger) czzle.Middleware {
	return func(next czzle.Service) czzle.Service {
		return logging{logger, next}
	}
}

func (mw logging) Begin(ctx context.Context, client *czzle.ClientInfo) (puzzle *czzle.Puzzle, err error) {
	logger := level.Debug(mw.logger)
	logger = log.With(logger, "method", "Begin")
	logger.Log(
		"status", "calling service",
		"client.id", client.GetID(),
		"client.time", client.GetTime(),
		"client.ip", client.GetIP(),
		"client.user_agent", client.GetUserAgent(),
	)
	start := time.Now()
	defer func() {
		kv := []interface{}{
			"status", "called service",
			"took", time.Since(start),
		}
		if puzzle != nil {
			kv = append(kv,
				"puzzle.level", puzzle.GetLevel(),
				"puzzle.token", puzzle.GetToken(),
				"puzzle.expires_at", puzzle.GetExpirationTime(),
				"puzzle.issued_at", puzzle.GetIssuedTime(),
			)
			if client := puzzle.GetClient(); client != nil {
				kv = append(kv,
					"puzzle.client_id", client.GetID(),
					"puzzle.client_time", client.GetTime(),
					"puzzle.client_ip", client.GetIP(),
					"puzzle.client_user_agent", client.GetUserAgent(),
				)
			}
			if tm := puzzle.GetTileMap(); tm != nil {
				kv = append(kv,
					"puzzle.tile_map.size", tm.GetSize(),
					"puzzle.tile_map.tiles_len", len(tm.GetTiles()),
				)
			}
		}
		if err != nil {
			kv = append(kv, "err", err)
			logger = level.Error(logger)
		}
		logger.Log(kv...)
	}()
	return mw.next.Begin(ctx, client)
}

func (mw logging) Solve(ctx context.Context, solution *czzle.Solution) (results *czzle.Results, err error) {
	logger := level.Debug(mw.logger)
	logger = log.With(logger, "method", "Solve")
	logger.Log(
		"status", "calling service",
		"solution.token", solution.GetToken(),
		"solution.actions_len", len(solution.GetActions()),
	)
	start := time.Now()
	defer func() {
		kv := []interface{}{
			"status", "called service",
			"took", time.Since(start),
		}
		if results != nil {
			kv = append(kv,
				"results.ok", results.IsOK(),
				"results.took", results.GetTook(),
			)
			if results.GetAccessToken() != "" {
				kv = append(kv,
					"results.access_token", results.GetAccessToken(),
				)
			}
			if puzzle := results.GetNext(); puzzle != nil {
				kv = append(kv,
					"results.next.level", puzzle.GetLevel(),
					"results.next.token", puzzle.GetToken(),
					"results.next.expires_at", puzzle.GetExpirationTime(),
					"results.next.issued_at", puzzle.GetIssuedTime(),
				)
				if client := puzzle.GetClient(); client != nil {
					kv = append(kv,
						"results.next.client.id", client.GetID(),
						"results.next.client.time", client.GetTime(),
						"results.next.client.ip", client.GetIP(),
						"results.next.client.user_agent", client.GetUserAgent(),
					)
				}
				if tm := puzzle.GetTileMap(); tm != nil {
					kv = append(kv,
						"results.next.tile_map.size", tm.GetSize(),
						"results.next.tile_map.tiles_len", len(tm.GetTiles()),
					)
				}
			}
		}
		if err != nil {
			kv = append(kv, "err", err)
			logger = level.Error(logger)
		}
		logger.Log(kv...)
	}()
	return mw.next.Solve(ctx, solution)
}

func (mw logging) Validate(ctx context.Context, accessToken string) (ok bool, err error) {
	logger := level.Debug(mw.logger)
	logger = log.With(logger, "method", "Validate")
	logger.Log(
		"status", "calling service",
		"access_token", accessToken,
	)
	start := time.Now()
	defer func() {
		kv := []interface{}{
			"status", "called service",
			"took", time.Since(start),
			"ok", ok,
		}
		if err != nil {
			kv = append(kv, "err", err)
			logger = level.Error(logger)
		}
		logger.Log(kv...)
	}()
	return mw.next.Validate(ctx, accessToken)
}
