package api

import (
	"log"
	"net/http"
	"netshop/main/db"
	"netshop/main/tools"
	"netshop/main/tools/router"
	"strings"

	"github.com/gorilla/mux"
)

type InitEndpointsOptions struct {
	DatabaseConnection *db.DatabaseConnection
}

func moveRouterToMux(router *router.Router, muxRouter *mux.Router) {
	for _, route := range router.Routes {
		muxRouter.HandleFunc(route.Options.Pattern, route.HandlerFunc).Methods(route.Options.Methods...)
	}

	for _, subrouter := range router.Subroutes {
		subMuxRouter := muxRouter.PathPrefix(subrouter.Path).Subrouter()
		moveRouterToMux(subrouter, subMuxRouter)
	}
}

type routerEndpointInfo struct {
	Methods     []string    `json:"methods"`
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	Description string      `json:"description"`
	Schema      interface{} `json:"schema"`
}

func getRoutersSchema(router *router.Router) []routerEndpointInfo {
	routes := []routerEndpointInfo{}
	for _, route := range router.Routes {
		routes = append(routes, routerEndpointInfo{
			Methods:     route.Options.Methods,
			Name:        route.Options.Name,
			Description: route.Options.Description,
			Path:        route.Options.Pattern,
			Schema:      route.Options.Schema,
		})
	}

	for _, subrouter := range router.Subroutes {
		subrouterSchemaList := getRoutersSchema(subrouter)
		routes = append(routes, subrouterSchemaList...)
	}

	return routes
}

func InitAndCreateRouter(opts *InitEndpointsOptions) http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.StrictSlash(true)

	muxRouter.Use(LoggingMiddleware)

	router := router.NewRouter()
	InitAuthRouter(router, opts)
	InitCategoryRouter(router, opts)
	InitProductsRouter(router, opts)

	// move all registered routes to the mux router to be able to use it
	moveRouterToMux(router, muxRouter)

	// get list of all routes in the router
	routerList := getRoutersSchema(router)
	for _, route := range routerList {
		log.Printf("Initialized route: [%s] %s", strings.Join(route.Methods, ", "), route.Path)
	}
	muxRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tools.RespondWithSuccess(w, routerList)
	})

	muxRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tools.RespondWithError(w, "Not found", http.StatusNotFound)
	})
	muxRouter.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tools.RespondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	return muxRouter
}
