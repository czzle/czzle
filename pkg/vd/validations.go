package vd

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/czzle/czzle/pkg/multierr"
)

func makeErr(param string, format string, args ...interface{}) error {
	if param != "" {
		format = "%s: " + format
		args = append([]interface{}{param}, args...)
	}
	return fmt.Errorf(format, args...)
}

func wrapParam(name string, v Validator) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if param != "" && name != "" {
			param = fmt.Sprintf("%s.%s", param, name)
		} else if name != "" {
			param = name
		}
		return v.Validate(ctx, param)
	}
}

func Param(param string) Validations {
	return validations{
		param: param,
	}
}

type Validations interface {
	Param(param string) Validations
	Email(email string) ValidateFunc
	Not(vf ValidateFunc) ValidateFunc
	NotNullAll(list ...interface{}) ValidateFunc
	NotNull(src interface{}) ValidateFunc
	NullAll(list ...interface{}) ValidateFunc
	Null(src interface{}) ValidateFunc
	EqualOneOf(src interface{}, ofList ...interface{}) ValidateFunc
	EqualAll(first interface{}, with ...interface{}) ValidateFunc
	Equal(a interface{}, b interface{}) ValidateFunc
	True(b bool) ValidateFunc
	False(b bool) ValidateFunc
	LengthMin(src interface{}, min int) ValidateFunc
	LengthMax(src interface{}, max int) ValidateFunc
	Length(src interface{}, min, max int) ValidateFunc
	OR(validators ...Validator) ValidateFunc
	AND(validators ...Validator) ValidateFunc
	Custom(validator Validator) ValidateFunc
	Async() AsyncValidations
	Timeout(timeout time.Duration) Validations
	RegExp(str, exp string) ValidateFunc
	IntMin(i, min int) ValidateFunc
	IntMax(i, max int) ValidateFunc
	Int(i, min, max int) ValidateFunc
	Int64Min(i, min int64) ValidateFunc
	IPv4(str string) ValidateFunc
	IPv6(str string) ValidateFunc
	IP(str string) ValidateFunc
	Int64Max(i, max int64) ValidateFunc
	Int64(i, min, max int64) ValidateFunc
	InstanceOf(a, b interface{}) ValidateFunc
	Hex(str string) ValidateFunc
	Base32(str string) ValidateFunc
	Unique(arr ...interface{}) ValidateFunc
	Subdomain(str string) ValidateFunc
	Name(str string) ValidateFunc
	Code(str string) ValidateFunc
}

type validations struct {
	param   string
	timeout time.Duration
}

func (v validations) Param(param string) Validations {
	if v.param != "" {
		param = fmt.Sprintf("%s.%s", v.param, param)
	}
	return validations{
		param: param,
	}
}

func (v validations) Custom(validator Validator) ValidateFunc {
	return wrap(v.param, v.timeout, validator)
}

func (v validations) RegExp(exp, str string) ValidateFunc {
	return wrap(v.param, v.timeout, RegExp(exp, str))
}

func (v validations) InstanceOf(a, b interface{}) ValidateFunc {
	return wrap(v.param, v.timeout, InstanceOf(a, b))
}

func (v validations) LengthMin(src interface{}, min int) ValidateFunc {
	return wrap(v.param, v.timeout, LengthMin(src, min))
}
func (v validations) LengthMax(src interface{}, max int) ValidateFunc {
	return wrap(v.param, v.timeout, LengthMax(src, max))
}
func (v validations) Length(src interface{}, min, max int) ValidateFunc {
	return wrap(v.param, v.timeout, Length(src, min, max))
}

func (v validations) Unique(arr ...interface{}) ValidateFunc {
	return wrap(v.param, v.timeout, Unique(arr...))
}

func (v validations) Email(email string) ValidateFunc {
	return wrap(v.param, v.timeout, Email(email))
}

func (v validations) True(b bool) ValidateFunc {
	return wrap(v.param, v.timeout, True(b))
}

func (v validations) False(b bool) ValidateFunc {
	return wrap(v.param, v.timeout, False(b))
}

func (v validations) Not(vf ValidateFunc) ValidateFunc {
	return wrap(v.param, v.timeout, Not(vf))
}

func (v validations) NotNullAll(list ...interface{}) ValidateFunc {
	return wrap(v.param, v.timeout, NotNullAll(list...))
}

func (v validations) NotNull(src interface{}) ValidateFunc {
	return wrap(v.param, v.timeout, NotNull(src))
}

func (v validations) NullAll(list ...interface{}) ValidateFunc {
	return wrap(v.param, v.timeout, NullAll(list...))
}

