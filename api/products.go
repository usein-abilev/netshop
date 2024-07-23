package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"netshop/main/db"
	"netshop/main/tools"
	"netshop/main/tools/router"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type productHandler struct {
	DatabaseConnection *db.DatabaseConnection
	EntityStore        *db.ProductEntityStore
}

type getAllQueryParams struct {
	CategoryIds []int64  `schema:"q_category_ids" json:"q_category_ids"`
	SizeIds     []int64  `schema:"q_size_ids" json:"q_size_ids"`
	ColorIds    []int64  `schema:"q_color_ids" json:"q_color_ids"`
	MinPrice    *float64 `schema:"q_min_price" json:"q_min_price"`
	MaxPrice    *float64 `schema:"q_max_price" json:"q_max_price"`
	Limit       int64    `schema:"limit,default:0" json:"limit"`
	Offset      int64    `schema:"offset,default:0" json:"offset"`
	OrderColumn string   `schema:"order_column,default:id" json:"order_column"`
	OrderAsc    bool     `schema:"order_asc,default:false" json:"order_asc"`
}

func InitProductsRouter(router *router.Router, opts *InitEndpointsOptions) {
	handler := productHandler{
		DatabaseConnection: opts.DatabaseConnection,
		EntityStore:        db.NewProductEntityStore(opts.DatabaseConnection),
	}
	productRouter := router.Subrouter()
	productRouter.AddRoute("/products", handler.handleGet).
		Methods("GET").
		Name("Get products").
		Description("Get all products. This endpoint supports filtering by category, size, color, price, and ordering.").
		Schema(getAllQueryParams{
			CategoryIds: []int64{1, 2},
			SizeIds:     []int64{3, 4},
			ColorIds:    []int64{5, 6},
			MinPrice:    nil,
			MaxPrice:    nil,
			Limit:       0,
			Offset:      0,
			OrderColumn: "id",
			OrderAsc:    true,
		})

	productRouter.AddRoute("/products", RequireAuth(handler.handleCreate)).
		Methods("POST").
		Name("Create product").
		Description("Create a new product").
		Schema(&db.ProductCreateUpdate{
			Name:        "Product name",
			Description: "Product description",
			BasePrice:   10.0,
			CategoryId:  1,
			EmployeeId:  1,
			Variants: []db.ProductVariantCreateUpdate{
				{
					SizeId:  1,
					ColorId: 1,
					Price:   10.0,
					Stock:   10,
				},
			},
		})

	productRouter.AddRoute("/products/{id:[0-9]+}", handler.handleGetById).
		Methods("GET").
		Name("Get product by id").
		Description("Gets product by id. The response includes product details and variants")

	productRouter.AddRoute("/products/{id:[0-9]+}", handler.handleEdit).
		Methods("PUT").
		Name("Edit product").
		Description("Edit product by id")

	productRouter.AddRoute("/products/{id:[0-9]+}", handler.handleDelete).
		Methods("DELETE").
		Name("Delete product").
		Description("Deletes product by id")

	productRouter.AddRoute("/products/{id:[0-9]+}/variants", handler.handleGetVariants).
		Methods("GET").
		Name("Get product variants").
		Description("Get product variants by product id")
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
