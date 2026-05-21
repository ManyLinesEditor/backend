package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ManyLinesEditor/backend/payment/config"
	_ "github.com/ManyLinesEditor/backend/payment/docs"
	"github.com/ManyLinesEditor/backend/payment/internal/db"
	"github.com/ManyLinesEditor/backend/payment/internal/handlers"
	"github.com/ManyLinesEditor/backend/payment/internal/middleware"
	"github.com/ManyLinesEditor/backend/payment/internal/migrations"
	"github.com/ManyLinesEditor/backend/payment/internal/repositories/postgres"
	"github.com/ManyLinesEditor/backend/payment/internal/services"
	"github.com/ManyLinesEditor/backend/payment/internal/sse"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           ManyLines Payment API
// @version         1.0
// @description     Payment service for ManyLines Editor
// @host            localhost:8082
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Format: Bearer <token>
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("postgres ping: %v", err)
	}

	migrations.RunMigrations(cfg.DatabaseURL)

	deviceRepo := postgres.NewDeviceRepo(db)
	paymentRepo := postgres.NewPaymentRepo(db)
	subscriptionRepo := postgres.NewSubscriptionRepo(db)

	broker := sse.NewBroker()

	paymentService := services.NewPaymentService(paymentRepo, cfg.AcquiremockURL, cfg.BaseURL)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)
	webhookService := services.NewWebhookService(paymentRepo, subscriptionRepo, broker, cfg.WebhookSecret)

	deviceH := handlers.NewDeviceHandler(deviceRepo)
	paymentH := handlers.NewPaymentHandler(paymentService)
	webhookH := handlers.NewWebhookHandler(webhookService)
	sseH := handlers.NewSSEHandler(broker)
	subH := handlers.NewSubscriptionHandler(subscriptionService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/payment/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/webhooks/acquiremock", webhookH.Handle)

	authorized := r.Group("/", middleware.GatewayAuth())
	{
		authorized.POST("/devices", deviceH.Register)
		authorized.POST("/payments", paymentH.Create)
		authorized.GET("/payments", paymentH.List)
		authorized.GET("/subscriptions", subH.ListActive)
		authorized.GET("/sse/features", sseH.Stream)
	}

	log.Printf("starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
