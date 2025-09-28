package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "telemetry_api/docs"

	"telemetry_api/providers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title TelemetryAPI
// @version 1.0
// @description Telemetry API for Smart Home System
// @host localhost:8083
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	initKafka()
	defer kafkaWriter.Close()

	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	deviceRepo := NewDeviceRepository(db)

	providerFactory := providers.NewProviderFactory()

	deviceService := NewDeviceService(deviceRepo)
	telemetryService := NewTelemetryService(deviceRepo, providerFactory)

	apiHandler := NewApiHandler(deviceService, telemetryService)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiRoutes := router.Group("/api/v1")
	{
		apiHandler.RegisterRoutes(apiRoutes)
	}

	port := getEnv("PORT", ":8083")
	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		log.Printf("Telemetry API server starting on %s\n", srv.Addr)
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
