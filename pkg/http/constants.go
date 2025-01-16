package http

import "errors"

type StatusCode int

const (
	NotFound         StatusCode = 404
	MethodNotAllowed            = 405
)

const (
	notFoundMessage         = "Not Found"
	methodNotAllowedMessage = "Method Not Allowed"
)

func (s StatusCode) String() string {
	switch s {
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
