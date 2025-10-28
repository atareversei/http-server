package http

import (
	"fmt"
	"net"
	"strings"
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

	// loggingEnabled indicates if logging is active.
	loggingEnabled bool

	// previousLogger holds the logger used before disabling logging.
	previousLogger Logger

	// Logger is the active logger for server messages.
	Logger Logger

	CorsConfig CORSConfig

	// shutdownTimeout would define how long to wait during graceful shutdown.
	shutdownTimeout time.Duration

	// MaxRequestPerConnection
	MaxRequestPerConnection int

	// KeepAliveFor
	KeepAliveFor time.Duration
}

// New creates and returns a new Server with the specified port and router.
func New() Server {
	return Server{
		KeepAliveFor: 30 * time.Second,
		Logger:       &DefaultLogger{},
		CorsConfig: CORSConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []Method{MethodGet, MethodHead, MethodOptions},
			AllowedHeaders: []string{"Content-Type"},
			AllowedMaxAge:  86400,
		},
	}
}

// Start begins listening on the server's configured port and handles incoming TCP connections.
func (s *Server) Start(port int) {
	cli.MadeInBasliqLabs()
	s.Port = port
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
	defer conn.Close()

	requestCount := 0

	for {
		conn.SetDeadline(time.Now().Add(s.KeepAliveFor))
		request := newRequest(conn, s.Logger)
		err := request.Parse()
		response := newResponse(conn, request)
		if err != nil {
			HTTPErrorWithMessage(response, StatusBadRequest, "400 Bad Request", fmt.Sprintf("<p>Digest: %s</p>", err))
			return
		}
		s.handleRequest(request, response)

		if !shouldKeepAlive(request) {
			return
		}

		requestCount++
		if requestCount >= s.MaxRequestPerConnection {
			return
		}
	}
}

// handleHttpRequest delegates request handling to the registered Router.
func (s *Server) handleHttpRequest(request Request, response Response) {
	s.Router.ServeHTTP(request, response)
}

func shouldKeepAlive(req Request) bool {
	connectionHeader, _ := req.Header("Connection")
	connectionHeader = strings.ToLower(strings.TrimSpace(connectionHeader))

	if req.Version() == "HTTP/1.1" {
		return connectionHeader != "close"
	}

	if req.Version() == "HTTP/1.0" {
		return connectionHeader == "keep-alive"
	}

	return false
}
