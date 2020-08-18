package app

import "net/http"

type Web struct {
	root http.Handler
}
