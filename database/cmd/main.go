package main

import (
	"context"
	"log"
	"net/http"

	"github.com/maksroxx/DeliveryService/database/internal/configs"
	"github.com/maksroxx/DeliveryService/database/internal/handlers"
	"github.com/maksroxx/DeliveryService/database/internal/middleware"
	"github.com/maksroxx/DeliveryService/database/internal/repository"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var (
		cfg           = configs.Load()
		mongoCfg      = cfg.Database.MongoDB
		clientOptions = options.Client().ApplyURI(mongoCfg.URI)
		client, err   = mongo.Connect(context.Background(), clientOptions)
	)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	var (
		db             = client.Database(mongoCfg.Database)
		repo           = repository.NewMongoRepository(db, "packages")
		mux            = http.NewServeMux()
		protected      = http.NewServeMux()
		packageHandler = handlers.NewPackageHandler(repo)
	)
	packageHandler.RegisterDefaultRoutes(mux)
	packageHandler.RegisterUserRoutes(protected)
	protectedHandler := middleware.AuthMiddleware(protected)

	mux.Handle("/", protectedHandler)
	mux.Handle("/metrics", promhttp.Handler())

	log.Printf("Server starting on %s", cfg.Server.Address)

	if err := http.ListenAndServe(cfg.Server.Address, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
