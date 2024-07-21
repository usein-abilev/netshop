package api

import (
	"context"
	"encoding/json"
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
	productRouter.HandleFunc("/products/create", RequireAuth(handler.handleCreate)).Methods("POST")
	productRouter.HandleFunc("/products/{id:[0-9]+}", handler.handleGetById).Methods("GET")
	productRouter.HandleFunc("/products/{id:[0-9]+}/edit", handler.handleEdit).Methods("PUT")
	productRouter.HandleFunc("/products/{id:[0-9]+}/delete", handler.handleDelete).Methods("DELETE")
}

func (ph *productHandler) handleGet(w http.ResponseWriter, req *http.Request) {
	products, err := ph.EntityStore.GetProducts()
	if err != nil {
		tools.RespondWithError(w, fmt.Sprintf("Unexpected error while received products: %s", err.Error()), http.StatusInternalServerError)
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

	product, err := ph.EntityStore.GetProductById(id)
	if err != nil {
		tools.RespondWithError(w, "Product not found", http.StatusNotFound)
		return
	}

	tools.RespondWithSuccess(w, product)
}

func (ph *productHandler) handleCreate(w http.ResponseWriter, req *http.Request) {
	user := req.Context().Value("user").(*tools.UserTokenClaims)

	createOpts := &db.ProductCreateUpdate{}
	if err := json.NewDecoder(req.Body).Decode(createOpts); err != nil {
		tools.RespondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createOpts.EmployeeId = user.Id

	if err := ph.EntityStore.CreateProduct(context.Background(), createOpts); err != nil {
		tools.RespondWithError(w, err.Error(), http.StatusBadRequest)
		return
	}

	tools.RespondWithSuccess(w, true)
}

func (ph *productHandler) handleEdit(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement product update
	tools.RespondWithError(w, "products.update not implemented yet", http.StatusNotImplemented)
}

func (ph *productHandler) handleDelete(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement product deletion
	tools.RespondWithError(w, "products.delete not implemented yet", http.StatusNotImplemented)
}
