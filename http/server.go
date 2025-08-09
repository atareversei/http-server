package http

import (
	"fmt"
	"net"
	"time"

	"github.com/atareversei/http-server/internal/cli"
)

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

	CorsConfig CORSConfig

	// shutdownTimeout would define how long to wait during graceful shutdown.
	shutdownTimeout time.Duration
}

// New creates and returns a new Server with the specified port and router.
func New(port int, router Router) Server {
	return Server{
		Port:   port,
		Router: router,
		Logger: &DefaultLogger{},
		CorsConfig: CORSConfig{
			AllowedMethods: []Method{MethodGet, MethodHead, MethodOptions},
			AllowedHeaders: []string{"Content-Type"},
		},
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

	if req.Method() == MethodHead {
		handler, handlerOk := resource[MethodGet]
		if !handlerOk {
			HTTPError(res, StatusNotFound)
			return
		}
		handler.ServeHTTP(req, res)
	}

	if req.Method() == MethodOptions {
	}

	handler, handlerOk := resource[req.Method()]
	if !handlerOk {
		catchAll, allOk := resource["*"]
		if allOk {
			catchAll.ServeHTTP(req, res)
			return
		}
		HTTPError(res, StatusMethodNotAllowed)
		return
	}
	handler.ServeHTTP(req, res)
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

// handleHttpRequest delegates request handling to the registered Router.
func (s *Server) handleHttpRequest(request Request, response Response) {
	s.Router.ServeHTTP(request, response)
}
