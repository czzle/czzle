package czzle

import (
	"github.com/czzle/czzle/pkg/multierr"
)

var (
	errg               = multierr.Group("czzle")
	ErrUnexpected      = errg.Code(0).Kind(multierr.Unexpected).New("unexpected")
	ErrInvalidArgument = errg.Code(1).Kind(multierr.InvalidArgument).New("invalid argument")
	ErrNotFound        = errg.Code(2).Kind(multierr.NotFound).New("not found")
)

func init() {
	multierr.Register(
		ErrUnexpected,
		ErrInvalidArgument,
		ErrNotFound,
	)
}
