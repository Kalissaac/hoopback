package api

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"hoopback.schwa.tech/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/session/v2"

	"github.com/oklog/ulid/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	app      *fiber.App
	sessions *session.Session
	client   *mongo.Client
	t        = time.Now()
	entropy  = ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
)

type webhookRequest struct {
	Destination     string   `json:"destination" bson:"destination" form:"destination" validate:"required"`
	Name            string   `json:"name" bson:"name" form:"name" validate:"required"`
	Transformations []string `json:"transformations" bson:"transformations" form:"transformations" validate:"required,min=1"`
}

// Setup API routes and such
func Setup(a *fiber.App, s *session.Session, c *mongo.Client) {
	app = a
	client = c
	sessions = s

	app.Post("/api/v1/webhooks/create", func(c *fiber.Ctx) error {
		store := sessions.Get(c)
		usersCollection := client.Database("data").Collection("users")
		var user auth.User
		err := usersCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusNotFound, "User not found! Are they registered?")
			}
			log.Fatal(err)
		}

		var body webhookRequest
		if err := c.BodyParser(&body); err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		// Correct array parsing
		if len(body.Transformations) > 0 {
			body.Transformations = strings.Split(body.Transformations[0], ",")
		}

		var id string
		idRaw, err := ulid.New(ulid.Timestamp(t), entropy)
		if err != nil {
			id = primitive.NewObjectID().Hex()
		} else {
			id = strings.ToLower(idRaw.String()) + primitive.NewObjectID().Hex()
		}

		newWebhook := auth.Webhook{
			Destination:     body.Destination,
			Name:            body.Name,
			Transformations: body.Transformations,
			ID:              id,
			Type:            "basic",
			Method:          "post",
		}
		user.Webhooks[newWebhook.ID] = newWebhook
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "webhooks", Value: user.Webhooks}}}}

		usersCollection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}, update)

		return c.JSON(&fiber.Map{
			"success": true,
			"webhook": newWebhook,
		})
	})
}
