package main

import (
	"context"
	"log"
	"os"
	"time"

	"hoopback.schwa.tech/auth"
	"hoopback.schwa.tech/user"
	"hoopback.schwa.tech/webhook"

	"github.com/go-resty/resty/v2"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/session/v2"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	app = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).SendString(err.Error())
		},
	})
	sessions = session.New(session.Config{
		Expiration: 7 * 24 * time.Hour,
		Secure:     false, // TODO: change to true for production
	})
	client *mongo.Client
	fetch  = resty.New()
)

func setupMiddleware() {
	app.Use(compress.New())
	app.Use(helmet.New())
	// app.Use(limiter.New())
	app.Use(recover.New())
}

func setupRoutes() {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to hoopback!")
	})

	auth.Setup(app, sessions, client)

	app.Get("/restricted", func(c *fiber.Ctx) error {
		return c.SendString("restricted area!1!!")
	})

	user.Setup(app, sessions, client)
	webhook.Setup(app, client)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Sorry can't find that!")
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	defer client.Disconnect(ctx)

	setupMiddleware()
	setupRoutes()

	app.Listen(":" + os.Getenv("PORT"))
}
