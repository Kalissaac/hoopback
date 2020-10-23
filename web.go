package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"
)

var (
	app = fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).SendString(err.Error())
		},
	})
)

func setupMiddleware() {
	app.Use(compress.New())
	app.Use(limiter.New())
}

func setupRoutes() {
	app.Get("/login", func(c *fiber.Ctx) error {
		fmt.Println(c.Params("code"))
		redirectURI := c.BaseURL() + "/login"
		if c.Cookies("authorization") != "" {
			// Authorization cookie found, user is authenticated
			return c.Redirect("/")
		} else if c.Params("code") != "" {
			// Authorization code from Discord found
			data := url.Values{
				"client_id":     {os.Getenv("CLIENT_ID")},
				"client_secret": {os.Getenv("CLIENT_SECRET")},
				"grant_type":    {"authorization_code"},
				"redirect_uri":  {redirectURI},
				"code":          {c.Params("code")},
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
			return c.Redirect("/")
		} else {
			// If no auth cookie is found, and it's not a redirect from Discord, send them to login
			return c.Redirect("https://discord.com/api/oauth2/authorize?client_id=768975292621783089&redirect_uri=" + url.QueryEscape(redirectURI) + "&response_type=code&scope=identify")
		}
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to hoopback!")
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
	port := os.Getenv("PORT")

	setupMiddleware()
	setupRoutes()
	app.Listen(":" + port)
}
