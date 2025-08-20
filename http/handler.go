package http

// Handler defines the contract for handling an HTTP request and generating a response.
type Handler interface {
	ServeHTTP(req Request, res Response)
}

type HandlerFunc func(req Request, res Response)

func (f HandlerFunc) ServeHTTP(req Request, res Response) {
	f(req, res)
}

// Middleware is a function that wraps a Handler, allowing preprocessing or modification.
type Middleware func(Handler) Handler
