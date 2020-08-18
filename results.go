package czzle

type Results struct {
	OK          bool    `json:"ok"`
	Took        int64   `json:"took"`
	Next        *Puzzle `json:"next,omitempty"`
	AccessToken string  `json:"access_token,omitempty"`
}

func (r *Results) IsOK() bool {
	if r == nil {
		return false
	}
	return r.OK
}

func (r *Results) SetOK(ok bool) {
	if r == nil {
		return
	}
	r.OK = ok
}

func (r *Results) GetTook() int64 {
	if r == nil {
		return 0
	}
	return r.Took
}

func (r *Results) SetTook(took int64) {
	if r == nil {
		return
	}
	r.Took = took
}

func (r *Results) GetNext() *Puzzle {
	if r == nil {
		return nil
	}
	return r.Next
}

func (r *Results) SetNext(next *Puzzle) {
	if r == nil {
		return
	}
	r.Next = next
}

func (r *Results) GetAccessToken() string {
	if r == nil {
		return ""
	}
	return r.AccessToken
}

func (r *Results) HasAccessToken() bool {
	if r == nil {
		return false
	}
	return r.AccessToken != ""
}

func (r *Results) SetAccessToken(tkn string) {
	if r == nil {
		return
	}
	r.AccessToken = tkn
}
