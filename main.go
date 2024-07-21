package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"netshop/main/api"
	"netshop/main/config"
	"netshop/main/db"
)

func main() {
	database, err := db.NewDatabaseConnection(context.Background(), &db.DatabaseConnectionOptions{
		ConnectionURL:  config.AppConfig.DatabaseURL,
		MaxConns:       10,
		MinConns:       1,
		ConnectTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed database connection by url '%s'", config.AppConfig.DatabaseURL)
	}
	defer database.Close()
	log.Println("Database connection pool initialized successfully", database.ConnectionString)

	if err = database.Connection.Ping(database.Context); err != nil {
		log.Fatalf("Failed database ping by url '%s'", config.AppConfig.DatabaseURL)
	}

	router := api.InitAndCreateRouter(&api.InitEndpointsOptions{
		DatabaseConnection: database,
	})

	httpServerStr := fmt.Sprintf("%s:%s", config.AppConfig.ServerURL, config.AppConfig.ServerPort)

	server := &http.Server{
		Addr:         httpServerStr,
		IdleTimeout:  time.Second * 15,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		Handler:      router,
	}

	log.Printf("Starting server on http://%s", httpServerStr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Cannot run server on %s. Error: %s", httpServerStr, err.Error())
	}
}
