// Package main Go Job Queue API
//
// @title Go Job Queue API
// @version 1.0
// @description Distributed Job Queue built using Go, Redis and PostgreSQL.
// @termsOfService http://swagger.io/terms/
//
// @contact.name Sagar Pardhi
//
// @license.name MIT
//
// @host localhost:8080
// @BasePath /
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sagar-pardhi/go-job-queue/config"
	"github.com/sagar-pardhi/go-job-queue/internal/database"
	"github.com/sagar-pardhi/go-job-queue/internal/jobs"
	redisclient "github.com/sagar-pardhi/go-job-queue/internal/redis"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/sagar-pardhi/go-job-queue/docs"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.DatabaseURL())
	if err != nil {
		log.Fatal(err)
	}

	// Redis
	redis := redisclient.NewRedis(cfg.RedisAddr)

	log.Println(cfg.DatabaseURL())
	log.Println(cfg.RedisAddr)

	// Repository
	repo := jobs.NewRepository(db)

	// Handler
	handler := jobs.NewHandler(repo, redis)

	// Router
	router := gin.Default()

	router.POST("/jobs", handler.CreateJob)
	router.GET("/jobs/:id", handler.GetJob)
	router.GET("/jobs", handler.ListJobs)
	router.GET("/metrics", handler.Metrics)
	router.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
		),
	)

	log.Println("API running on :8080")
	router.Run(":8080")
}
