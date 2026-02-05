package handler

import (
	"github.com/atareversei/http-server/http"
)

type Handler struct{}

func New() Handler {
	return Handler{}
}

func (h Handler) Ping(req http.Request, res http.Response) {}
