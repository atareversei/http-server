package http

import (
	"errors"
	"fmt"
)

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusNotFound            StatusCode = 404
	StatusMethodNotAllowed    StatusCode = 405
	StatusInternalServerError StatusCode = 500
)

const (
	okMessage                  = "OK"
	badRequestMessage          = "Bad Request"
	notFoundMessage            = "Not Found"
	methodNotAllowedMessage    = "Method Not Allowed"
	internalServerErrorMessage = "Internal Server Error"
)

func (s StatusCode) String() string {
	switch s {
	case StatusOk:
		return okMessage
	case StatusBadRequest:
		return badRequestMessage
	case StatusNotFound:
		return notFoundMessage
	case StatusMethodNotAllowed:
		return methodNotAllowedMessage
	case StatusInternalServerError:
		return internalServerErrorMessage
	default:
		return ""

	}
}

var MalformedBodyError = errors.New("body couldn't be parsed")

type Method string

const (
	MethodGet     Method = "GET"
	MethodHead    Method = "HEAD"
	MethodPost    Method = "POST"
	MethodPut     Method = "PUT"
	MethodPatch   Method = "PATCH"
	MethodDelete  Method = "DELETE"
	MethodConnect Method = "CONNECT"
	MethodOptions Method = "OPTIONS"
	MethodTrace   Method = "TRACE"
)

func IsMethodValid(s string) (Method, error) {
	switch s {
	case "GET":
		return MethodGet, nil
	case "HEAD":
		return MethodHead, nil
	case "POST":
		return MethodPost, nil
	case "PUT":
		return MethodPut, nil
	case "PATCH":
		return MethodPatch, nil
	case "DELETE":
		return MethodDelete, nil
	case "CONNECT":
		return MethodConnect, nil
	case "OPTIONS":
		return MethodOptions, nil
	case "TRACE":
		return MethodTrace, nil
	default:
		return "", fmt.Errorf("unsupported method %q", s)
	}
}

type Version string

const (
	Version10 Version = "HTTP/1.0"
	Version11 Version = "HTTP/1.1"
)

func IsVersionValid(v string) (Version, error) {
	switch v {
	case "HTTP/1.0":
		return Version10, nil
	case "HTTP/1.1":
		return Version11, nil
	default:
		return "", fmt.Errorf("unsupported version %q", v)
	}
}
