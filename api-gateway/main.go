package main

import (
	"api-gateway/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	configCors := cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Allow-Headers, Access-Control-Allow-Credentials, Origin, Accept, X-Requested-With, access-control-allow-origin, access-control-allow-methods, access-control-allow-headers, x-api-token, x-timestamp, x-signature",
		AllowMethods:     "POST, OPTIONS, GET, PUT, DELETE",
		AllowCredentials: false,
	})

	app := fiber.New()
	app.Use(configCors)
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Asia/Jakarta",
	}))

	routes.GatewayRoutes(app)
	log.Fatal(app.Listen(":4000"))
}
