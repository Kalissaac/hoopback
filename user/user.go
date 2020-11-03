package user

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/session/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	app      *fiber.App
	sessions *session.Session
	client   *mongo.Client
)

// Webhook object
type Webhook struct {
	ID              string             `bson:"_id"`
	Name            string             `bson:"name"`
	Destination     string             `bson:"destination"`
	Transformations []string           `bson:"transformations"`
	Type            string             `bson:"type"`
	Method          string             `bson:"method"`
	LastSent        primitive.DateTime `bson:"lastSent"`
	Status          string             `bson:"status"`
}

// User object
type User struct {
	ID                 string             `bson:"_id,omitempty"`
	Username           string             `bson:"username"`
	Avatar             string             `bson:"avatar"`
	AccessToken        string             `bson:"access_token"`
	AccessTokenExpires primitive.DateTime `bson:"access_token_expires"`
	RefreshToken       string             `bson:"refresh_token"`
	Webhooks           map[string]Webhook `bson:"webhooks"`
}

// Setup User routes and such
func Setup(a *fiber.App, s *session.Session, c *mongo.Client) {
	app = a
	sessions = s
	client = c

	app.Get("/home", func(c *fiber.Ctx) error {
		store := sessions.Get(c)
		var user User
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
