package http

type HandlerFunc func(req Request, res *Response)

// All is a catch-all handler registrar.
func All(pattern string, handler HandlerFunc) {
	checkResourceEntry(pattern)
	router[pattern]["ALL"] = handler
}

// Get is a GET method handler registrar.
func Get(pattern string, handler HandlerFunc) {
	checkResourceEntry(pattern)
	router[pattern]["GET"] = handler
}

// Post is a POST method handler registrar.
func Post(pattern string, handler HandlerFunc) {
	checkResourceEntry(pattern)
	router[pattern]["POST"] = handler
}

// checkResourceEntry is used to initialize the inner map of a router
// if it has not yet been initialized.
func checkResourceEntry(pattern string) {
	_, ok := s.router[pattern]
	if !ok {
		s.router[pattern] = make(map[string]HandlerFunc)
	}
}
