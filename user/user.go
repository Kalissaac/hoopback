package user

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	app      *fiber.App
	sessions *session.Store
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
func Setup(a *fiber.App, s *session.Store, cl *mongo.Client) {
	app = a
	sessions = s
	client = cl

	app.Get("/home", func(c *fiber.Ctx) error {
		store, _ := sessions.Get(c)
		var user User
		err := client.Database("data").Collection("users").FindOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusServiceUnavailable, "User not found! Are you registered?")
			}
			log.Fatal(err)
		}

		return c.Render("home", fiber.Map{
			"title": "Your Webhooks",
			"user":  user,
		}) // , "layouts/main"
	})

	app.Get("/activity", func(c *fiber.Ctx) error {
		store, _ := sessions.Get(c)
		var user User
		err := client.Database("data").Collection("users").FindOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusServiceUnavailable, "User not found! Are you registered?")
			}
			log.Fatal(err)
		}

		return c.Render("activity", fiber.Map{
			"title": "Webhook Activity",
			"user":  user,
		}) // , "layouts/main"
	})

	app.Get("/settings", func(c *fiber.Ctx) error {
		store, _ := sessions.Get(c)
		var user User
		err := client.Database("data").Collection("users").FindOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusServiceUnavailable, "User not found! Are you registered?")
			}
			log.Fatal(err)
		}

		return c.Render("settings", fiber.Map{
			"title": "User Settings",
			"user":  user,
		}) // , "layouts/main"
	})

	app.Get("/webhooks/new", func(c *fiber.Ctx) error {
		store, _ := sessions.Get(c)
		var user User
		err := client.Database("data").Collection("users").FindOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusServiceUnavailable, "User not found! Are you registered?")
			}
			log.Fatal(err)
		}

		return c.Render("newhook", fiber.Map{
			"title": "New Webhook",
			"user":  user,
		}) // , "layouts/main"
	})

	app.Get("/webhooks/success", func(c *fiber.Ctx) error {
		store, _ := sessions.Get(c)
		var user User
		err := client.Database("data").Collection("users").FindOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusServiceUnavailable, "User not found! Are you registered?")
			}
			log.Fatal(err)
		}

		webhook, ok := user.Webhooks[c.Query("id")]
		if !ok {
			return c.Redirect("/home")
		}

		return c.Render("newhooksuccess", fiber.Map{
			"title":   "Webhook Creation Success",
			"user":    user,
			"webhook": webhook,
		}) // , "layouts/main"
	})

	app.Get("/webhooks/edit/:webhook", func(c *fiber.Ctx) error {
		store, _ := sessions.Get(c)
		var user User
		err := client.Database("data").Collection("users").FindOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusNotFound, "User not found! Are they registered?")
			}
			log.Fatal(err)
		}

		webhook, ok := user.Webhooks[c.Params("webhook")]
		if !ok {
			return fiber.NewError(fiber.StatusNotFound, "Webhook not found!")
		}

		return c.Render("edithook", fiber.Map{
			"title":   "Edit Webhook",
			"user":    user,
			"webhook": webhook,
		}) // , "layouts/main"
	})
}
