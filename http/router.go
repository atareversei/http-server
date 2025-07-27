package http

// Router holds the routing information.
// The structure can be simplified as -> [PATH][METHOD]handler
type Router struct {
	mapper map[string]map[string]HandlerFunc
}
