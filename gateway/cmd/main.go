package main

import (
	"context"
	"log"

	"github.com/ManyLinesEditor/backend/gateway/config"
	_ "github.com/ManyLinesEditor/backend/gateway/docs"
	"github.com/ManyLinesEditor/backend/gateway/internal/handlers"
	"github.com/ManyLinesEditor/backend/gateway/internal/middleware"
	"github.com/ManyLinesEditor/backend/gateway/internal/migrations"
	"github.com/ManyLinesEditor/backend/gateway/internal/proxy"
	"github.com/ManyLinesEditor/backend/gateway/internal/repositories/postgres"
	"github.com/ManyLinesEditor/backend/gateway/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           ManyLines Gateway API
// @version         1.0
// @description     API gateway for ManyLines Editor
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Format: Bearer <token>
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("env variables load: %v", err)
	}

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("postgres connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("postgres ping: %v", err)
	}

	migrations.RunMigrations(cfg.DatabaseURL)

	userRepo := postgres.NewUserRepo(db)

	tokens := services.NewTokenService(cfg.TokenSecret, cfg.TokenTTL)
	authSvc := services.NewAuthService(userRepo, tokens)

	storageProxy := proxy.Handler(cfg.StorageURL)
	paymentProxy := proxy.Handler(cfg.PaymentURL)

	authH := handlers.NewAuthHandler(authSvc)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/payment/swagger/*any", proxy.Handler(cfg.PaymentURL))
	r.GET("/storage/swagger/*any", proxy.Handler(cfg.StorageURL))

	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)

	r.POST("/webhooks/*path", paymentProxy)

	authed := r.Group("/", middleware.Auth(tokens))
	{
		authed.Any("/files", storageProxy)
		authed.Any("/files/*path", storageProxy)

		authed.Any("/devices", paymentProxy)
		authed.Any("/payments", paymentProxy)
		authed.Any("/payments/*path", paymentProxy)
		authed.Any("/subscriptions", paymentProxy)
		authed.Any("/sse/features", paymentProxy)
	}

	log.Printf("starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
