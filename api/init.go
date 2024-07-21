package api

import (
	"net/http"
	"netshop/main/db"
	"netshop/main/tools"

	"github.com/gorilla/mux"
)

type InitEndpointsOptions struct {
	DatabaseConnection *db.DatabaseConnection
}

func InitAndCreateRouter(opts *InitEndpointsOptions) *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.Use(LoggingMiddleware)

	InitAuthRouter(router, opts)
	InitCategoryRouter(router, opts)
	InitProductsRouter(router, opts)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tools.RespondWithError(w, "Not found", http.StatusNotFound)
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tools.RespondWithError(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	return router
}
