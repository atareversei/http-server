package http

import (
	"fmt"
	"io"
	"net"
)

// Method represents an HTTP request method (e.g., GET, POST).
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

// IsMethodValid checks if a string corresponds to a supported HTTP method.
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

// String returns the string representation of the Method.
func (m Method) String() string {
	return string(m)
}

// handleHeadMethod handles a HEAD request. The server should behave as if it's
// a GET request but must not send a response body. The streaming response writer
// handles this implicitly by sending headers first.
func handleHeadMethod(req Request, res Response, resource map[Method]Handler) {
	handler, handlerOk := resource[MethodGet]
	if !handlerOk {
		HTTPError(res, StatusNotFound)
		return
	}
	handler.ServeHTTP(req, res)
}

// handleOptionsMethod handles an OPTIONS request, which is used for both
// general server capability discovery and for CORS preflight requests.
func handleOptionsMethod(req Request, res Response, router *DefaultRouter) {
	isCORSAllowed := router.cors.isCORSAllowed(req)
	if !isCORSAllowed {
		return
	}

	res.SetStatus(StatusNoContent)
	if req.Path() == "*" {
		res.SetHeader("Allow", router.getAllAvailableMethodsHeader())
	} else {
		res.SetHeader("Allow", router.getAvailableMethodsForResourceHeader(req.Path()))
	}

	res.SetHeader("Access-Control-Allow-Origin", router.cors.getAllowedOriginsHeader())
	res.SetHeader("Access-Control-Allow-Methods", router.cors.getAllowedMethodsHeader())
	res.SetHeader("Access-Control-Allow-Headers", router.cors.getAllowedHeadersHeader())
	res.SetHeader("Access-Control-Allow-Credentials", router.cors.getAllowedCredentialsHeader())
	res.SetHeader("Access-Control-Max-Age", router.cors.getAllowedMaxAgeHeader())
}

// handleConnectMethod handles a CONNECT request, establishing a two-way tunnel
// between the client and the requested destination. This is primarily used for HTTPS proxies.
func handleConnectMethod(req Request, res Response) {
	path := req.Path()
	srcConn := req.conn
	dstConn, err := net.Dial("tcp", path)
	if err != nil {
		HTTPErrorWithMessage(res, StatusServiceUnavailable, "Service Unavailable", "service refused to connect")
		return
	}

	res.SetStatusWithMessage(StatusOk, "Connection Established")

	go func() {
		defer dstConn.Close()
		defer srcConn.Close()
		io.Copy(dstConn, srcConn)
	}()

	io.Copy(srcConn, dstConn)
}
