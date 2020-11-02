package user

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/session/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"hoopback.schwa.tech/auth"
)

var (
	app      *fiber.App
	sessions *session.Session
	client   *mongo.Client
)

// Setup User routes and such
func Setup(a *fiber.App, s *session.Session, c *mongo.Client) {
	app = a
	sessions = s
	client = c

	app.Get("/home", func(c *fiber.Ctx) error {
		store := sessions.Get(c)
		var user auth.User
		err := client.Database("data").Collection("users").FindOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusServiceUnavailable, "User not found! Are you registered?")
			}
			log.Fatal(err)
		}

		return c.Render("home", fiber.Map{
			"user": user,
		}) // , "layouts/main"
	})

	app.Get("/w/:user/:webhook/edit", func(c *fiber.Ctx) error {
		var user auth.User
		err := client.Database("data").Collection("users").FindOne(context.TODO(), bson.D{{Key: "_id", Value: c.Params("user")}}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusNotFound, "User not found! Are they registered?")
			}
			log.Fatal(err)
		}

		webhook, ok := user.Webhooks[c.Params("webhook")]
		if ok == false {
			return fiber.NewError(fiber.StatusNotFound, "Webhook not found!")
		}

		store := sessions.Get(c)
		if store.Get("user") != user.ID {
			return fiber.NewError(fiber.StatusNotFound, "Webhook not found or user is unauthorized to access this page!")
		}

		return c.Render("edithook", fiber.Map{
			"user":    user,
			"webhook": webhook,
		}) // , "layouts/main"
	})
}
