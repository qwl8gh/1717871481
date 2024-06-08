package main

import (
	"context"
	"log"
	"net/http"
	"time"
	_ "web-messaging-service/docs"
	"web-messaging-service/pkg/api/routes"
	"web-messaging-service/pkg/config"
	"web-messaging-service/pkg/db"
	"web-messaging-service/pkg/ws"

	"web-messaging-service/pkg/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// @title Web Messaging Service API
// @version 1.0
// @description This is a web messaging service API.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	config.Load()
	err := db.Connect()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.DB.Close()

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	ws.SetupWebSocket(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	messageService := services.NewMessageService()
	routes.Setup(router, messageService)

	router.Static("/client1", "./clients/client1")
	router.Static("/client2", "./clients/client2")
	router.Static("/client3", "./clients/client3")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan struct{})
	defer close(quit)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %s", err)
	}
	log.Println("Server exiting")
}
