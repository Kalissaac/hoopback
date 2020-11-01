package user

import (
	"context"
	"fmt"
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

	app.Get("/@me", func(c *fiber.Ctx) error {
		store := sessions.Get(c)
		var user auth.User
		err := client.Database("data").Collection("users").FindOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusServiceUnavailable, "User not found! Are you registered?")
			}
			log.Fatal(err)
		}
		return c.SendString(fmt.Sprintf("Welcome to hoopback, %s!", user.Username))
	})
}
