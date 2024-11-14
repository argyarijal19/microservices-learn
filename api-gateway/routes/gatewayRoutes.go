package routes

import (
	"api-gateway/handlers"
	"api-gateway/middleware"
	"api-gateway/repository"
	"api-gateway/service"

	"github.com/gofiber/fiber/v2"
)

func GatewayRoutes(route *fiber.App) {
	repo := repository.NewServiceRepository()
	serv := service.NewServiceGateway(repo)
	handlers := handlers.NewHandlersGateway(serv)

	route.All("/:service/:path", middleware.ValidateMiddleware(), handlers.ProxyHandler)
}
