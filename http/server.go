package http

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/atareversei/http-server/internal/cli"
)

var defaultRouter Router

// Server is used to spawn an HTTP server.
type Server struct {
	// Domain indicates on which domain the server will start. If empty, "localhost:" (127.0.0.1:)
	// will be used.
	Domain string
	// Port indicates on which port the server will start. If empty, ":http" (:80)
	// will be used.
	Port   int
	Router Router
	// fileHandler holds the information about file handlers.
	// The structure can be simplified as [REQUEST_PATH]FILESYSTEM_DIRECTORY_PATH
	fileHandler map[string]string
}

// New returns a server structure.
func New(port int) Server {
	return Server{Port: port, Router: make(map[string]map[string]HandlerFunc), fileHandler: make(map[string]string)}
}

// Start is used to listen on the specified port.
func (s *Server) Start() {
	cli.Success(fmt.Sprintf("tcp server started at :%d", s.Port))
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		cli.Error("server could not be started", err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			cli.Error("connection resulted in an error", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// FileHandler is file handler registrar.
func (s *Server) FileHandler(pattern string, directory string) {
	s.fileHandler[pattern] = directory
}

// handleConnection is used to parse streams of bytes received from a
// TCP connection into a meaningful HTTP request.
func (s *Server) handleConnection(conn net.Conn) {
	httpRequest := NewRequest(conn)
	httpRequest.Parse()
	response := NewResponse(conn, httpRequest.Version())
	s.handleRequest(httpRequest, &response)
	conn.Close()
}

// handleRequest is used to handle the routing phase of the request and
// deliver the request-response information to the right handler.
// TODO - refactor the code
func (s *Server) handleRequest(request Request, response *Response) {
	cli.Info(fmt.Sprintf("%s %s", request.Method(), request.Path()))
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
		resource, resOk := s.Router[request.Path()]
		if !resOk {
			// TODO - add proper response
			response.WriteHeader(NotFound)
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
			response.WriteHeader(MethodNotAllowed)
			return
		}
		handler(request, response)
	}
}