func (v validations) Null(src interface{}) ValidateFunc {
	return wrap(v.param, v.timeout, Null(src))
}

func (v validations) EqualOneOf(src interface{}, ofList ...interface{}) ValidateFunc {
	return wrap(v.param, v.timeout, EqualOneOf(src, ofList...))
}

func (v validations) EqualAll(first interface{}, with ...interface{}) ValidateFunc {
	return wrap(v.param, v.timeout, EqualAll(first, with...))
}

func (v validations) Equal(a interface{}, b interface{}) ValidateFunc {
	return wrap(v.param, v.timeout, Equal(a, b))
}

func (v validations) OR(validators ...Validator) ValidateFunc {
	return wrap(v.param, v.timeout, OR(validators...))
}

func (v validations) AND(validators ...Validator) ValidateFunc {
	return wrap(v.param, v.timeout, AND(validators...))
}

func (v validations) Async() AsyncValidations {
	return asyncValidations{
		param: v.param,
	}
}

func (v validations) Timeout(timeout time.Duration) Validations {
	return validations{
		param:   v.param,
		timeout: timeout,
	}
}
func (v validations) IntMin(i, min int) ValidateFunc {
	return wrap(v.param, v.timeout, IntMin(i, min))
}

func (v validations) IntMax(i, max int) ValidateFunc {
	return wrap(v.param, v.timeout, IntMax(i, max))
}

func (v validations) Int(i, min, max int) ValidateFunc {
	return wrap(v.param, v.timeout, Int(i, min, max))
}
func (v validations) Int64Min(i, min int64) ValidateFunc {
	return wrap(v.param, v.timeout, Int64Min(i, min))
}

func (v validations) Int64Max(i, max int64) ValidateFunc {
	return wrap(v.param, v.timeout, Int64Max(i, max))
}

func (v validations) Int64(i, min, max int64) ValidateFunc {
	return wrap(v.param, v.timeout, Int64(i, min, max))
}

func (v validations) Hex(str string) ValidateFunc {
	return wrap(v.param, v.timeout, Hex(str))
}

func (v validations) Base32(str string) ValidateFunc {
	return wrap(v.param, v.timeout, Base32(str))
}

func (v validations) IPv4(str string) ValidateFunc {
	return wrap(v.param, v.timeout, IPv4(str))
}

func (v validations) IPv6(str string) ValidateFunc {
	return wrap(v.param, v.timeout, IPv6(str))
}

func (v validations) IP(str string) ValidateFunc {
	return wrap(v.param, v.timeout, IP(str))
}

func (v validations) Name(str string) ValidateFunc {
	return wrap(v.param, v.timeout, Name(str))
}

func (v validations) Subdomain(str string) ValidateFunc {
	return wrap(v.param, v.timeout, Subdomain(str))
}

func (v validations) Code(str string) ValidateFunc {
	return wrap(v.param, v.timeout, Code(str))
}

func InstanceOf(a, b interface{}) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if reflect.TypeOf(a) != reflect.TypeOf(b) {
			return makeErr(param, "invalid instance type")
		}
		return nil
	}
}

var rxHex = regexp.MustCompile("^[a-fA-F0-9]+$")

func Hex(str string) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if !rxHex.MatchString(str) {
			return makeErr(param, "invalid hex fromat")
		}
		return nil
	}
}

var rxBase32 = regexp.MustCompile("^(?:[A-Z2-7]{8})*(?:[A-Z2-7]{2}={6}|[A-Z2-7]{4}={4}|[A-Z2-7]{5}={3}|[A-Z2-7]{7}=)?$")

func Base32(str string) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if !rxBase32.MatchString(str) {
			return makeErr(param, "invalid base32 fromat")
		}
		return nil
	}
}

func IPv4(str string) ValidateFunc {
	return func(ctx context.Context, param string) error {
		ip := net.ParseIP(str)
		if ip == nil || ip.To4() == nil {
			return makeErr(param, "invalid ipv4 fromat")
		}
		return nil
	}
}

func IPv6(str string) ValidateFunc {
	return func(ctx context.Context, param string) error {
		ip := net.ParseIP(str)
		if ip == nil || !strings.Contains(str, ":") {
			return makeErr(param, "invalid ipv6 fromat")
		}
		return nil
	}
}

func IP(str string) ValidateFunc {
	return func(ctx context.Context, param string) error {
		ip := net.ParseIP(str)
		if ip == nil {
			return makeErr(param, "invalid ipv fromat")
		}
		return nil
	}
}

