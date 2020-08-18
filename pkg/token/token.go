package token

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/czzle/czzle"
	"github.com/czzle/czzle/pkg/uuid"
)

var (
	ErrMissingSignature = errors.New("missing signature")
	ErrInvalidVersion   = errors.New("invalid version")
	ErrMalformed        = errors.New("malformed token")
)

type Version int

const (
	V1 Version = 1
)

const sep = "."

var encoding = base64.RawURLEncoding.WithPadding(base64.NoPadding)

type Token struct {
	raw     string
	header  header
	payload Payload
}

type header struct {
	Version Version `json:"ver"`
}

func (tkn Token) Payload() Payload {
	return tkn.payload
}

type Payload struct {
	Solved    bool  `json:"solved"`
	IssuedAt  int64 `json:"iat"`
	ExpiresAt int64 `json:"exp"`

	Level     czzle.Level `json:"lvl"`
	ClientID  uuid.UUID   `json:"cid"`
	ClientIP  string      `json:"cip"`
	Generator string      `json:"gen"`
}

func (tkn Token) String() string {
	return tkn.raw
}

func (tkn *Token) Sign(secret string) error {
	if tkn.header.Version != V1 {
		return ErrInvalidVersion
	}
	secret = strings.TrimSpace(secret)
	hd, err := json.Marshal(tkn.header)
	if err != nil {
		return err
	}
	hdStr := encoding.EncodeToString(hd)
	pl, err := json.Marshal(tkn.payload)
	if err != nil {
		return err
	}
	plStr := encoding.EncodeToString(pl)
	msg := strings.Join([]string{hdStr, plStr}, sep)
	var sigStr string
	{
		var msgh []byte
		{
			h := sha256.New()
			h.Write([]byte(msg))
			msgh = h.Sum(nil)
		}
		var key []byte
		{
			h := sha256.New()
			h.Write([]byte(secret))
			key = h.Sum(nil)
		}
		h := sha256.New()
		h.Write(msgh)
		h.Write(key)
		sig := h.Sum(nil)
		sigStr = encoding.EncodeToString(sig)
	}
	tkn.raw = strings.Join([]string{msg, sigStr}, sep)
	return nil
}

func (tkn Token) Validate(secret string) bool {
	if tkn.raw == "" {
		return false
	}
	splits := strings.Split(tkn.raw, sep)
	msg := strings.Join([]string{splits[0], splits[1]}, sep)
	var sigStr string
	{
		var msgh []byte
		{
			h := sha256.New()
			h.Write([]byte(msg))
			msgh = h.Sum(nil)
		}
		var key []byte
		{
			h := sha256.New()
			h.Write([]byte(secret))
			key = h.Sum(nil)
		}
		h := sha256.New()
		h.Write(msgh)
		h.Write(key)
		sig := h.Sum(nil)
		sigStr = encoding.EncodeToString(sig)
	}
	return sigStr == splits[2]
}

func New(payload Payload) Token {
	return Token{
		header: header{
			Version: V1,
		},
		payload: payload,
	}
}

func Parse(str string) (tkn Token, err error) {
	splits := strings.Split(str, sep)
	if len(splits) != 3 {
		err = ErrMalformed
		return
	}
	hdData, err := encoding.DecodeString(splits[0])
	if err != nil {
		err = ErrMalformed
		return
	}
	var hd header
	err = json.Unmarshal(hdData, &hd)
	if err != nil {
		err = ErrMalformed
		return
	}
	if hd.Version != V1 {
		err = ErrInvalidVersion
		return
	}
	plData, err := encoding.DecodeString(splits[1])
	if err != nil {
		err = ErrMalformed
		return
	}
	var pl Payload
	err = json.Unmarshal(plData, &pl)
	if err != nil {
		err = ErrMalformed
		return
	}
	sigData, err := encoding.DecodeString(splits[2])
	if err != nil || len(sigData) != 32 {
		err = ErrMalformed
		return
	}
	tkn = Token{
		raw:     str,
		header:  hd,
		payload: pl,
	}
	return

}
