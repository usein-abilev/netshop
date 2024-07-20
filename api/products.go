package api

import (
	"fmt"
	"net/http"
	"netshop/main/db"
	"netshop/main/tools"
	"strconv"

	"github.com/gorilla/mux"
)

type productHandler struct {
	DatabaseConnection *db.DatabaseConnection
	EntityStore        *db.ProductEntityStore
}

func InitProductsRouter(router *mux.Router, opts *InitEndpointsOptions) {
	handler := productHandler{
		DatabaseConnection: opts.DatabaseConnection,
		EntityStore:        db.NewProductEntityStore(opts.DatabaseConnection),
	}
	productRouter := router.NewRoute().Subrouter()
	productRouter.HandleFunc("/products", handler.handleGet).Methods("GET")
	productRouter.HandleFunc("/products", handler.handleCreate).Methods("POST")
	productRouter.HandleFunc("/products/{id:[0-9]+}", handler.handleGetById).Methods("GET")
	productRouter.HandleFunc("/products/{id:[0-9]+}", handler.handleUpdate).Methods("PUT")
	productRouter.HandleFunc("/products/{id:[0-9]+}", handler.handleDelete).Methods("DELETE")
}

func (ph *productHandler) handleGet(w http.ResponseWriter, req *http.Request) {
	products, err := ph.EntityStore.GetProducts()
	if err != nil {
		tools.RespondWithError(w, "Unexpected error while received products", http.StatusInternalServerError)
		return
	}

	tools.RespondWithSuccess(w, products)
}

func (ph *productHandler) handleGetById(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		tools.RespondWithError(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("Handle products.getById query: %v\n", id)

	product, err := ph.EntityStore.GetProductById(id)
	if err != nil {
		tools.RespondWithError(w, "Product not found", http.StatusNotFound)
		return
	}

	tools.RespondWithSuccess(w, product)
}

func (ph *productHandler) handleCreate(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement product creation
	tools.RespondWithError(w, "products.create not implemented yet", http.StatusNotImplemented)
}

func (ph *productHandler) handleUpdate(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement product update
	tools.RespondWithError(w, "products.update not implemented yet", http.StatusNotImplemented)
}

func (ph *productHandler) handleDelete(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement product deletion
	tools.RespondWithError(w, "products.delete not implemented yet", http.StatusNotImplemented)
}
