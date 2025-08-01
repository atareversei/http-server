package http

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/atareversei/http-server/internal/cli"
)

// Request is the object that gets populated when TCP
// connection hits the server.
type Request struct {
	// method can be GET and POST.
	// TODO - implement CONNECT, HEAD, PUT, PATCH, and DELETE
	method string
	// path is the resource that the request is trying
	// to retrieve from the server.
	// e.g. www.example.com/user?page=10 -> path: /user
	path string
	// version is the HTTP version that the request has
	// been created in.
	version string
	// headers holds most of the headers that the request
	// has, except the headers that have occurred multiple times.
	// TODO - support multiple headers
	headers map[string]string
	// body holds the body of the request if it is not GET.
	body []byte
	// params holds the query parameters of the request.
	// e.g. www.example.com/user?page=10 -> the map holds
	// a key of page with the value of 10 (string)
	params map[string]string
	// conn holds the TCP connection information.
	// required for reading the request data.
	conn net.Conn
}

// newRequestFromTCPConn creates a new request struct that can be used
// to invoke receiver functions to populate the struct.
func newRequestFromTCPConn(conn net.Conn) Request {
	return Request{conn: conn}
}

// Parse is used to parse a TCP byte stream into
// an HTTP Request struct.
func (req *Request) Parse() {
	r := bufio.NewReader(req.conn)
	req.parseStartLine(r)
	req.parseHeaders(r)
	req.parseBody(r)
}

// parseStartLine parses the very first line of an HTTP request.
// e.g. GET /users?page=10 HTTP/1.1
func (req *Request) parseStartLine(r *bufio.Reader) {
	request, err := r.ReadString('\n')
	if err != nil {
		cli.Error("failed to read the request", err)
		return
	}
	requestParts := strings.Split(request, " ")
	fullPath := strings.TrimSpace(requestParts[1])
	queryParamStart := strings.Index(fullPath, "?")
	var path string
	var queryString string
	if queryParamStart != -1 {
		path = fullPath[:queryParamStart]
		queryString = fullPath[queryParamStart+1:]
	} else {
		path = fullPath
		queryString = ""
	}
	req.parseQueryParams(queryString)
	req.path = path
	req.method = strings.TrimSpace(requestParts[0])
	req.version = strings.TrimSpace(requestParts[2])
}

// parseQueryParams parses the query string of an HTTP request.
// e.g. /users?page=10&per_page=&sort& will result in 3 entries:
// "page": "10"
// "per_page": ""
// "sort": ""
func (req *Request) parseQueryParams(qs string) {
	qp := make(map[string]string)
	queries := strings.Split(qs, "&")
	for _, query := range queries {
		if query == "" {
			continue
		}
		queryParts := strings.Split(query, "=")
		if queryParts[0] == "" {
			continue
		}
		if len(queryParts) == 1 || queryParts[1] == "" {
			qp[queryParts[0]] = ""
		} else {
			qp[queryParts[0]] = queryParts[1]
		}
	}
	req.params = qp
}

// parseHeaders parses the headers of an HTTP request.
func (req *Request) parseHeaders(r *bufio.Reader) {
	if req.headers == nil {
		req.headers = make(map[string]string)
	}
	for {
		header, err := r.ReadString('\n')
		header = strings.TrimSpace(header)
		if err != nil {
			cli.Error(fmt.Sprintf("couldn't read header, skipping header (%s)", header), err)
			continue
		}
		// empty line which indicates the header section is over
		if header == "" {
			break
		}
		headerParts := strings.SplitN(header, ":", 2)
		if len(headerParts) != 2 {
			cli.Error(fmt.Sprintf("malformed header: %s", header), fmt.Errorf("not abiding by key:value format"))
			continue
		}
		req.headers[headerParts[0]] = headerParts[1]
	}
}

// parseBody parses the body of an HTTP request if there is any.
func (req *Request) parseBody(r *bufio.Reader) {
	if req.method == "GET" {
		return
	}
	contentLengthValue, ok := req.headers["Content-Length"]
	if !ok {
		cli.Error("couldn't find `Content-Length` in the headers", fmt.Errorf("headers map returned `!ok` for `Content-Length`"))
		return
	}
	contentLength, err := strconv.Atoi(contentLengthValue)
	if err != nil {
		cli.Error("couldn't parse `Content-Length` out of headers", err)
		return
	}
	if contentLength == 0 {
		return
	}
	body := make([]byte, contentLength)
	_, err = r.Read(body)
	req.body = body
}

// Method returns the method of the request.
func (req *Request) Method() string {
	return req.method
}

// Path returns the path of the request.
func (req *Request) Path() string {
	return req.path
}

// Version returns the version of the request.
func (req *Request) Version() string {
	return req.version
}

// Header returns the headers of the request.
func (req *Request) Header(key string) (string, bool) {
	h, ok := req.headers[key]
	return h, ok
}

// Params returns the query parameters of a request.
func (req *Request) Params(key string) (string, bool) {
	v, ok := req.params[key]
	return v, ok
}
