package api

import (
	"net/http"
	"netshop/main/db"
	"netshop/main/tools"
	"strconv"

	"github.com/gorilla/mux"
)

type categoryHandler struct {
	DatabaseConnection *db.DatabaseConnection
	EntityStore        *db.CategoryEntityStore
}

func InitCategoryRouter(router *mux.Router, opts *InitEndpointsOptions) {
	handler := categoryHandler{
		DatabaseConnection: opts.DatabaseConnection,
		EntityStore:        db.NewCategoryEntityStore(opts.DatabaseConnection),
	}
	productRouter := router.NewRoute().Subrouter()
	productRouter.HandleFunc("/categories", handler.handleGet).Methods("GET")
	productRouter.HandleFunc("/categories/{id:[0-9]+}", handler.handleGetById).Methods("GET")
}

func (c *categoryHandler) handleGet(w http.ResponseWriter, req *http.Request) {
	items, err := c.EntityStore.GetCategories()
	if err != nil {
		tools.RespondWithError(w, "Unexpected error while received categories", http.StatusInternalServerError)
		return
	}

	tools.RespondWithSuccess(w, items)
}

func (c *categoryHandler) handleGetById(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		tools.RespondWithError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	product, err := c.EntityStore.GetCategoryById(id)
	if err != nil {
		tools.RespondWithError(w, "Category not found", http.StatusNotFound)
		return
	}

	tools.RespondWithSuccess(w, product)
}
