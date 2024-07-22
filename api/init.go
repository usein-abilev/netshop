package api

import (
	"log"
	"net/http"
	"netshop/main/db"
	"netshop/main/tools"
	"strings"

	"github.com/gorilla/mux"
)

type InitEndpointsOptions struct {
	DatabaseConnection *db.DatabaseConnection
}

type routerEndpointInfo struct {
	Name    string   `json:"name"`
	Methods []string `json:"methods"`
}

func InitAndCreateRouter(opts *InitEndpointsOptions) *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.Use(LoggingMiddleware)

	InitAuthRouter(router, opts)
	InitCategoryRouter(router, opts)
	InitProductsRouter(router, opts)

	// get list of all routes in the router
	routerList := []routerEndpointInfo{}
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		routerName, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()

		if len(routerName) == 0 {
			return nil
		}

		routerList = append(routerList, routerEndpointInfo{
			Name:    routerName,
			Methods: methods,
		})
		log.Printf("Route: [%s] %s", strings.Join(methods, ", "), routerName)
		return nil
	})

	if err != nil {
		log.Printf("Error getting routes: %v", err)
	}
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tools.RespondWithSuccess(w, routerList)
	})

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tools.RespondWithError(w, "Not found", http.StatusNotFound)
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tools.RespondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	return router
}
