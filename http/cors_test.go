package http

import (
	"testing"
)

// TestIsOriginAllowed validates the origin checking logic.
func TestIsOriginAllowed(t *testing.T) {
	testCases := []struct {
		name     string
		config   CORSConfig
		origin   string
		expected bool
	}{
		{"Wildcard allows any origin", CORSConfig{AllowedOrigins: []string{"*"}}, "https://example.com", true},
		{"Exact match is allowed", CORSConfig{AllowedOrigins: []string{"https://good.com"}}, "https://good.com", true},
		{"Case is ignored for incoming origin", CORSConfig{AllowedOrigins: []string{"https://good.com"}}, "https://GOOD.com", true},
		{"Non-listed origin is not allowed", CORSConfig{AllowedOrigins: []string{"https://good.com"}}, "https://evil.com", false},
		{"Empty origin is not allowed", CORSConfig{AllowedOrigins: []string{"https://good.com"}}, "", false},
		{"Empty config allows nothing", CORSConfig{}, "https://any.com", false},
		{"Multiple entries accept matches", CORSConfig{AllowedOrigins: []string{"https://good.com", "https://any.net"}}, "https://any.net", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.config.isOriginAllowed(tc.origin); got != tc.expected {
				t.Errorf("for origin '%s', expected %v but got %v", tc.origin, tc.expected, got)
			}
		})
	}
}

// TestIsMethodAllowed validates the method checking logic for preflight requests.
func TestIsMethodAllowed(t *testing.T) {
	config := CORSConfig{AllowedMethods: []Method{MethodGet, MethodPost}}
	testCases := []struct {
		name     string
		method   string
		ok       bool
		expected bool
	}{
		{"Allowed method passes", "POST", true, true},
		{"Allowed method (case difference) passes", "get", true, true},
		{"Disallowed method fails", "DELETE", true, false},
		{"Disallowed empty method", "", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// In a real scenario, the method from the header would be normalized to uppercase
			if got := config.isMethodAllowed(tc.method, tc.ok); got != tc.expected {
				t.Errorf("for method '%s', expected %v but got %v", tc.method, tc.expected, got)
			}
		})
	}
}

// TestAreHeadersAllowed validates the header checking logic for preflight requests.
func TestAreHeadersAllowed(t *testing.T) {
	config := CORSConfig{AllowedHeaders: []string{"content-type", "authorization"}}
	testCases := []struct {
		name     string
		headers  string
		ok       bool
		expected bool
	}{
		{"All requested headers are allowed", "Content-Type, Authorization", true, true},
		{"One requested header is allowed", "authorization", true, true},
		{"Case is ignored for incoming headers", "AUTHORIZATION", true, true},
		{"A non-allowed header fails the check", "Content-Type, X-Custom-Header", true, false},
		{"Empty header string (no custom headers) passes", "", false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := config.areHeadersAllowed(tc.headers, tc.ok); got != tc.expected {
				t.Errorf("for headers '%s', expected %v but got %v", tc.headers, tc.expected, got)
			}
		})
	}
}
