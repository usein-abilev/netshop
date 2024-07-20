package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"netshop/main/api"
	"netshop/main/db"
	"netshop/main/tools"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Cannot read .env file")
	}

	var (
		dbConnectionStr = os.Getenv("DATABASE_URL")
		serverUrl       = tools.TryGetEnv("SERVER_URL", "localhost")
		serverPort      = tools.TryGetEnv("SERVER_PORT", "6900")
	)

	database, dbError := db.NewDatabaseConnection(context.Background(), dbConnectionStr)
	if dbError != nil {
		log.Fatalf("Failed database connection by url '%s'", dbConnectionStr)
	}
	defer database.Close()
	log.Println("Database connected successfully", database.ConnectionString)

	router := api.InitAndCreateRouter(&api.InitEndpointsOptions{
		DatabaseConnection: database,
	})

	httpServerStr := fmt.Sprintf("%s:%s", serverUrl, serverPort)
	if err := http.ListenAndServe(httpServerStr, router); err != nil {
		log.Fatalf("Cannot run server on %s port", serverPort)
	} else {
		log.Printf("Server is running on %s", httpServerStr)
	}
}
