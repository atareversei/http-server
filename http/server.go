package http

import (
	"fmt"
	"io"
	"mime"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atareversei/http-server/internal/cli"
)

type Logger interface {
	Success(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string, err error)
}

type noOpLogger struct{}

func (n *noOpLogger) Success(msg string)          {}
func (n *noOpLogger) Info(msg string)             {}
func (n *noOpLogger) Warning(msg string)          {}
func (n *noOpLogger) Error(msg string, err error) {}

type Server struct {
	Host        string
	Port        int
	Router      Router
	Middlewares []Middleware
	Static      map[string]string

	// TLSCertFile string
	// TLSKeyFile	 sting

	loggingEnabled bool
	previousLogger Logger
	Logger         Logger

	// httpServer      *http.Server
	shutdownTimeout time.Duration
}

func (s *Server) DisableLogging() {
	s.previousLogger = s.Logger
	s.Logger = &noOpLogger{}
	s.loggingEnabled = false
}

func (s *Server) EnableLogging() {
	s.Logger = s.previousLogger
	s.loggingEnabled = true
}

type Router interface {
	Register(method string, path string, handler Handler)
	Handler
}

type Handler interface {
	ServeHTTP(req Request, res Response)
}

type DefaultRouter struct {
	routes map[string]map[string]Handler
	logger Logger
}

func (s *Server) NewRouter() DefaultRouter {
	return DefaultRouter{
		routes: make(map[string]map[string]Handler),
		logger: s.Logger,
	}
}

func (dr *DefaultRouter) Register(method string, path string, handler Handler) {
	m := strings.ToUpper(method)
	switch m {
	case "GET":
		dr.Get(path, handler)
	case "POST":
		dr.Post(path, handler)
	default:
		dr.logger.Warning("Unknown method: handler wasn't registered")
	}
}

func (dr *DefaultRouter) All(path string, handler Handler) {
	dr.checkResourceEntry(path)
	dr.routes[path]["*"] = handler
}

func (dr *DefaultRouter) Get(path string, handler Handler) {
	dr.checkResourceEntry(path)
	dr.routes[path]["GET"] = handler
}

func (dr *DefaultRouter) Post(path string, handler Handler) {
	dr.checkResourceEntry(path)
	dr.routes[path]["POST"] = handler
}

// checkResourceEntry is used to initialize the inner map of a router
// if it has not yet been initialized.
func (dr *DefaultRouter) checkResourceEntry(path string) {
	_, ok := dr.routes[path]
	if !ok {
		dr.routes[path] = make(map[string]Handler)
	}
}

func (dr *DefaultRouter) ServeHTTP(req Request, res Response) {
	resource, resOk := dr.routes[req.Path()]
	if !resOk {
		res.WriteHeader(StatusNotFound)
		res.SetHeader("Content-Type", "text/html")
		res.Write([]byte("<h1>404 Not Found</h1>"))
		return
	}
	handler, handlerOk := resource[strings.ToUpper(req.Method())]
	if !handlerOk {
		catchAll, allOk := resource["*"]
		if allOk {
			catchAll.ServeHTTP(req, res)
			return
		}
		res.WriteHeader(StatusMethodNotAllowed)
		res.SetHeader("Content-Type", "text/html")
		res.Write([]byte("<h1>405 Method Not Allowed</h1>"))
		return
	}
	handler.ServeHTTP(req, res)
}

type Middleware func(Handler) Handler

func New(port int, router Router) Server {
	return Server{
		Port:   port,
		Router: router,
	}
}

func (s *Server) Start() {
	cli.MadeInBasliqLabs()
	s.Logger.Success(fmt.Sprintf("tcp server is starting at :%d", s.Port))
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		s.Logger.Error("server could not be started", err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			s.Logger.Error("connection resulted in an error", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) FileHandler(pattern string, directory string) {
	s.Static[pattern] = directory
}

func (s *Server) handleConnection(conn net.Conn) {
	request := newRequestFromTCPConn(conn)
	request.Parse()
	response := newResponse(conn, request.Version())
	s.handleRequest(request, response)
	// TODO: check for errors
	conn.Close()
}

func (s *Server) handleRequest(request Request, response Response) {
	s.Logger.Info(fmt.Sprintf("%s %s", request.Method(), request.Path()))

	for prefix, _ := range s.Static {
		if strings.Contains(request.Path(), prefix) {
			s.handleFileRequest(prefix, request, response)
			return
		}
	}

	s.handleHttpRequest(request, response)
}

func (s *Server) handleFileRequest(prefix string, request Request, response Response) {
	i := strings.Index(request.Path(), prefix)
	filePath := request.Path()[i+len(prefix):]
	fullPath := s.Static[prefix] + filePath
	f, err := os.Open(fullPath)

	if os.IsNotExist(err) {
		response.WriteHeader(StatusNotFound)
		response.SetHeader("Content-Type", "text/html")
		response.Write([]byte("<h1>404 Not Found</h1>"))
		return
	} else if err != nil {
		response.WriteHeader(StatusInternalServerError)
		response.SetHeader("Content-Type", "text/html")
		response.Write([]byte("<h1>500 Internal Server Error</h1>"))
		return
	}

	defer f.Close()

	ext := filepath.Ext(fullPath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	response.SetHeader("Content-Type", mimeType)

	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			response.WriteHeader(StatusInternalServerError)
			response.SetHeader("Content-Type", "text/html")
			response.Write([]byte("<h1>500 Internal Server Error</h1>"))
			return
		}
		if n == 0 {
			break
		}
		response.Write(buf[:n])
	}
}

func (s *Server) handleHttpRequest(request Request, response Response) {
	s.Router.ServeHTTP(request, response)
}
