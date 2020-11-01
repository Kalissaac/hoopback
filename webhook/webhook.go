package webhook

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hoopback.schwa.tech/auth"
)

var (
	app    *fiber.App
	client *mongo.Client
	fetch  = resty.New()
)

func executeWebhook(destination string, body map[string]interface{}, transformations []string) {
	transformedBody := make(map[string]interface{})
	for _, rawOperator := range transformations {
		operator := strings.Split(rawOperator, ":")
		if field := body[operator[0]]; field != nil {
			transformedBody[operator[1]] = field
		}
	}
	fetch.R().
		SetBody(transformedBody).
		Post(destination)
	fmt.Println("here webhook, to", destination)
}

// Setup webhook routes and such
func Setup(a *fiber.App, c *mongo.Client) {
	app = a
	client = c

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

		body := make(map[string]interface{})
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		if webhook.Type == "basic" {
			executeWebhook(webhook.Destination, body, webhook.Transformations)
		}
		return c.SendString("Success")
	})
}
