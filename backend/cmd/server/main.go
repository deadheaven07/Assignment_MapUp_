package main

import (
	"log"
	"os"

	"geofencing-alerts/backend/internal/database"
	"geofencing-alerts/backend/internal/handlers"
	"geofencing-alerts/backend/internal/middleware"
	"geofencing-alerts/backend/internal/repositories"
	"geofencing-alerts/backend/internal/services"
	"geofencing-alerts/backend/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	ws := websocket.NewManager()
	go ws.Run()

	repo := repositories.New(db)
	service := services.New(repo, ws)
	handler := handlers.New(service, ws)

	router := gin.Default()
	router.Use(middleware.CORS())
	handler.Register(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
