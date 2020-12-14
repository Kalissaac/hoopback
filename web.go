package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"time"

	"hoopback.schwa.tech/api"
	"hoopback.schwa.tech/auth"
	"hoopback.schwa.tech/user"
	"hoopback.schwa.tech/webhook"

	"github.com/go-resty/resty/v2"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/storage/mongodb"
	"github.com/gofiber/template/html"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	engine = html.New("./views", ".html").AddFunc("getDomain", func(url string) string {
		return regexp.MustCompile(`([\w]+\.){1}([\w]+\.?)+`).FindString(url)
	}).AddFunc("trimString", func(str string, length int) string {
		if len(str) < length {
			return str
		}
		a := []rune(str)
		return string(a[0:length]) + "..."
	}) // fix this mess
	app = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).SendString(err.Error())
		},
		Views: engine,
	})
	sessions *session.Store
	client   *mongo.Client
	fetch    = resty.New()
)

func setupMiddleware() {
	app.Use(compress.New())
	app.Use(helmet.New())
	// app.Use(limiter.New())
	app.Use(recover.New())

	sessions = session.New(session.Config{
		Expiration: 7 * 24 * time.Hour,
		Storage: mongodb.New(mongodb.Config{
			ConnectionURI: os.Getenv("MONGO_URI"),
			Database:      "data",
			Collection:    "fiber_storage",
			Reset:         false,
		}),
	})
}

func setupRoutes() {
	app.Static("/", "./static")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./static/index.html")
	})

	webhook.Setup(app, client)

	auth.Setup(app, sessions, client)

	api.Setup(app, sessions, client)

	user.Setup(app, sessions, client)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Sorry can't find that!")
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
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
