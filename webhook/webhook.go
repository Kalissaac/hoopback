package webhook

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hoopback.schwa.tech/auth"
)

var (
	app    *fiber.App
	client *mongo.Client
)

func executeWebhook(destination string, transformations []string) {
	fmt.Println("here webhook, to", destination)
	return
}

// Setup webhook routes and such
func Setup(a *fiber.App, c *mongo.Client) {
	app.Post("/w/:user/:webhook", func(c *fiber.Ctx) error {
		var user auth.User
		err := client.Database("data").Collection("users").FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: c.Params("user")}}).Decode(&user)
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

		executeWebhook(webhook.Destination, webhook.Transformations)
		return c.SendString("Success")
	})
}
