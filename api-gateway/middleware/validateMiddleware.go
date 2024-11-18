package middleware

import (
	"api-gateway/lib"
	"bytes"
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
)

func ValidateMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		receivedSignature := ctx.Get("x-signature")
		timestamp := ctx.Get("x-timestamp")
		apiToken := ctx.Get("x-api-token")

		privateKey := os.Getenv("PRIVATE_KEY")

		if receivedSignature == "" || timestamp == "" || apiToken == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing required headers",
			})
		}

		var formattedBody bytes.Buffer

		if ctx.Body() != nil {
			err := json.Compact(&formattedBody, ctx.Body())
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to parse request body",
				})
			}
		}

		// Validasi x-signature
		isValid := lib.ValidateXSignature(receivedSignature, timestamp, formattedBody.String(), apiToken, privateKey)
		if !isValid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid signature",
			})
		}

		return ctx.Next()
	}
}
