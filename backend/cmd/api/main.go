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
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.InitPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	commentLikeRepo := repository.NewCommentLikeRepository(db)
	feedRepo := repository.NewFeedRepository(db)

	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	userSvc := service.NewUserService(userRepo)
	commentLikeSvc := service.NewCommentLikeService(commentLikeRepo)
	commentSvc := service.NewCommentService(commentRepo, commentLikeRepo)
	commentTreeSvc := service.NewCommentTreeService(commentRepo, commentLikeRepo)
	postSvc := service.NewPostService(postRepo, commentSvc, commentTreeSvc)

	fileSvc := service.NewFileService("uploads", "/uploads/", 10*1024*1024, []string{"jpg", "jpeg", "png", "gif", "mp4"})

	feedSvc := service.NewFeedService(feedRepo)
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	postHandler := handler.NewPostHandler(postSvc, fileSvc)
	commentHandler := handler.NewCommentHandler(commentSvc)
	commentLikeHandler := handler.NewCommentLikeHandler(commentLikeSvc)
	feedHandler := handler.NewFeedHandler(feedSvc)
	fileHandler := handler.NewFileHandler(fileSvc)

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	api := e.Group("/api")

	feedGroup := api.Group("/feed")
	feedGroup.Use(middleware.JWT(cfg.JWTSecret))
	feedGroup.GET("", feedHandler.Get)
	api.GET("/debug/headers", func(c echo.Context) error {
		return c.JSON(200, c.Request().Header)
	})
	api.GET("/debug/token", func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		previewLen := 10
		if len(cfg.JWTSecret) < previewLen {
			previewLen = len(cfg.JWTSecret)
		}
		return c.JSON(200, map[string]interface{}{
			"authorization_header": authHeader,
			"jwt_secret_length":    len(cfg.JWTSecret),
			"jwt_secret_preview":   cfg.JWTSecret[:previewLen] + "...",
		})
	})

	authGroup := api.Group("/auth")
	authGroup.POST("/register", authHandler.Register)
	authGroup.POST("/login", authHandler.Login)

	userGroup := api.Group("/user")
	userGroup.Use(middleware.JWT(cfg.JWTSecret))
	userGroup.GET("/me", userHandler.GetProfile)
	userGroup.PATCH("/me", userHandler.Update)
	userGroup.DELETE("/me", userHandler.Delete)
	userGroup.GET("/search", userHandler.Search)
	userGroup.POST("/:id/follow", userHandler.Follow)
	userGroup.DELETE("/:id/follow", userHandler.Unfollow)
	userGroup.POST("/avatar", fileHandler.Upload)

	postGroup := api.Group("/posts")
	postGroup.Use(middleware.JWT(cfg.JWTSecret))
	postGroup.POST("", postHandler.Create)
	postGroup.GET("/:id", postHandler.Get)
	postGroup.PATCH("/:id", postHandler.Update)
	postGroup.DELETE("/:id", postHandler.Delete)
	postGroup.POST("/:id/files", postHandler.AddFiles)
	postGroup.POST("/:id/like", postHandler.Like)
	postGroup.DELETE("/:id/like", postHandler.Unlike)
	postGroup.GET("/:id/comments", commentHandler.GetTree)
	postGroup.POST("/:id/comments", commentHandler.Add)
	postGroup.POST("/comments/:comment_id/like", commentLikeHandler.Like)
	postGroup.DELETE("/comments/:comment_id/like", commentLikeHandler.Unlike)

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
