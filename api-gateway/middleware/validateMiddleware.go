package middleware

import (
	"api-gateway/lib"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ValidateMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		receivedSignature := ctx.Get("x-signature")
		timestamp := ctx.Get("x-timestamp")
		apiToken := ctx.Get("x-api-token")

		privateKey := os.Getenv("PRIVATE_KEY")

		if receivedSignature == "" || timestamp == "" || apiToken == "" {
			return lib.ErrorHandler(ctx, &lib.BadRequestError{
				Message:    "Missing required headers",
				MessageDev: "Missing required headers",
			})
		}

		var formattedBody bytes.Buffer

		if ctx.Body() != nil {
			err := json.Compact(&formattedBody, ctx.Body())
			if err != nil {
				return lib.ErrorHandler(ctx, &lib.InternalServerError{
					Message:    "Failed to format request body",
					MessageDev: fmt.Sprintf("Failed to format request body: %v", err),
				})
			}
		}

		// Validasi x-signature
		isValid := lib.ValidateXSignature(receivedSignature, timestamp, formattedBody.String(), apiToken, privateKey)
		if !isValid {
			return lib.ErrorHandler(ctx, &lib.UnauthorizedError{
				Message:    "Invalid signature",
				MessageDev: "Invalid signature",
			})
		}

		if err := validateTimestamp(timestamp, 5*time.Minute); err != nil {
			return lib.ErrorHandler(ctx, &lib.UnauthorizedError{
				Message:    "Invalid timestamp",
				MessageDev: fmt.Sprintf("Invalid timestamp: %v", err),
			})
		}

		return ctx.Next()
	}
}

func validateTimestamp(timestamp string, allowedSkew time.Duration) error {
	clientTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return err
	}

	serverTime := time.Now().UTC()

	timeDiff := serverTime.Sub(clientTime)

	if timeDiff > allowedSkew || timeDiff < -allowedSkew {
		return fmt.Errorf("timestamp out of allowed range")
	}

	return nil
}
