package main

import (
	"database/sql"
	"fmt"
	"log"
	"ozinse-backend/internal/config"
	"ozinse-backend/internal/handler"
	"ozinse-backend/internal/repository"
	"ozinse-backend/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("База данных недоступна: %v", err)
	}

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	authRepo := repository.NewAuthRepository(db)
	authSvc := service.NewAuthService(authRepo, cfg)
	authHandler := handler.NewAuthHandler(authSvc)

	projectRepo := repository.NewProjectRepository(db)
	projectSvc := service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectSvc)

	favRepo := repository.NewFavoriteRepository(db)
	favSvc := service.NewFavoriteService(favRepo, projectRepo)
	favHandler := handler.NewFavoriteHandler(favSvc)

	refRepo := repository.NewReferenceRepository(db)
	refSvc := service.NewReferenceService(refRepo)
	refHandler := handler.NewReferenceHandler(refSvc)

	adminRepo := repository.NewAdminRepository(db)
	adminSvc := service.NewAdminService(adminRepo)
	adminHandler := handler.NewAdminHandler(adminSvc)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		projects := v1.Group("/projects")
		{
			projects.GET("", projectHandler.GetAll)
			projects.GET("/featured", projectHandler.GetFeatured)
			projects.GET("/:id", projectHandler.GetByID)
			projects.GET("/:id/seasons", projectHandler.GetSeasons)
		}

		v1.GET("/categories", refHandler.GetCategories)
		v1.GET("/genres", refHandler.GetGenres)
		v1.GET("/age-ratings", refHandler.GetAgeRatings)

		secureGroup := v1.Group("")
		secureGroup.Use(handler.AuthMiddleware(cfg))
		{
			profile := secureGroup.Group("/profile")
			{
				profile.GET("/me", authHandler.GetMe)
				profile.PUT("/me", authHandler.UpdateMe)
				profile.PUT("/password", authHandler.UpdatePassword)
			}

			favorites := secureGroup.Group("/favorites")
			{
				favorites.GET("", favHandler.GetFavorites)
				favorites.POST("/:project_id", favHandler.AddFavorite)
				favorites.DELETE("/:project_id", favHandler.DeleteFavorite)
			}

			admin := secureGroup.Group("/admin")
			admin.Use(handler.AdminMiddleware(*adminRepo))
			{
				admin.POST("/upload", handler.UploadFile)

				admin.POST("/categories", adminHandler.CreateCategory)
				admin.PUT("/categories/:id", adminHandler.UpdateCategory)
				admin.DELETE("/categories/:id", adminHandler.DeleteCategory)

				admin.POST("/genres", adminHandler.CreateGenre)
				admin.PUT("/genres/:id", adminHandler.UpdateGenre)
				admin.DELETE("/genres/:id", adminHandler.DeleteGenre)

				admin.POST("/age-ratings", adminHandler.CreateAgeRating)
				admin.PUT("/age-ratings/:id", adminHandler.UpdateAgeRating)
				admin.DELETE("/age-ratings/:id", adminHandler.DeleteAgeRating)

				admin.GET("/users", adminHandler.GetUsers)
				admin.POST("/users/:user_id/assign-role", adminHandler.AssignRole)

				admin.GET("/roles", adminHandler.GetRoles)
				admin.POST("/roles", adminHandler.CreateRole)
				admin.PUT("/roles/:id", adminHandler.UpdateRole)
				admin.DELETE("/roles/:id", adminHandler.DeleteRole)

				admin.GET("/projects", adminHandler.GetProjects)
				admin.POST("/projects", adminHandler.CreateProject)
				admin.PUT("/projects/:id", adminHandler.UpdateProject)
				admin.DELETE("/projects/:id", adminHandler.DeleteProject)
				admin.POST("/projects/:id/seasons", adminHandler.CreateSeason)
				admin.PUT("/projects/featured-order", adminHandler.UpdateFeaturedOrder)

				admin.POST("/seasons/:season_id/episodes", adminHandler.CreateEpisode)
			}
		}
	}

	fmt.Printf("Сервер запущен на порту %s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}