package handlers

import (
	"auth-service/lib"
	"auth-service/model"
	"auth-service/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type authHandlers struct {
	service service.AuthService
}

func NewAuthHandlers(s service.AuthService) *authHandlers {
	return &authHandlers{
		service: s,
	}
}

func (h *authHandlers) GenerateToken(ctx *fiber.Ctx) error {
	var request model.User

	if err := ctx.BodyParser(&request); err != nil {
		return &lib.BadRequestError{
			Message:    "invalid request",
			MessageDev: err.Error(),
		}
	}

	if request.ClientID == "" || request.ClientSecret == "" {
		return &lib.BadRequestError{
			Message:    "client_id and client_secret are required",
			MessageDev: "client_id and client_secret are required",
		}
	}

	tokenResponse, err := h.service.GenerateToken(request)
	if err != nil {
		return lib.ErrorHandler(ctx, err)
	}

	res := &lib.ResponseData{
		Code:    http.StatusOK,
		Message: "success",
		Data:    tokenResponse,
		Status:  true,
	}

	return ctx.Status(http.StatusOK).JSON(res)
}

func (h *authHandlers) CheckToken(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")
	if token == "" {
		return lib.ErrorHandler(ctx, &lib.UnauthorizedError{
			Message:    "token is required",
			MessageDev: "token is required",
		})
	}

	userData, err := h.service.CheckToken(token)
	if err != nil {
		return lib.ErrorHandler(ctx, err)
	}

	res := &lib.ResponseData{
		Code:    http.StatusOK,
		Message: "success",
		Data:    userData,
	}
	return ctx.Status(http.StatusOK).JSON(res)
}

func (h *authHandlers) RevokeToken(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")
	if token == "" {
		return lib.ErrorHandler(ctx, &lib.UnauthorizedError{
			Message:    "token is required",
			MessageDev: "token is required",
		})
	}

	if err := h.service.RevokeToken(token); err != nil {
		return lib.ErrorHandler(ctx, err)
	}

	res := &lib.ResponseData{
		Code:    http.StatusOK,
		Message: "success",
		Status:  true,
	}

	return ctx.Status(http.StatusOK).JSON(res)
}
