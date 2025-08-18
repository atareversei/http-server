package http

import (
	"strings"
	"sync"
)

// Router is the interface for registering and handling HTTP routes.
type Router interface {
	Register(method string, path string, handler Handler)
	Handler
}

// DefaultRouter provides a basic implementation of the Router interface.
type DefaultRouter struct {
	routes            map[string]map[Method]Handler
	logger            Logger
	cors              CORSConfig
	methodsCache      []Method
	methodsCacheMutex sync.Mutex
}

// NewRouter creates and returns a new DefaultRouter with initialized route map and logger.
func (s *Server) NewRouter() *DefaultRouter {
	router := &DefaultRouter{
		routes: make(map[string]map[Method]Handler),
		logger: s.Logger,
		cors:   s.CorsConfig,
	}
	s.Router = router
	return router
}

// Register adds a handler and maps it to an HTTP method and a path.
func (dr *DefaultRouter) Register(method string, path string, handler Handler) {
	m, err := IsMethodValid(strings.ToUpper(method))
	if err != nil {
		dr.logger.Warning("Unknown method: handler wasn't registered")
		return
	}
	switch m {
	case MethodGet:
		dr.Get(path, handler)
	case MethodPost:
		dr.Post(path, handler)
	case MethodPatch:
		dr.Patch(path, handler)
	case MethodPut:
		dr.Put(path, handler)
	case MethodDelete:
		dr.Delete(path, handler)
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
	dr.routes[path][MethodGet] = handler
}

// Post registers a handler for HTTP POST requests on the given path.
func (dr *DefaultRouter) Post(path string, handler Handler) {
	dr.checkResourceEntry(path)
	dr.routes[path][MethodPost] = handler
}

// Post registers a handler for HTTP Patch requests on the given path.
func (dr *DefaultRouter) Patch(path string, handler Handler) {
	dr.checkResourceEntry(path)
	dr.routes[path][MethodPatch] = handler
}

// Post registers a handler for HTTP Put requests on the given path.
func (dr *DefaultRouter) Put(path string, handler Handler) {
	dr.checkResourceEntry(path)
	dr.routes[path][MethodPut] = handler
}

// Post registers a handler for HTTP Delete requests on the given path.
func (dr *DefaultRouter) Delete(path string, handler Handler) {
	dr.checkResourceEntry(path)
	dr.routes[path][MethodDelete] = handler
}

// ServeHTTP handles incoming HTTP requests by dispatching to the appropriate route handler.
func (dr *DefaultRouter) ServeHTTP(req Request, res Response) {
	resource, resOk := dr.routes[req.Path()]
	if !resOk &&
		!(req.Method() == MethodOptions && req.Path() == "*") &&
		!(req.Method() == MethodConnect) {
		res.WriteHeader(StatusNotFound)
		res.SetHeader("Content-Type", "text/html")
		res.Write([]byte("<h1>404 Not Found</h1>"))
		return
	}

	if req.Method() == MethodHead {
		handleHeadMethod(req, res, resource)
		return
	}

	if req.Method() == MethodOptions {
		handleOptionsMethod(req, res, dr, resource)
		return
	}

	if req.Method() == MethodConnect {
		handleConnectMethod(req, res)
		return
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

// checkResourceEntry ensures the inner map for a path exists before assigning a method handler.
func (dr *DefaultRouter) checkResourceEntry(path string) {
	_, ok := dr.routes[path]
	if !ok {
		dr.routes[path] = make(map[Method]Handler)
	}
}

func (dr *DefaultRouter) getAllAvailableMethods() []Method {
	dr.methodsCacheMutex.Lock()
	defer dr.methodsCacheMutex.Unlock()

	if dr.methodsCache != nil {
		return dr.methodsCache
	}

	var methods = make(map[Method]bool)
	for p := range dr.routes {
		for m := range dr.routes[p] {
			methods[m] = true
		}
	}

	for m := range methods {
		dr.methodsCache = append(dr.methodsCache, m)
	}

	return dr.methodsCache
}
