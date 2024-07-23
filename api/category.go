package api

import (
	"net/http"
	"netshop/main/db"
	"netshop/main/tools"
	"netshop/main/tools/router"
	"strconv"

	"github.com/gorilla/mux"
)

type categoryHandler struct {
	DatabaseConnection *db.DatabaseConnection
	EntityStore        *db.CategoryEntityStore
}

func InitCategoryRouter(parent *router.Router, opts *InitEndpointsOptions) {
	handler := categoryHandler{
		DatabaseConnection: opts.DatabaseConnection,
		EntityStore:        db.NewCategoryEntityStore(opts.DatabaseConnection),
	}
	router := parent.Subrouter()

	router.AddRoute("/categories", handler.handleGet).
		Methods("GET").
		Name("Get all categories").
		Description("Get all categories")

	router.AddRoute("/categories/{id:[0-9]+}", handler.handleGetById).
		Methods("GET").
		Name("Get category entity").
		Description("Get category by given id")
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
