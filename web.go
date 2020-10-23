package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	port := os.Getenv("PORT")
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("welcome to hoopback")
	})

	app.Listen(":" + port)
}
