package handlers

import (
	"api-gateway/lib"
	"api-gateway/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

type serviceGateway struct {
	serv service.ServiceGateway
}

func NewHandlersGateway(service service.ServiceGateway) *serviceGateway {
	return &serviceGateway{
		serv: service,
	}
}

func (s *serviceGateway) ProxyHandler(ctx *fiber.Ctx) error {
	serviceName := ctx.Params("service")
	path := ctx.Params("path")

	data, ok := s.serv.GetDataByServiceName(serviceName)
	if !ok {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "service not found",
		})
	}

	host := data["host"].(string)
	ctx.Request().URI().SetHost(host)
	ctx.Request().URI().SetScheme("http")
	ctx.Request().URI().SetPath(path)
	fullPath := ctx.Request().URI().String()

	log.Println("ProxyHandler fullPath: ", fullPath)

	if err := proxy.Do(ctx, fullPath); err != nil {
		return lib.ErrorHandler(ctx, &lib.InternalServerError{Message: "Failed to proxy request", MessageDev: err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
