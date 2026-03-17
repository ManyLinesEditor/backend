package main

import (
	"ManyLines/internal/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", handlers.Default)
	if err := router.Run(); err != nil {
		log.Panicf("failed to run server: %v", err)
	}

}
