package api

import (
	"encoding/json"
	"fmt"
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

type authHTTPError struct {
	Message string
	Code    int
}

func (e *authHTTPError) Error() string {
	return e.Message
}

var (
	errInvalidUserType    = &authHTTPError{"Invalid user type", http.StatusBadRequest}
	errInvalidCredentials = &authHTTPError{"Invalid username or password", http.StatusBadRequest}
	errInternalError      = &authHTTPError{"Internal Server Error", http.StatusInternalServerError}
)

type authHandler struct {
	DatabaseConnection *db.DatabaseConnection
	EmployeeStore      *db.EmployeeEntityStore
	CustomerStore      *db.CustomerEntityStore
}

type commonEntityData struct {
	Id       int64
	Username string
	Password string
}

func InitAuthRouter(parentRouter *router.Router, opts *InitEndpointsOptions) {
	handler := authHandler{
		DatabaseConnection: opts.DatabaseConnection,
	}
	router := parentRouter.Subrouter()

	router.AddRoute("/auth/login", RequireGuest(handler.handleAuth)).
		Methods("POST").
		Name("User Authorization").
		Description("Authorize the user as a customer or employee").
		Schema(map[string]interface{}{
			"type":     "<customer | employee>",
			"username": "<string>",
			"password": "<string>",
		})

	router.AddRoute("/auth/customer/signup", RequireGuest(handler.handleCustomerSignup)).
		Methods("POST").
		Name("Customer Registration").
		Description("Sign up as a customer").
		Schema(map[string]interface{}{
			"first_name": "John",
			"last_name":  "Doe",
			"username":   "john32",
			"password":   "<string>",
			"email":      "example@gmail.com",
			"phone":      "+380000000001",
			"address":    "Lesi Ukrainky Blvd, 26",
			"zipcode":    "01133",
			"city":       "Kyiv",
			"country":    "Ukraine",
		})

	router.AddRoute("/auth/employee/signup", RequireGuest(handler.handleEmployeeSignup)).
		Methods("POST").
		Name("Employee Registration").
		Description("Sign up as an employee").
		Schema(map[string]interface{}{
			"first_name": "Admin",
			"last_name":  "Admin",
			"username":   "admin",
			"password":   "<string>",
			"email":      "admin@netshop.example.com",
			"phone":      "+380000000001",
			"address":    "Lesi Ukrainky Blvd, 26",
			"zipcode":    "01133",
			"city":       "Kyiv",
			"country":    "Ukraine",
		})

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

	var (
		userType = req.Form.Get("type")
		username = req.Form.Get("username")
		password = req.Form.Get("password")
	)

	if username == "" || password == "" {
		tools.RespondWithError(w, errInvalidCredentials.Message, errInvalidCredentials.Code)
		return
	}

	var (
		token    string
		tokenErr error
	)

	if userType == authCustomerTypeStr {
		customerStore := db.NewCustomerEntityStore(handler.DatabaseConnection)
		customer, err := customerStore.GetByUsername(username)
		if err != nil {
			tools.RespondWithError(w, errInvalidCredentials.Message, errInvalidCredentials.Code)
			return
		}
		token, tokenErr = tryGenerateToken(authCustomerTypeStr, password, commonEntityData{
			customer.Id,
			customer.Username,
			customer.Password,
		})
	} else if userType == authEmployeeTypeStr {
		employeeStore := db.NewEmployeeEntityStore(handler.DatabaseConnection)
		employee, err := employeeStore.GetByUsername(username)
		if err != nil {
			tools.RespondWithError(w, errInvalidCredentials.Message, errInvalidCredentials.Code)
			return
		}
		token, tokenErr = tryGenerateToken(authEmployeeTypeStr, password, commonEntityData{
			employee.Id,
			employee.Username,
			employee.Password,
		})
	} else {
		tools.RespondWithError(w, errInvalidUserType.Message, errInvalidUserType.Code)
		return
	}

	if tokenErr != nil {
		if authErr, ok := tokenErr.(*authHTTPError); ok {
			tools.RespondWithError(w, authErr.Message, authErr.Code)
			return
		} else {
			tools.RespondWithError(w, tokenErr.Error(), http.StatusInternalServerError)
			return
		}
	}

	tools.RespondWithSuccess(w, map[string]string{"token": token})
}

func tryGenerateToken(userType string, queryPassword string, data commonEntityData) (string, error) {
	equal, err := tools.ComparePasswordAndHash(queryPassword, data.Password)
	if err != nil || !equal {
		return "", errInvalidCredentials
	}

	token, err := tools.NewJWTToken(map[string]interface{}{
		"type":     userType,
		"id":       data.Id,
		"username": data.Username,
	})
	if err != nil {
		log.Printf("Error creating JWT token: %v", err)
		return "", errInternalError
	}
	return token, nil
}

func (handler *authHandler) handleCustomerSignup(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		tools.RespondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createOpts := &db.CustomerCreateUpdate{}
	if err := json.NewDecoder(req.Body).Decode(createOpts); err != nil {
		tools.RespondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if createOpts.Person.Email == "" {
		tools.RespondWithError(w, "Property 'person.email' is required", http.StatusBadRequest)
		return
	}

	if createOpts.Person.Phone == "" {
		tools.RespondWithError(w, "Property 'person.phone' is required", http.StatusBadRequest)
		return
	}

	if createOpts.Username == "" {
		tools.RespondWithError(w, "Property 'username' is required", http.StatusBadRequest)
		return
	}

	if createOpts.Password == "" {
		tools.RespondWithError(w, "Property 'password' is required", http.StatusBadRequest)
		return
	} else {
		hash, err := tools.HashPassword(createOpts.Password)
		if err != nil {
			tools.RespondWithError(w, "Invalid password", http.StatusBadRequest)
			return
		}
		createOpts.Password = hash
	}

	customerStore := db.NewCustomerEntityStore(handler.DatabaseConnection)
	customer, err := customerStore.Create(req.Context(), createOpts)
	if err != nil {
		tools.RespondWithError(w, fmt.Sprintf("Failed to create customer: %s", err.Error()), http.StatusBadRequest)
		return
	}

	tools.RespondWithSuccess(w, customer)
}

func (handler *authHandler) handleEmployeeSignup(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement employee signup
	tools.RespondWithSuccess(w, "Employee signup is not implemented yet")
}

func (handler *authHandler) handleVerify(w http.ResponseWriter, req *http.Request) {
	tools.RespondWithSuccess(w, "Authorized")
}
