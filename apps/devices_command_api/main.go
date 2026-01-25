package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "devices_command_api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title DevicesCommandAPI
// @version 1.0
// @host localhost:8082
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	initKafka()
	defer kafkaWriter.Close()
	if kafkaConn != nil {
		defer kafkaConn.Close()
	}

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	apiRoutes := router.Group("/api/v1")
	{
		deviceHandler := NewDeviceHandler()
		deviceHandler.RegisterRoutes(apiRoutes)
	}

	port := getEnv("PORT", ":8082")
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		log.Printf("Devices Command API server starting on %s\n", srv.Addr)
		log.Printf("Swagger documentation available at http://localhost%s/swagger/index.html\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited properly")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
