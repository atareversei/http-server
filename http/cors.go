package http

import (
	"strconv"
	"strings"
)

// CORSConfig specifies the configuration for Cross-Origin Resource Sharing.
// It allows the user to define which external domains are allowed to access
// the server's resources.
type CORSConfig struct {
	// AllowedOrigins is a list of origins that are allowed to make cross-site
	// requests. An origin of "*" allows all origins.
	AllowedOrigins []string

	// AllowedMethods is a list of HTTP methods that are allowed when
	// accessing the resource.
	AllowedMethods []Method

	// AllowedHeaders is a list of HTTP headers that can be used when
	// making the actual request.
	AllowedHeaders []string

	// AllowedCredentials indicates whether the response to the request can be
	// exposed when the credentials flag is true.
	AllowedCredentials bool

	// AllowedMaxAge indicates how long the results of a preflight request
	// (an OPTIONS request) can be cached.
	AllowedMaxAge int
}

// isCORSAllowed checks if a given request has permitted origin, method, and headers.
func (c *CORSConfig) isCORSAllowed(req Request) bool {
	origin, _ := req.Header("Origin")
	isOriginAllowed := c.isOriginAllowed(origin)
	isMethodAllowed := c.isMethodAllowed(req.Header("Access-Control-Request-Method"))
	areHeadersAllowed := c.areHeadersAllowed(req.Header("Access-Control-Request-Headers"))

	if isOriginAllowed && isMethodAllowed && areHeadersAllowed {
		return true
	}
	return false
}

// isOriginAllowed checks if a given origin is permitted by the CORS policy.
func (c *CORSConfig) isOriginAllowed(value string) bool {
	for _, o := range c.AllowedOrigins {
		if o == strings.ToLower(strings.TrimSpace(value)) || o == "*" {
			return true
		}
	}
	return false
}

// isMethodAllowed checks if a given method is permitted by the CORS policy.
func (c *CORSConfig) isMethodAllowed(value string, isValueOk bool) bool {
	if !isValueOk && len(c.AllowedMethods) > 0 {
		return false
	}
	for _, m := range c.AllowedMethods {
		if m.String() == strings.ToUpper(strings.TrimSpace(value)) {
			return true
		}
	}
	return false
}

// areHeadersAllowed checks if all requested headers are permitted by the CORS policy.
func (c *CORSConfig) areHeadersAllowed(value string, isValueOk bool) bool {
	if !isValueOk {
		return true
	}
	headerParts := strings.Split(value, ",")
	for _, h := range headerParts {
		headerIsAllowed := false
		for _, ah := range c.AllowedHeaders {
			if strings.TrimSpace(strings.ToLower(h)) == ah {
				headerIsAllowed = true
			}
		}
		if !headerIsAllowed {
			return false
		}
	}
	return true
}

// getAllowedOriginsHeader returns the value for the Access-Control-Allow-Origin header.
func (c *CORSConfig) getAllowedOriginsHeader() string {
	return strings.Join(c.AllowedOrigins, ", ")
}

// getAllowedMethodsHeader returns the value for the Access-Control-Allow-Methods header.
func (c *CORSConfig) getAllowedMethodsHeader() string {
	corsMtdStrArr := make([]string, len(c.AllowedMethods))
	for i, m := range c.AllowedMethods {
		corsMtdStrArr[i] = m.String()
	}
	return strings.Join(corsMtdStrArr, ", ")
}

// getAllowedHeadersHeader returns the value for the Access-Control-Allow-Headers header.
func (c *CORSConfig) getAllowedHeadersHeader() string {
	return strings.Join(c.AllowedHeaders, ", ")
}

// getAllowedCredentialsHeader returns the value for the Access-Control-Allow-Credentials header.
func (c *CORSConfig) getAllowedCredentialsHeader() string {
	return strconv.FormatBool(c.AllowedCredentials)
}

// getAllowedMaxAgeHeader returns the value for the Access-Control-Max-Age header.
func (c *CORSConfig) getAllowedMaxAgeHeader() string {
	return strconv.Itoa(c.AllowedMaxAge)
}
