// This package is used to create self-documented API endpoints.
// It independent router package that can be used in any project.
package router

import "net/http"

type (
	RouteOptions struct {
		// Methods is a list of HTTP methods that the route should match
		Methods []string

		// Pattern is the URL pattern of the route
		Pattern string

		// Name is the documentation name of the route
		Name string

		// Description is the documentation description of the route
		Description string

		// Schema is the JSON schema of the request body
		Schema interface{}
	}

	// Route represents a single route in the router
	Route struct {
		Options     RouteOptions
		HandlerFunc http.HandlerFunc
	}

	// Router represents the router
	Router struct {
		Path      string
		Routes    []*Route
		Subroutes []*Router
	}
)

// NewRouter creates a new router
func NewRouter() *Router {
	return &Router{}
}

func (router *Router) Subrouter() *Router {
	newRouter := &Router{}
	router.Subroutes = append(router.Subroutes, newRouter)
	return newRouter
}

func (router *Router) PathPrefix(path string) *Router {
	router.Path = path
	return router
}

// AddRoute adds a new route to the router
func (router *Router) AddRoute(pattern string, handler http.HandlerFunc) *Route {
	route := &Route{
		Options: RouteOptions{
			Methods:     []string{"GET"},
			Name:        pattern,
			Pattern:     pattern,
			Description: "",
			Schema:      nil,
		},
		HandlerFunc: handler,
	}
	router.Routes = append(router.Routes, route)
	return route
}

func (route *Route) Methods(methods ...string) *Route {
	route.Options.Methods = methods
	return route
}

func (route *Route) Name(name string) *Route {
	route.Options.Name = name
	return route
}

func (route *Route) Description(description string) *Route {
	route.Options.Description = description
	return route
}

func (route *Route) Schema(schema interface{}) *Route {
	route.Options.Schema = schema
	return route
}
