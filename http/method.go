package http

import "fmt"

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

func IsMethodValid(m string) (Method, error) {
	switch m {
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
		return "", fmt.Errorf("unsupported method %q", m)
	}
}

func (m Method) String() string {
	return string(m)
}

func handleHeadMethod(req Request, res Response, resource map[Method]Handler) {
	handler, handlerOk := resource[MethodGet]
	if !handlerOk {
		HTTPError(res, StatusNotFound)
		return
	}
	handler.ServeHTTP(req, res)
}

func handlerOptionsMethod(req Request, res Response, resource map[Method]Handler) {
	// TODO: implement
}
