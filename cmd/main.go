package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-mgmt/config"
	"user-mgmt/database"
	"user-mgmt/logger"
	"user-mgmt/server"

	"github.com/gin-gonic/gin"
)

func main() {

	// Init logger
	log := logger.New()

	// Load config
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal().Msgf("failed to load config: %v", err)
	}

	// Create context with cancel for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init MongoDB
	mongoClient, err := database.NewMongoDBClient(cfg.MongoDBConfig)
	if err != nil {
		log.Fatal().Msgf("failed to init MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	db := database.GetDatabase(mongoClient, cfg.MongoDBConfig)
	if err := database.EnsureIndexes(ctx, db); err != nil {
		log.Fatal().Msgf("failed to ensure database indexes: %v", err)
	}

	log.Info().Msgf("Successfully connected to MongoDB database: %s", cfg.MongoDBConfig.DatabaseName)

	gin.SetMode(cfg.Server.GinMode)
	// Init router
	srv := server.NewServer(cfg, db, log)
	router := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server
	go func() {
		log.Info().Msgf("Starting server on port: %s", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("failed to start server: %v", err)
		}
	}()

	// Wait for signal to stop server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Gracefully shutdown server
	log.Printf("Shutting down server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("failed to shutdown server: %v", err)
	}
	log.Info().Msg("Server shutdown successfully")

}
