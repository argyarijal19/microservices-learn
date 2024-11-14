package middleware

import (
	"auth-service/lib"
	"crypto/sha256"
	"encoding/hex"
	"os"

	"github.com/gofiber/fiber/v2"
)

func VerifyAPIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key") // Mengambil X-API-Key dari header request
		if apiKey == "" {
			return lib.ErrorHandler(c, &lib.UnauthorizedError{Message: "API Key is required"})
		}

		if validKey := compareKeyWithHash(os.Getenv("API_KEY"), apiKey, 30); validKey {
			c.Next()
			return nil
		}

		if apiKey == os.Getenv("API_KEY") {
			c.Next()
			return nil
		}
		return lib.ErrorHandler(c, &lib.UnauthorizedError{Message: "Invalid API Key"})
	}
}

func AuthTokenJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return lib.ErrorHandler(c, &lib.UnauthorizedError{Message: "Token is required"})
		}

		result, err := lib.VerifyJWTToken(token, os.Getenv("SECRET_KEY_JWT"))
		if err != nil {
			return lib.ErrorHandler(c, &lib.UnauthorizedError{Message: "Invalid Token", MessageDev: err.Error()})
		}

		customClaims, ok := result.(*lib.JWTCustomClaims)
		if !ok {
			return lib.ErrorHandler(c, &lib.UnauthorizedError{Message: "Invalid Token Structure"})
		}

		if customClaims != nil {
			c.Next()
			return nil
		}

		return c.Next()
	}
}

func ShortHash(key string, length int) (string, error) {
	hash := sha256.Sum256([]byte(key))
	shortHash := hex.EncodeToString(hash[:])[:length]
	return shortHash, nil
}

func compareKeyWithHash(key, hashToCompare string, length int) bool {
	generatedHash, _ := ShortHash(key, length)
	return generatedHash == hashToCompare
}
