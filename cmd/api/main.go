package main

import (
	"log"

	"github.com/sahilkarwasra/moepay/internal/config"
	"github.com/sahilkarwasra/moepay/internal/database"
	"github.com/sahilkarwasra/moepay/internal/repository/mongo"
	"github.com/sahilkarwasra/moepay/internal/routes"
	"github.com/sahilkarwasra/moepay/internal/service"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		log.Fatal("failed to load config", err)
	}

	client, db, err := database.Connect(cfg)

	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	defer func() {
		if err := database.Disconnect(client); err != nil {
			log.Printf("failed to disconnect from database %v", err)
		}
	}()

	userCollection := db.Collection(config.UsersCollections)
	otpCollection := db.Collection(config.OtpCollections)
	// tokenCollection := database.Collection(config.TokenCollection)

	userRepo := mongo.NewUserMongoRepository(userCollection)
	otpRepo := mongo.NewOtpMongoRepository(otpCollection)

	userService := service.NewUserService(userRepo, otpRepo)

	routes.SetupRouter(userService)

	log.Println("..--..")

	router := routes.SetupRouter(userService)

	port := cfg.ServerPort

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
