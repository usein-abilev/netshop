package api

import (
	"fmt"
	"net/http"
	"netshop/main/db"
	"netshop/main/tools"
	"netshop/main/tools/router"
)

type orderHandler struct {
	DatabaseConnection *db.DatabaseConnection
	EntityStore        *db.OrderEntityStore
}

func InitOrderRouter(parent *router.Router, opts *InitEndpointsOptions) {
	handler := orderHandler{
		DatabaseConnection: opts.DatabaseConnection,
		EntityStore:        db.NewOrderEntity(opts.DatabaseConnection),
	}

	router := parent.Subrouter()

	router.AddRoute("/orders", RequireAuth(handler.handleGet)).
		Methods("GET").
		Name("Get user's orders").
		Description("Get all orders of the current user")
}

func (handler *orderHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	customerId := r.Context().Value("user").(*tools.UserTokenClaims).Id
	items, err := handler.EntityStore.GetAll(&db.OrderGetAllOptions{CustomerId: &customerId})
	if err != nil {
		tools.RespondWithError(w, fmt.Sprintf("Cannot get orders information: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	tools.RespondWithSuccess(w, items)
}
