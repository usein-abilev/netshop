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
	"time"

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

	server := &http.Server{
		Addr:         httpServerStr,
		IdleTimeout:  time.Second * 15,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		Handler:      router,
	}

	log.Printf("Starting server on http://%s", httpServerStr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Cannot run server on %s port %s", serverPort, err.Error())
	}
}
