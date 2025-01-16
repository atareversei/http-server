package http

import (
	"bufio"
	"fmt"
	"github.com/atareversei/network-course-projects/pkg/cli"
	"net"
	"strconv"
	"strings"
)

// Request is the object that gets populated when tcp
// connection hits the server.
type Request struct {
	// method can be GET and POST.
	method string
	// path is the resource that the request is trying
	// to get on the server.
	// e.g. www.example.com/user?page=10 -> path: /user
	path string
	// version is the HTTP version that the request has
	// been created with.
	version string
	// headers holds most of the headers that the request
	// has, expect headers that have occurred multiple times.
	headers map[string]string
	// body holds the body of the request if it is not GET.
	body []byte
	// params holds the query parameters of the request.
	// e.g. www.example.com/user?page=10 -> the map holds
	// a key of page with the value of 10 (string)
	params map[string]string
	// conn holds the tcp connection information.
	// required for reading the request data.
	conn net.Conn
}

// NewRequest creates a new request struct, that can be used
// to invoke receiver functions to populate the struct.
func NewRequest(conn net.Conn) Request {
	return Request{conn: conn}
}

// Parse is used to parse a tcp byte stream into
// HTTP Request struct.
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
	path := fullPath[:queryParamStart]
	queryString := fullPath[queryParamStart+1:]
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

func (req *Request) Method() string {
	return req.method
}

func (req *Request) Path() string {
	return req.path
}

func (req *Request) Version() string {
	return req.version
}

func (req *Request) Header() map[string]string {
	return req.headers
}

func (req *Request) Params() map[string]string {
	return req.params
}