func AND(validators ...Validator) ValidateFunc {
	return func(ctx context.Context, param string) error {
		for _, v := range validators {
			err := v.Validate(ctx, param)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func Unique(arr ...interface{}) ValidateFunc {
	return func(ctx context.Context, param string) error {
		for i, v := range arr {
			for j, vv := range arr {
				if i != j && v == vv {
					return makeErr(param, "duplicate found: %v", v)
				}
			}
		}
		return nil
	}
}

type AsyncValidations interface {
	OR(validators ...Validator) ValidateFunc
	AND(validators ...Validator) ValidateFunc
	Timeout(timeout time.Duration) AsyncValidations
}

func Async() AsyncValidations {
	return asyncValidations{}
}

type asyncValidations struct {
	param   string
	timeout time.Duration
}

func (v asyncValidations) Timeout(timeout time.Duration) AsyncValidations {
	return asyncValidations{
		param:   v.param,
		timeout: timeout,
	}
}

func (v asyncValidations) AND(validators ...Validator) ValidateFunc {
	return wrap(v.param, v.timeout, asyncAND(validators...))
}

func (v asyncValidations) OR(validators ...Validator) ValidateFunc {
	return wrap(v.param, v.timeout, asyncOR(validators...))
}

func asyncAND(validators ...Validator) ValidateFunc {
	return func(ctx context.Context, param string) error {
		errch := make(chan error)
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		count := len(validators)
		for _, v := range validators {
			go func(ch chan<- error, v Validator) {
				err := v.Validate(ctx, param)
				errch <- err
			}(errch, v)
		}
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-errch:
				if err != nil {
					return err
				}
				count--
				if count == 0 {
					return nil
				}
			}
		}
	}
}

func asyncOR(validators ...Validator) ValidateFunc {
	return func(ctx context.Context, param string) error {
		var merr *multierr.MultiErr
		errch := make(chan error)
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		count := len(validators)
		for _, v := range validators {
			go func(ch chan<- error, v Validator) {
				errch <- v.Validate(ctx, param)
			}(errch, v)
		}
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-errch:
				if err == nil {
					return err
				}
				if merr == nil {
					merr = multierr.From(err)
				} else {
					merr = merr.With(err)
				}
				count--
				if count == 0 {
					return merr
				}
			}
		}
	}
}

func OR(validators ...Validator) ValidateFunc {
	return func(ctx context.Context, param string) error {
		var merr *multierr.MultiErr
		for _, v := range validators {
			err := v.Validate(ctx, param)
			if err == nil {
				return nil
			}
			if merr == nil {
				merr = multierr.From(err)
			} else {
				merr = merr.With(err)
			}
		}
		return merr
	}
}

func LengthMin(src interface{}, min int) ValidateFunc {
	return func(ctx context.Context, param string) error {
		l, ok := length(src)
		if !ok {
			return makeErr(param, "invalid type for length check")
		}
		if l < min {
			return makeErr(param, "length min: %d", min)
		}
		return nil
	}
}

func LengthMax(src interface{}, max int) ValidateFunc {
	return func(ctx context.Context, param string) error {
		l, ok := length(src)
		if !ok {
			return makeErr(param, "invalid type for length check")
		}
		if l > max {
			return makeErr(param, "length max: %d", max)
		}
		return nil
	}
}

func Length(src interface{}, min, max int) ValidateFunc {
	return AND(
		LengthMin(src, min),
		LengthMax(src, max),
	)
}

func length(src interface{}) (int, bool) {
	rt := reflect.ValueOf(src)
	switch rt.Kind() {
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		fallthrough
	case reflect.String:
		return rt.Len(), true
	default:
		return 0, false
	}
}

type Comperable interface {
	Equals(interface{}) bool
}

func Equal(a interface{}, b interface{}) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if a != b {
			return makeErr(param, "%v != %v", a, b)
		}
		return nil
	}
}

