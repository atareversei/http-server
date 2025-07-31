package http

import "errors"

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusNotFound            StatusCode = 404
	StatusMethodNotAllowed    StatusCode = 405
	StatusInternalServerError StatusCode = 500
)

const (
	okMessage               = "OK"
	notFoundMessage         = "Not Found"
	methodNotAllowedMessage = "Method Not Allowed"
	internalServerError     = "Internal Server Error"
)

func (s StatusCode) String() string {
	switch s {
	case StatusOk:
		return okMessage
	case StatusNotFound:
		return notFoundMessage
	case StatusMethodNotAllowed:
		return methodNotAllowedMessage
	case StatusInternalServerError:
		return internalServerError
	default:
		return ""

	}
}

var MalformedBodyError = errors.New("body couldn't be parsed")
