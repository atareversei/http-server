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

// DefaultLogger provides a basic implementation of the Logger interface using the cli package.
type DefaultLogger struct{}

func (dl *DefaultLogger) Success(msg string) {
	cli.Success(msg)
}
func (dl *DefaultLogger) Info(msg string) {
	cli.Info(msg)
}
func (dl *DefaultLogger) Warning(msg string) {
	cli.Warning(msg)
}
func (dl *DefaultLogger) Error(msg string, err error) {
	cli.Error(msg, err)
}

// Logger defines the interface for logging messages at various levels.
type Logger interface {
	Success(msg string)
	Info(msg string)
	Warning(msg string)
	Error(msg string, err error)
}

// NoOpLogger implements Logger but discards all log output.
// Useful for disabling logs without changing logic.
type NoOpLogger struct{}

func (n *NoOpLogger) Success(msg string)          {}
func (n *NoOpLogger) Info(msg string)             {}
func (n *NoOpLogger) Warning(msg string)          {}
func (n *NoOpLogger) Error(msg string, err error) {}

// Server represents an HTTP server with routing, middleware,
// static file serving, and basic logging support.
type Server struct {
	// Host specifies the host address to bind to.
	Host string

	// Port is the TCP port on which the server listens.
	Port int

	// Router handles HTTP route registrations and dispatching.
	Router Router

	// Middlewares is a chain of middleware applied to requests.
	Middlewares []Middleware

	// Static maps URL path prefixes to filesystem directories for file serving.
	Static map[string]string

	// TLSCertFile string
	// TLSKeyFile	 sting

	// loggingEnabled indicates if logging is active.
	loggingEnabled bool

	// previousLogger holds the logger used before disabling logging.
	previousLogger Logger

	// Logger is the active logger for server messages.
	Logger Logger

	// httpServer      *http.Server

	// shutdownTimeout would define how long to wait during graceful shutdown.
	shutdownTimeout time.Duration
}

// DisableLogging disables all logging by replacing the logger with a no-op logger.
func (s *Server) DisableLogging() {
	s.previousLogger = s.Logger
	s.Logger = &NoOpLogger{}
	s.loggingEnabled = false
}

// EnableLogging restores the previously used logger and re-enables logging.
func (s *Server) EnableLogging() {
	s.Logger = s.previousLogger
	s.loggingEnabled = true
}

// Router is the interface for registering and handling HTTP routes.
type Router interface {
	Register(method string, path string, handler Handler)
	Handler
}

// Handler defines the contract for handling an HTTP request and generating a response.
type Handler interface {
	ServeHTTP(req Request, res Response)
}

// DefaultRouter provides a basic implementation of the Router interface.
type DefaultRouter struct {
	routes map[string]map[Method]Handler
	logger Logger
}

// NewRouter creates and returns a new DefaultRouter with initialized route map and logger.
func (s *Server) NewRouter() DefaultRouter {
	return DefaultRouter{
		routes: make(map[string]map[Method]Handler),
		logger: s.Logger,
	}
}

// Register adds a handler and maps it to an HTTP method and a path.
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

// All registers a handler for all HTTP methods on the given path.
func (dr *DefaultRouter) All(path string, handler Handler) {
	dr.checkResourceEntry(path)
	dr.routes[path]["*"] = handler
}

// Get registers a handler for HTTP GET requests on the given path.
func (dr *DefaultRouter) Get(path string, handler Handler) {
	dr.checkResourceEntry(path)
	dr.routes[path]["GET"] = handler
}

// Post registers a handler for HTTP POST requests on the given path.
func (dr *DefaultRouter) Post(path string, handler Handler) {
	dr.checkResourceEntry(path)
	dr.routes[path]["POST"] = handler
}

// checkResourceEntry ensures the inner map for a path exists before assigning a method handler.
func (dr *DefaultRouter) checkResourceEntry(path string) {
	_, ok := dr.routes[path]
	if !ok {
		dr.routes[path] = make(map[Method]Handler)
	}
}

// ServeHTTP handles incoming HTTP requests by dispatching to the appropriate route handler.
func (dr *DefaultRouter) ServeHTTP(req Request, res Response) {
	resource, resOk := dr.routes[req.Path()]
	if !resOk {
		res.WriteHeader(StatusNotFound)
		res.SetHeader("Content-Type", "text/html")
		res.Write([]byte("<h1>404 Not Found</h1>"))
		return
	}
	if req.Method() == MethodOptions {
		handler, handlerOk := resource[MethodGet]
		if !handlerOk {
			res.WriteHeader(StatusNotFound)
			res.SetHeader("Content-Type", "text/html")
			res.Write([]byte("<h1>404 Not Found</h1>"))
			return
		}
		handler.ServeHTTP(req, res)
	}
	handler, handlerOk := resource[req.Method()]
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

// Middleware is a function that wraps a Handler, allowing preprocessing or modification.
type Middleware func(Handler) Handler

// New creates and returns a new Server with the specified port and router.
func New(port int, router Router) Server {
	return Server{
		Port:   port,
		Router: router,
		Logger: &DefaultLogger{},
	}
}

// Start begins listening on the server's configured port and handles incoming TCP connections.
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

// FileHandler maps a URL path prefix to a directory on disk for serving static files.
func (s *Server) FileHandler(pattern string, directory string) {
	s.Static[pattern] = directory
}

// handleConnection reads an HTTP request from a raw TCP connection and dispatches it.
func (s *Server) handleConnection(conn net.Conn) {
	request := newRequestFromTCPConn(conn, s.Logger)
	err := request.Parse()
	response := newResponse(conn, request.Version(), request.Method())
	if err != nil {
		response.WriteHeader(StatusBadRequest)
		response.SetHeader("Content-Type", "text/html")
		response.Write([]byte("<h1>400 Bad Request</h1>"))
		response.Write([]byte(fmt.Sprintf("<p>Digest: %s</p>", err)))
	}
	s.handleRequest(request, response)
	conn.Close()
}

// handleRequest routes a request to either file serving or HTTP handling logic.
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

// handleFileRequest serves files from a static directory mapped by the given prefix.
// It sets proper MIME types and handles 404/500 errors.
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

// handleHttpRequest delegates request handling to the registered Router.
func (s *Server) handleHttpRequest(request Request, response Response) {
	s.Router.ServeHTTP(request, response)
}
