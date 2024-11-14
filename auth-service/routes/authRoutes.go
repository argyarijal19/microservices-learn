package routes

import (
	"auth-service/handlers"
	"auth-service/middleware"
	"auth-service/repository"
	"auth-service/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router *fiber.App) {
	repo := repository.NewAuthRepository()
	service := service.NewAuthService(repo)
	handlers := handlers.NewAuthHandlers(service)

	groupRoutes := router.Group("/auth")
	groupRoutes.Post("/", middleware.AuthTokenJWT(), middleware.VerifyAPIKey(), handlers.GenerateToken)
	groupRoutes.Get("/", middleware.AuthTokenJWT(), middleware.VerifyAPIKey(), handlers.CheckToken)
	groupRoutes.Delete("/", middleware.AuthTokenJWT(), middleware.VerifyAPIKey(), handlers.RevokeToken)
}
