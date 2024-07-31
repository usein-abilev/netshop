package api

import (
	"log"
	"net/http"
	"netshop/main/db"
	"netshop/main/tools"
	"netshop/main/tools/router"
	"path"
	"strings"

	"github.com/gorilla/mux"
	cors "github.com/rs/cors"
)

type InitEndpointsOptions struct {
	DatabaseConnection *db.DatabaseConnection
}

func moveRouterToMux(router *router.Router, muxRouter *mux.Router) {
	for _, route := range router.Routes {
		muxRouter.HandleFunc(path.Join(router.Path, route.Options.Pattern), route.HandlerFunc).Methods(route.Options.Methods...)
	}

	for _, subrouter := range router.Subroutes {
		subMuxRouter := muxRouter.PathPrefix(path.Join(router.Path, subrouter.Path)).Subrouter()
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

func getRoutersSchema(router *router.Router, parentPath string) []routerEndpointInfo {
	if parentPath == "" {
		parentPath = router.Path
	}

	routes := []routerEndpointInfo{}
	for _, route := range router.Routes {
		routes = append(routes, routerEndpointInfo{
			Methods:     route.Options.Methods,
			Name:        route.Options.Name,
			Description: route.Options.Description,
			Path:        path.Join(parentPath, route.Options.Pattern),
			Schema:      route.Options.Schema,
		})
	}

	for _, subrouter := range router.Subroutes {
		subrouterSchemaList := getRoutersSchema(subrouter, path.Join(parentPath, subrouter.Path))
		routes = append(routes, subrouterSchemaList...)
	}

	return routes
}

func InitAndCreateRouter(opts *InitEndpointsOptions) http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.StrictSlash(true)

	// Init file server for ./static/files/*.webp files
	fileServer := http.FileServer(http.Dir("./static/files"))
	muxRouter.HandleFunc("/static/files/{filename}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename := vars["filename"]
		if !strings.HasSuffix(filename, ".webp") {
			tools.RespondWithError(w, "Invalid file format", http.StatusBadRequest)
			return
		}
		http.StripPrefix("/static/files/", fileServer).ServeHTTP(w, r)
	})

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	router := router.NewRouter()
	router.PathPrefix("/api/v1")

	InitAuthRouter(router, opts)
	InitCategoryRouter(router, opts)
	InitProductsRouter(router, opts)
	InitFileRouter(router, opts)
	InitOrderRouter(router, opts)

	// move all registered routes to the mux router to be able to use it
	moveRouterToMux(router, muxRouter)

	// get list of all routes in the router
	routerList := getRoutersSchema(router, router.Path)
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

	handler := corsConfig.Handler(muxRouter)
	return handler
}
