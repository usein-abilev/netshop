package api

import (
	"netshop/main/db"

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
	InitProductsRouter(router, opts)

	return router
}