func EqualAll(first interface{}, with ...interface{}) ValidateFunc {
	return func(ctx context.Context, param string) error {
		for _, v := range with {
			err := Equal(first, v)(ctx, param)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func EqualOneOf(src interface{}, ofList ...interface{}) ValidateFunc {
	return func(ctx context.Context, param string) error {
		for _, v := range ofList {
			err := Equal(src, v)(ctx, param)
			if err == nil {
				return nil
			}
		}
		return makeErr(param, "%v did not match any of %v", src, ofList)
	}
}

type Nullable interface {
	IsNull() bool
}

func Null(src interface{}) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if src == nil || (reflect.ValueOf(src).Kind() == reflect.Ptr && reflect.ValueOf(src).IsNil()) {
			return nil
		}
		imp, ok := src.(Nullable)
		if !ok || !imp.IsNull() {
			return makeErr(param, "should be null %v", src)
		}
		return nil
	}
}

func NullAll(list ...interface{}) ValidateFunc {
	return func(ctx context.Context, param string) error {
		for _, v := range list {
			err := Null(v)(ctx, param)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func NotNull(src interface{}) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if src == nil || (reflect.ValueOf(src).Kind() == reflect.Ptr && reflect.ValueOf(src).IsNil()) {
			return makeErr(param, "should not be null")
		}
		imp, ok := src.(Nullable)
		if ok && imp.IsNull() {
			return makeErr(param, "should not be null")
		}
		return nil
	}
}

func NotNullAll(list ...interface{}) ValidateFunc {
	return func(ctx context.Context, param string) error {
		for _, v := range list {
			err := NotNull(v)(ctx, param)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func Not(vf ValidateFunc) ValidateFunc {
	return func(ctx context.Context, param string) error {
		err := vf(ctx, param)
		if err != nil {
			return nil
		}
		return makeErr(param, "not")
	}
}

var rxEmail = regexp.MustCompile(
	"^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|" +
		"[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFE" +
		"F}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\" +
		"|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\" +
		"x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\" +
		"x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\" +
		"x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F" +
		"900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b" +
		"\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}" +
		"\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\" +
		"x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF" +
		"}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[" +
		"\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]" +
		")([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x" +
		"{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D" +
		"7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA" +
		"-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{F" +
		"FEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\" +
		"x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}" +
		"\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A" +
		"0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$",
)

func Email(email string) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if !rxEmail.MatchString(email) {
			return makeErr(param, "invalid email address: %s", email)
		}
		return nil
	}
}
func True(b bool) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if !b {
			return makeErr(param, "invalid value")
		}
		return nil
	}
}

func False(b bool) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if b {
			return makeErr(param, "invalid value")
		}
		return nil
	}
}

func RegExp(exp, str string) ValidateFunc {
	rx := regexp.MustCompile(exp)
	return func(ctx context.Context, param string) error {
		if !rx.MatchString(str) {
			return makeErr(param, "invalid regex match")
		}
		return nil
	}
}

func IntMin(i, min int) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if i < min {
			return makeErr(param, "%d < %d", i, min)
		}
		return nil
	}
}

func IntMax(i, max int) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if i > max {
			return makeErr(param, "%d > %d", i, max)
		}
		return nil
	}
}

func Int(i, min, max int) ValidateFunc {
	return AND(
		IntMin(i, min),
		IntMax(i, max),
	)
}

func Int64Min(i, min int64) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if i < min {
			return makeErr(param, "%d < %d", i, min)
		}
		return nil
	}
}

func Int64Max(i, max int64) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if i > max {
			return makeErr(param, "%d > %d", i, max)
		}
		return nil
	}
}

func Int64(i, min, max int64) ValidateFunc {
	return AND(
		Int64Min(i, min),
		Int64Max(i, max),
	)
}

var rxSubdomain = regexp.MustCompile("^[a-z0-9][a-z0-9-]+[a-z0-9]$")

func Subdomain(name string) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if !rxSubdomain.MatchString(name) {
			return makeErr(param, "invalid subdomain name: %s", name)
		}
		return nil
	}
}

var rxName = regexp.MustCompile("^[\\p{L}0-9\\s'.\\-+]{1,50}$")

func Name(name string) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if !rxName.MatchString(name) {
			return makeErr(param, "invalid name: %s", name)
		}
		return nil
	}
}

var rxCode = regexp.MustCompile("^[A-Z0-9]{6}$")

func Code(str string) ValidateFunc {
	return func(ctx context.Context, param string) error {
		if !rxCode.MatchString(str) {
			return makeErr(param, "invalid code: %s", str)
		}
		return nil
	}
}

func wrapTimeout(dur time.Duration, v Validator) ValidateFunc {
	return func(ctx context.Context, param string) error {
		ctx, cancel := context.WithTimeout(ctx, dur)
		defer cancel()
		errch := make(chan error)
		go func(ch chan<- error, v Validator) {
			ch <- v.Validate(ctx, param)
		}(errch, v)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errch:
			return err
		}
	}
}

func wrap(param string, timeout time.Duration, v Validator) ValidateFunc {
	if param != "" {
		v = wrapParam(param, v)
	}
	if timeout > 0 {
		v = wrapTimeout(timeout, v)
	}
	return func(ctx context.Context, param string) error {
		return v.Validate(ctx, param)
	}
}
