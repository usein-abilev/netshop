package api

import (
	"log"
	"net/http"
	"netshop/main/db"
	"netshop/main/tools"
	"netshop/main/tools/router"
)

const (
	authCustomerTypeStr = "customer"
	authEmployeeTypeStr = "employee"
)

type authHandler struct {
	DatabaseConnection *db.DatabaseConnection
}

func InitAuthRouter(parentRouter *router.Router, opts *InitEndpointsOptions) {
	handler := authHandler{
		DatabaseConnection: opts.DatabaseConnection,
	}
	router := parentRouter.Subrouter()

	router.AddRoute("/auth", handler.handleAuth).
		Methods("POST").
		Name("User Authorization").
		Description("Authorize the user as a customer or employee")

	router.AddRoute("/auth/verify", RequireAuth(handler.handleVerify)).
		Methods("POST").
		Name("Verify Authorization").
		Description("Verify the user's authorization. Returns 200 if authorized, 401 if not")
}

// Authorizes the customer or employee
func (handler *authHandler) handleAuth(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		tools.RespondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Header.Get("Authorization") != "" {
		tools.RespondWithError(w, "User already authorized", http.StatusBadRequest)
		return
	}

	var (
		userType = req.Form.Get("type")
		username = req.Form.Get("username")
		password = req.Form.Get("password")
	)

	if username == "" || password == "" {
		tools.RespondWithError(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	if userType == authCustomerTypeStr {
		// TODO: Add verification and authorization for customers
	} else if userType == authEmployeeTypeStr {
		employeeStore := db.NewEmployeeEntityStore(handler.DatabaseConnection)
		employee, err := employeeStore.GetEmployeeByUsername(username)
		if err != nil {
			tools.RespondWithError(w, "Invalid username or password", http.StatusBadRequest)
			return
		}
		equal, err := tools.ComparePasswordAndHash(password, employee.Password)
		if err != nil {
			log.Println("Error comparing password and hash:", err)
		}

		if err != nil || !equal {
			tools.RespondWithError(w, "Invalid username or password", http.StatusBadRequest)
			return
		}

		token, err := tools.NewJWTToken(map[string]interface{}{
			"type":     userType,
			"id":       employee.Id,
			"username": employee.Username,
		})
		if err != nil {
			log.Printf("Error creating JWT token: %v", err)
			tools.RespondWithError(w, "Internal error", http.StatusInternalServerError)
			return
		}
		tools.RespondWithSuccess(w, map[string]interface{}{"token": token})
	} else {
		tools.RespondWithError(w, "Invalid user type", http.StatusBadRequest)
		return
	}
}

func (handler *authHandler) handleVerify(w http.ResponseWriter, req *http.Request) {
	tools.RespondWithSuccess(w, "Authorized")
}
