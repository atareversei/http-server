package server

import (
	"fmt"
	"github.com/atareversei/network-course-projects/pkg/cli"
	"github.com/atareversei/network-course-projects/pkg/http"
	"io"
	"net"
	"os"
	"strings"
)

type Server struct {
	port        int
	router      map[string]map[string]HandlerFunc
	fileHandler map[string]string
}

func New(port int) Server {
	return Server{port: port, router: make(map[string]map[string]HandlerFunc), fileHandler: make(map[string]string)}
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
	s.checkResourceEntry(pattern)
	s.router[pattern]["ALL"] = handler
}

func (s *Server) Get(pattern string, handler func(req http.Request, res *http.Response)) {
	s.checkResourceEntry(pattern)
	s.router[pattern]["GET"] = handler
}

func (s *Server) Post(pattern string, handler func(req http.Request, res *http.Response)) {
	s.checkResourceEntry(pattern)
	s.router[pattern]["POST"] = handler
}

func (s *Server) FileHandler(pattern string, directory string) {
	s.fileHandler[pattern] = directory
}

func (s *Server) parseHTTP(conn net.Conn) {
	httpRequest := http.NewRequest(conn)
	httpRequest.Parse()
	response := http.NewResponse(conn, httpRequest.Version())
	s.handleRequest(httpRequest, &response)
	conn.Close()
}

func (s *Server) handleRequest(request http.Request, response *http.Response) {
	// TODO - enhance type checking
	// simple check to see if the request points to a file
	if strings.Contains(request.Path(), ".") {
		// TODO - enhance inefficient lookup
		for k, _ := range s.fileHandler {
			if strings.Contains(request.Path(), k) {
				i := strings.Index(request.Path(), k)
				filePath := request.Path()[i+len(k):]
				directoryPath := s.fileHandler[k]
				fullPath := directoryPath + filePath
				f, err := os.Open(fullPath)
				// TODO - check for 404 and faulty files
				if err != nil {
					response.Write([]byte{})
				}
				defer f.Close()
				buf := make([]byte, 1024)
				for {
					n, err := f.Read(buf)
					if err != nil && err != io.EOF {
						// TODO - add proper response
						response.Write([]byte{})
					}
					if n == 0 {
						break
					}
					response.Write(buf[:n])
				}
			}
		}
	} else {
		fmt.Println(request.Path())
		resource, resOk := s.router[request.Path()]
		if !resOk {
			// TODO - add proper response
			response.WriteHeader(http.NotFound)
			return
		}
		handler, handlerOk := resource[strings.ToUpper(request.Method())]
		if !handlerOk {
			catchAll, allOk := resource["ALL"]
			if allOk {
				catchAll(request, response)
				return
			}
			// TODO - add proper response
			response.WriteHeader(http.MethodNotAllowed)
			return
		}
		handler(request, response)
	}
}

func (s *Server) checkResourceEntry(pattern string) {
	_, ok := s.router[pattern]
	if !ok {
		s.router[pattern] = make(map[string]HandlerFunc)
	}
}
