package http

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []Method
	AllowedHeaders   []string
	AllowCredentials bool
}
