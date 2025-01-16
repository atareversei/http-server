package server

import (
	"fmt"
	"github.com/atareversei/network-course-projects/pkg/cli"
	"github.com/atareversei/network-course-projects/pkg/http"
	"net"
)

type Server struct {
	port   int
	router map[string]map[string]HandlerFunc
}

func New(port int) Server {
	return Server{port: port, router: make(map[string]map[string]HandlerFunc)}
}

func (s *Server) Start() {
	cli.Success(fmt.Sprintf("tcp server started at :%d", s.port))
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		cli.Error("server could not be started", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			cli.Error("connection resulted in an error", err)
			continue
		}
		go s.parseHTTP(conn)
	}
}

type HandlerFunc func(req http.Request, res *http.Response)

func (s *Server) All(pattern string, handler func(req http.Request, res *http.Response)) {
	s.router[pattern]["all"] = handler
}

func (s *Server) Get(pattern string, handler func(req http.Request, res *http.Response)) {
	s.router[pattern]["get"] = handler
}

func (s *Server) Post(pattern string, handler func(req http.Request, res *http.Response)) {
	s.router[pattern]["post"] = handler
}

func (s *Server) parseHTTP(conn net.Conn) {
	httpRequest := http.NewRequest(conn)
	httpRequest.Parse()
	response := http.NewResponse(conn, httpRequest.Version())
	s.handleRequest(httpRequest, &response)
	conn.Close()
}

func (s *Server) handleRequest(request http.Request, response *http.Response) {
	resource, resOk := s.router[request.Path()]
	if !resOk {
		response.WriteHeader(404)
		return
	}
	handler, handlerOk := resource[request.Method()]
	if !handlerOk {
		handler, allOk := resource["all"]
		if allOk {
			handler(request, response)
			return
		}
		response.WriteHeader(http.MethodNotAllowed)
		return
	}
	handler(request, response)
}
