package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"netshop/main/db"
	"netshop/main/tools"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type productHandler struct {
	DatabaseConnection *db.DatabaseConnection
	EntityStore        *db.ProductEntityStore
}

type getAllQueryParams struct {
	CategoryIds []int64  `schema:"q_category_ids"`
	SizeIds     []int64  `schema:"q_size_ids"`
	ColorIds    []int64  `schema:"q_color_ids"`
	MinPrice    *float64 `schema:"q_min_price"`
	MaxPrice    *float64 `schema:"q_max_price"`
	Limit       int64    `schema:"limit"`
	Offset      int64    `schema:"offset"`
	OrderColumn string   `schema:"order_column"`
	OrderAsc    bool     `schema:"order_asc"`
}

func InitProductsRouter(router *mux.Router, opts *InitEndpointsOptions) {
	handler := productHandler{
		DatabaseConnection: opts.DatabaseConnection,
		EntityStore:        db.NewProductEntityStore(opts.DatabaseConnection),
	}
	productRouter := router.NewRoute().Subrouter()
	productRouter.Schemes()

	productRouter.HandleFunc("/products", handler.handleGet).Methods("GET")
	productRouter.HandleFunc("/products", RequireAuth(handler.handleCreate)).Methods("POST")
	productRouter.HandleFunc("/products/{id:[0-9]+}", handler.handleGetById).Methods("GET")
	productRouter.HandleFunc("/products/{id:[0-9]+}", handler.handleEdit).Methods("PUT")
	productRouter.HandleFunc("/products/{id:[0-9]+}", handler.handleDelete).Methods("DELETE")
	productRouter.HandleFunc("/products/{id:[0-9]+}/variants", handler.handleGetVariants).Methods("GET")
}

func (ph *productHandler) handleGet(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	queryParams := getAllQueryParams{}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(&queryParams, query); err != nil {
		tools.RespondWithError(w, "Invalid query params", http.StatusBadRequest)
		return
	}

	products, err := ph.EntityStore.GetEntities(&db.ProductGetEntitiesOptions{
		Query: &db.ProductGetEntitiesQueryOpts{
			CategoryIds: queryParams.CategoryIds,
			SizeIds:     queryParams.SizeIds,
			ColorIds:    queryParams.ColorIds,
			MinPrice:    queryParams.MinPrice,
			MaxPrice:    queryParams.MaxPrice,
		},
		Limit:       queryParams.Limit,
		Offset:      queryParams.Offset,
		OrderColumn: queryParams.OrderColumn,
		OrderAsc:    queryParams.OrderAsc,
	})
	if err != nil {
		tools.RespondWithError(w, fmt.Sprintf("Unexpected error while received products: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	tools.RespondWithSuccess(w, products)
}

func (ph *productHandler) handleGetById(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		tools.RespondWithError(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	product, err := ph.EntityStore.GetById(id)
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

	if err := ph.EntityStore.Create(context.Background(), createOpts); err != nil {
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

func (ph *productHandler) handleGetVariants(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	if err != nil {
		tools.RespondWithError(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	variants, err := ph.EntityStore.GetVariants(id)
	if err != nil {
		log.Printf("Error while getting product variants: %s", err.Error())
		tools.RespondWithError(w, "Product not found", http.StatusNotFound)
		return
	}

	tools.RespondWithSuccess(w, variants)
}
