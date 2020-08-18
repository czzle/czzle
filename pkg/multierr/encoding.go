package multierr

import (
	"encoding/json"
	"errors"
	"net/http"
)

type jsonErrorRes struct {
	Error *jsonError `json:"error"`
}

type jsonError struct {
	Code    int        `json:"code"`
	Group   string     `json:"group"`
	Message string     `json:"message"`
	Prev    *jsonError `json:"prev"`
}

func encJSONErr(merr *MultiErr) *jsonError {
	if merr == nil || merr.value == nil {
		return nil
	}
	return &jsonError{
		Code:    merr.code,
		Group:   merr.group,
		Message: merr.value.Error(),
		Prev:    encJSONErr(merr.prev),
	}
}

func decJSONErr(merr *jsonError) *MultiErr {
	if merr == nil {
		return nil
	}
	parent := decJSONErr(merr.Prev)
	res := &MultiErr{
		value: errors.New(merr.Message),
		prev:  parent,
	}
	g, ok := registry[merr.Group]
	if !ok {
		return res
	}
	err, ok := g[merr.Code]
	if !ok {
		return res
	}
	res = clone(err)
	res.prev = parent
	return res
}

func catchFirst(root *MultiErr) *MultiErr {
	for merr := root; merr != nil; merr = merr.prev {
		if merr.kind != 0 && merr.group != "" {
			return merr
		}
	}
	return root
}

func ToHTTP(w http.ResponseWriter, err error) {
	merr := From(err)
	first := catchFirst(merr)
	w.WriteHeader(first.kind.HTTP())
	raw := encJSONErr(merr)
	json.NewEncoder(w).Encode(jsonErrorRes{raw})
}

func FromHTTP(src []byte) error {
	if src == nil {
		return errors.New("empty src for parsing error")
	}
	res := new(jsonErrorRes)
	json.Unmarshal(src, res)
	if res == nil || res.Error == nil {
		return nil
	}
	merr := decJSONErr(res.Error)
	return merr
}
