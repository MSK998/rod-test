package router

import (
	"net/http"
)

type RegisterFunc func(*http.ServeMux)

type RegistrationList []RegisterFunc

// Global variable that can be used to register new routes to the router
var RegisteredHandlers RegistrationList

// Called by handlers to register routes - keeping all of the handlers
// functionality in one place
func (r *RegistrationList) Register(hf RegisterFunc) {
	RegisteredHandlers = append(RegisteredHandlers, hf)
}

// Applies the handlers to the passed in router
func RegisterRoutes(router *http.ServeMux) {
	for _, v := range RegisteredHandlers {
		v(router)
	}
}

// Create a new router and attach passed in middlewares globally
func New(globalMiddlewares ...func(http.Handler) http.Handler) http.Handler {
	baseRouter := http.NewServeMux()
	RegisterRoutes(baseRouter)
	router := registerMiddleware(baseRouter, globalMiddlewares...)
	return router
}

func registerMiddleware(r http.Handler, gm ...func(http.Handler) http.Handler) http.Handler {
	if len(gm) == 0 {
		return r
	}

	for _, m := range gm {
		r = m(r)
	}

	return r
}
