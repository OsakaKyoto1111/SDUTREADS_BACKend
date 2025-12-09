package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// -------------------- Load config --------------------
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// -------------------- Init DB --------------------
	db, err := database.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// -------------------- Repos & Services --------------------
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	userSvc := service.NewUserService(userRepo)

	// -------------------- Handlers --------------------
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	api := e.Group("/api")

	authGroup := api.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	userGroup := api.Group("/user")
	userGroup.Use(middleware.JWT(cfg.JWTSecret))

	userGroup.GET("/me", userHandler.GetMe)
	userGroup.PATCH("/me", userHandler.UpdateMe)
	userGroup.POST("/avatar", userHandler.UploadAvatar)
	userGroup.GET("/search", userHandler.SearchUsers)
	userGroup.POST("/:id/follow", userHandler.Follow)
	userGroup.DELETE("/:id/follow", userHandler.Unfollow)

	postRepo := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postService)

	postGroup := api.Group("/posts")
	postGroup.Use(middleware.JWT(cfg.JWTSecret))

	postGroup.POST("", postHandler.Create)
	postGroup.PATCH("/:id", postHandler.Update)
	postGroup.DELETE("/:id", postHandler.Delete)

	postGroup.POST("/:id/files", postHandler.AddFiles)

	postGroup.POST("/:id/comments", postHandler.AddComment)

	postGroup.POST("/:id/like", postHandler.LikePost)
	postGroup.DELETE("/:id/like", postHandler.UnlikePost)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- e.Start(":" + cfg.AppPort)
	}()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server start error: %v", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := e.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown failed: %v", err)
		}
	}
}
