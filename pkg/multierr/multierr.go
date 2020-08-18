package multierr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Kind int

func (k Kind) HTTP() int {
	switch k {
	case InvalidArgument:
		return http.StatusBadRequest
	case NotFound:
		return http.StatusNotFound
	case AlreadyExists:
		return http.StatusConflict
	case PermissionDenied:
		return http.StatusForbidden
	case Unauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

const (
	Unexpected       Kind = 0
	InvalidArgument  Kind = 1
	NotFound         Kind = 2
	AlreadyExists    Kind = 3
	PermissionDenied Kind = 4
	Unauthorized     Kind = 5
)

var registry map[string]map[int]*MultiErr

func Register(arr ...*MultiErr) {
	for _, merr := range arr {
		if merr.Group() == "" || merr.Code() == 0 || merr.prev != nil {
			continue
		}
		gm, ok := registry[merr.Group()]
		if !ok {
			gm = make(map[int]*MultiErr)
			registry[merr.Group()] = gm
		}
		gm[merr.Code()] = merr
	}
}

func init() {
	registry = make(map[string]map[int]*MultiErr)
}

// MultiErr allows to stack multiple errors
type MultiErr struct {
	kind  Kind
	group string
	code  int
	value error
	prev  *MultiErr
}

// New create a new MultiErr with given format
func New(format string, a ...interface{}) *MultiErr {
	return &MultiErr{
		value: fmt.Errorf(format, a...),
		prev:  nil,
	}
}

type Builder struct {
	group string
	kind  Kind
	code  int
}

func Group(group string) Builder {
	return Builder{
		group: group,
		kind:  Unexpected,
		code:  0,
	}
}

func (builder Builder) Kind(kind Kind) Builder {
	return Builder{
		group: builder.group,
		kind:  kind,
		code:  builder.code,
	}
}

func (builder Builder) Code(code int) Builder {
	return Builder{
		group: builder.group,
		kind:  builder.kind,
		code:  code,
	}
}

func (builder Builder) New(format string, a ...interface{}) *MultiErr {
	return &MultiErr{
		group: builder.group,
		code:  builder.code,
		kind:  builder.kind,
		value: fmt.Errorf(format, a...),
		prev:  nil,
	}
}

// From converts error to MultiErr
func From(err error) *MultiErr {
	if err == nil {
		return nil
	}
	mr, ok := err.(*MultiErr)
	if !ok {
		mr = &MultiErr{
			value: err,
			prev:  nil,
		}
	}
	return mr
}

func format(err *MultiErr) string {
	if err == nil || err.value == nil {
		return ""
	}
	msg := err.value.Error()
	if err.group != "" {
		return fmt.Sprintf("[group:%s,code:%d] %s", err.group, err.code, msg)
	}
	return msg
}

// Error implements error interface
func (mr *MultiErr) Error() string {
	if mr == nil {
		return ""
	}
	list := make([]string, 0)
	for err := mr; err != nil; err = err.prev {
		list = append([]string{format(err)}, list...)
	}
	return strings.Join(list, "; ")
}

func (mr MultiErr) Kind() Kind {
	return mr.kind
}

func (mr MultiErr) Group() string {
	return mr.group
}

func (mr MultiErr) Code() int {
	return mr.code
}

func (mr *MultiErr) MarshalJSON() ([]byte, error) {
	raw := encJSONErr(mr)
	return json.Marshal(raw)
}

func (mr *MultiErr) UnmarshalJSON(src []byte) error {
	raw := new(jsonError)
	err := json.Unmarshal(src, raw)
	if err != nil {
		return err
	}
	tmp := decJSONErr(raw)
	*mr = *tmp
	return nil
}

// Has checks if MultiErr contains given error
func (mr *MultiErr) Has(err error) bool {
	for e := mr; e != nil; e = e.prev {
		if e.Equal(err) {
			return true
		}
	}
	return false
}

// Equal checks if MultiErr is equal to given error
func (mr *MultiErr) Equal(err error) bool {
	// nil safe-check mr == err == nil
	if mr == nil && err == nil {
		return true
	} else if mr == nil {
		return false
	} else if err == nil {
		return false
	}
	target := From(err)
	return target.value == mr.value
}

func clone(err *MultiErr) *MultiErr {
	if err == nil {
		return nil
	}
	tmp := *err
	if tmp.prev != nil {
		tmp.prev = clone(tmp.prev)
	}
	return &tmp
}

// With combines mr with given error and returns new MultiErr
func (mr *MultiErr) With(err error) *MultiErr {
	if mr == nil && err == nil {
		return nil
	}
	if err == nil {
		return mr
	}
	child := From(err)
	child = clone(child)
	child.prev = mr
	return child
}
