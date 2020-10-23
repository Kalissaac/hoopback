package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	app = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).SendString(err.Error())
		},
	})
	client *mongo.Client
)

func setupMiddleware() {
	app.Use(compress.New())
	// app.Use(limiter.New())
}

func setupAuth() {
	app.Get("/login", func(c *fiber.Ctx) error {
		redirectURI := c.BaseURL() + "/login"

		if c.Cookies("authorization") != "" {
			// Authorization cookie found, user is authenticated
			return c.Redirect("/restricted")
		} else if c.Query("code") != "" {
			// Authorization code from Discord found
			data := url.Values{
				"client_id":     {os.Getenv("CLIENT_ID")},
				"client_secret": {os.Getenv("CLIENT_SECRET")},
				"grant_type":    {"authorization_code"},
				"redirect_uri":  {redirectURI},
				"code":          {c.Query("code")},
				"scope":         {"identify"},
			}

			resp, err := http.PostForm("https://discordapp.com/api/oauth2/token", data)
			if err != nil {
				panic(err)
			}

			var res map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&res)

			accessToken, ok := res["access_token"].(string)
			if ok == false {
				panic("access_token casting failed!")
			}

			// Create cookie
			cookie := new(fiber.Cookie)
			cookie.Name = "authorization"
			cookie.Value = accessToken
			cookie.Expires = time.Now().Add(7 * 24 * time.Hour)

			// Set cookie
			c.Cookie(cookie)

			return c.Redirect("/restricted")
		} else {
			// If no auth cookie is found, and it's not a redirect from Discord, send them to login
			return c.Redirect("https://discord.com/api/oauth2/authorize?client_id=768975292621783089&redirect_uri=" + url.QueryEscape(redirectURI) + "&response_type=code&scope=identify")
		}
	})

	app.Use(func(c *fiber.Ctx) error {
		// TODO: Add authorization token verification
		if c.Cookies("authorization") == "" {
			return c.Redirect("/login")
		}
		return c.Next()
	})
}

// User object
type User struct {
	ID           string `json:"_id,omitempty"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	RefreshToken string `json:"refresh_token"`
}

func _createUser(userID *string) {
	collection := client.Database("data").Collection("users")

	user := User{
		ID:    *userID,
		Title: "bob",
		Body:  "john",
	}

	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
}

func initUser(userID *string) {
	collection := client.Database("data").Collection("users")

	var user User

	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: userID}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_createUser(userID)
			return
		}
		log.Fatal(err)
	}
	return
}

func setupRoutes() {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to hoopback!")
	})

	setupAuth()

	app.Get("/restricted", func(c *fiber.Ctx) error {
		return c.SendString("restricted area!1!!")
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Sorry can't find that!")
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	setupMiddleware()
	setupRoutes()

	port := os.Getenv("PORT")
	app.Listen(":" + port)

	client, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	defer cancel()
}
