package http

// Handler defines the contract for handling an HTTP request and generating a response.
type Handler interface {
	ServeHTTP(req Request, res Response)
}

// Middleware is a function that wraps a Handler, allowing preprocessing or modification.
type Middleware func(Handler) Handler
