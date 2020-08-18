package vd

import (
	"context"
)

type ValidateFunc func(ctx context.Context, param string) error

func (fn ValidateFunc) Validate(ctx context.Context, param string) error {
	return fn(ctx, param)
}

type Validator interface {
	Validate(ctx context.Context, param string) error
}

type RootValidator interface {
	Check(validators ...Validator) error
}

type root struct {
	ctx context.Context
}

func (r root) Check(validators ...Validator) error {
	for _, v := range validators {
		err := v.Validate(r.ctx, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func Check(validators ...Validator) error {
	ctx := context.Background()
	for _, v := range validators {
		err := v.Validate(ctx, "")
		if err != nil {
			return err
		}
	}
	return nil
}

func Context(ctx context.Context) RootValidator {
	return root{ctx}
}
