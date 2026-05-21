package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ManyLinesEditor/backend/storage/config"
	_ "github.com/ManyLinesEditor/backend/storage/docs"
	"github.com/ManyLinesEditor/backend/storage/internal/handlers"
	"github.com/ManyLinesEditor/backend/storage/internal/middleware"
	"github.com/ManyLinesEditor/backend/storage/internal/migrations"
	"github.com/ManyLinesEditor/backend/storage/internal/repository/postgres"
	"github.com/ManyLinesEditor/backend/storage/internal/services"
	"github.com/ManyLinesEditor/backend/storage/internal/storage"
)

// @title           ManyLines Storage API
// @version         1.0
// @description     File storage service for ManyLines Editor
// @host            localhost:8081
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Format: Bearer <token>
func main() {
	cfg := config.Load()

	db, err := pgxpool.New(context.Background(), cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("postgres connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("postgres ping: %v", err)
	}

	migrations.RunMigrations(cfg.Postgres.DSN)

	minio, err := storage.NewMinioStorage(cfg.MinIO)
	if err != nil {
		log.Fatalf("minio connect: %v", err)
	}

	fileRepo := postgres.NewFileRepo(db)

	fileSvc := services.NewFileService(fileRepo, minio)

	fileH := handlers.NewFileHandler(fileSvc)

	r := gin.Default()

	r.GET("/storage/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	files := r.Group("/files", middleware.GatewayAuth())
	{
		files.POST("", fileH.Upload)
		files.GET("/:id", fileH.Download)
	}

	r.GET("/", handlers.Default)

	if err := r.Run(cfg.HTTP.Addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
