package middleware

import "net/http"

type Middleware func(http.HandlerFunc) http.HandlerFunc

// chains multiple middleware functions
func ChainMiddleware(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
