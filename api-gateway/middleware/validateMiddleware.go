package middleware

import (
	"api-gateway/lib"
	"os"

	"github.com/gofiber/fiber/v2"
)

func ValidateMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		receivedSignature := ctx.Get("x-signature")
		timestamp := ctx.Get("x-timestamp")
		apiToken := ctx.Get("x-api-token")

		privateKey := os.Getenv("PRIVATE_KEY")

		body := ctx.Body()

		if receivedSignature == "" || timestamp == "" || apiToken == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing required headers",
			})
		}
		// Validasi x-signature
		isValid := lib.ValidateXSignature(receivedSignature, timestamp, string(body), apiToken, privateKey)
		if !isValid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid signature",
			})
		}

		return ctx.Next()
	}
}
