package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sahilkarwasra/moepay/config"
	"github.com/sahilkarwasra/moepay/routes"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoUri := os.Getenv("mongoUri")
	if mongoUri == "" {
		mongoUri = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatalf("Error creating the Client %s", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to Ping MongoDb", err)
	}

	log.Println("Connected To MongoDB")

	dbName := os.Getenv("dbName")
	if dbName == "" {
		dbName = "moepay"
	}

	config.MongoDB = client.Database(dbName)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	routes.RegisterRoutes(r)

	port := os.Getenv("port")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start the server", err)
	}

}
