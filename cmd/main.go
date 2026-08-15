package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-mgmt/config"
	"user-mgmt/database"
	"user-mgmt/server"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func main() {

	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init MongoDB
	mongoClient, err := database.NewMongoDBClient(cfg.MongoDBConfig)
	if err != nil {
		log.Fatalf("failed to init MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	db := database.GetDatabase(mongoClient, cfg.MongoDBConfig)

	log.Printf("Successfully connected to MongoDB database: %s", cfg.MongoDBConfig.DatabaseName)

	gin.SetMode(cfg.Server.GinMode)
	// Init router
	srv := server.NewServer(cfg, db)
	router := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Default().Printf("Starting server on port: %s", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for signal to stop server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Gracefully shutdown server
	log.Printf("Shutting down server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}
	log.Printf("Server shutdown successfully")

}
