package http

import "errors"

type StatusCode int

const (
	Ok               StatusCode = 200
	NotFound         StatusCode = 404
	MethodNotAllowed StatusCode = 405
)

const (
	okMessage               = "OK"
	notFoundMessage         = "Not Found"
	methodNotAllowedMessage = "Method Not Allowed"
)

func (s StatusCode) String() string {
	switch s {
	case Ok:
		return okMessage
	case NotFound:
		return notFoundMessage
	case MethodNotAllowed:
		return methodNotAllowedMessage
	// This should never happen
	default:
		return ""

	}
}

var MalformedBodyError = errors.New("body couldn't be parsed")
