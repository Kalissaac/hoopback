package webhook

import (
	"context"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"hoopback.schwa.tech/auth"

	"github.com/go-resty/resty/v2"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	fetch  = resty.New()
	app    *fiber.App
	client *mongo.Client
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
}

// Setup webhook routes and such
func Setup(a *fiber.App, c *mongo.Client) {
	app = a
	client = c

	app.Post("/w/:user/:webhook", func(c *fiber.Ctx) error {
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

		webhook.LastSent = primitive.NewDateTimeFromTime(time.Now())
		user.Webhooks[webhook.ID] = webhook

		update := bson.D{{Key: "$set", Value: bson.D{{Key: "webhooks", Value: user.Webhooks}}}}
		client.Database("data").Collection("users").UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: user.ID}}, update)

		body := make(map[string]interface{})
		if err := c.BodyParser(&body); err != nil {
			return err
		}

		if webhook.Type == "basic" {
			executeWebhook(webhook.Destination, body, webhook.Transformations)
		}

		return c.JSON(&fiber.Map{
			"success": true,
		})
	})
}
