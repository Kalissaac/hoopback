package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

var (
	app = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).SendString(err.Error())
		},
	})
)

func setupMiddleware() {
	app.Use(compress.New())
	app.Use(limiter.New())
}

func setupRoutes() {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to hoopback!")
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Sorry can't find that!")
	})
}

func main() {
	port := os.Getenv("PORT")

	setupMiddleware()
	setupRoutes()
	app.Listen(":" + port)
}
