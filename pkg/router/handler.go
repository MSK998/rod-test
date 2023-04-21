package router

import "net/http"

type MiddlewareHandlerFunc func(http.HandlerFunc) http.HandlerFunc

type HandlerWithMiddlewareFunc func(http.HandlerFunc, ...MiddlewareHandlerFunc) http.HandlerFunc

// Can be used to add a middleware to only a specific endpoint/handler
func HandlerWithMiddleware(handler http.HandlerFunc, middlewares ...MiddlewareHandlerFunc) http.HandlerFunc {
	for _, middlware := range middlewares {
		handler = middlware(handler)
	}
	return handler
}
