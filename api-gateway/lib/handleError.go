package lib

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var statusCode int

	// Type assertion switch untuk error handling
	switch e := err.(type) {
	case *NotFoundError:
		statusCode = http.StatusNotFound
		logError("NotFoundError", e.MessageDev, e.Error())
	case *BadRequestError:
		statusCode = http.StatusBadRequest
		logError("BadRequestError", e.MessageDev, e.Error())
	case *InternalServerError:
		statusCode = http.StatusInternalServerError
		logError("InternalServerError", e.MessageDev, e.Error())
	case *UnauthorizedError:
		statusCode = http.StatusUnauthorized
		logError("UnauthorizedError", e.MessageDev, e.Error())
	default:
		statusCode = http.StatusInternalServerError
		logError("UnknownError", "Unknown developer message", e.Error())
	}

	// Mengirim response error ke client
	response := Response(ResponseParams{StatusCode: statusCode, Message: err.Error()})
	return c.Status(statusCode).JSON(response)
}

// Log error dengan detail
func logError(errorType, developerMessage, errorMessage string) {
	fmt.Println()
	log.Println("--------Error Handler-------------")
	log.Printf("%s, Message: %s and Developer: %s", errorType, errorMessage, developerMessage)
	log.Println("---------------------")
}
