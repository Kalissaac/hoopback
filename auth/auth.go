package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/session/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	fetch    = resty.New()
	app      *fiber.App
	sessions *session.Session
	client   *mongo.Client
)

// DiscordTokenResponse represents the OAuth token recieved from Discord
type DiscordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// DiscordUserResponse represents the user object recieved from Discord
type DiscordUserResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

// Webhook object
type Webhook struct {
	ID              string   `bson:"_id"`
	Name            string   `bson:"name"`
	Destination     string   `bson:"destination"`
	Transformations []string `bson:"transformations"`
}

// User object
type User struct {
	ID                 string             `bson:"_id,omitempty"`
	Username           string             `bson:"username"`
	AccessToken        string             `bson:"access_token"`
	AccessTokenExpires int64              `bson:"access_token_expires"`
	RefreshToken       string             `bson:"refresh_token"`
	Webhooks           map[string]Webhook `bson:"webhooks"`
}

// Setup Auth routes and such
func Setup(a *fiber.App, s *session.Session, c *mongo.Client) {
	app = a
	sessions = s
	client = c

	app.Get("/login", func(c *fiber.Ctx) error {
		redirectURI := c.BaseURL() + "/login"

		store := sessions.Get(c)
		defer store.Save()

		if store.Get("user") != nil {
			// Authorization cookie found, user is authenticated
			return c.Redirect("/restricted")
		} else if c.Query("code") != "" {
			// Authorization code from Discord found

			resp, err := fetch.R().
				SetFormData(map[string]string{
					"client_id":     os.Getenv("CLIENT_ID"),
					"client_secret": os.Getenv("CLIENT_SECRET"),
					"grant_type":    "authorization_code",
					"redirect_uri":  redirectURI,
					"code":          c.Query("code"),
					"scope":         "identify",
				}).
				Post("https://discordapp.com/api/oauth2/token")

			if err != nil {
				panic(err)
			}

			res := DiscordTokenResponse{}

			err = json.Unmarshal(resp.Body(), &res)
			if err != nil {
				panic(err)
			}

			userInfo := User{
				AccessToken:        res.AccessToken,
				RefreshToken:       res.RefreshToken,
				AccessTokenExpires: time.Now().Unix() + res.ExpiresIn,
			}

			initUser(&userInfo)

			store.Set("user", userInfo.RefreshToken)

			return c.Redirect("/restricted")
		} else {
			// If no auth cookie is found, and it's not a redirect from Discord, send them to login
			return c.Redirect("https://discord.com/api/oauth2/authorize?client_id=" + url.QueryEscape(os.Getenv("CLIENT_ID")) + "&redirect_uri=" + url.QueryEscape(redirectURI) + "&response_type=code&scope=identify")
		}
	})

	app.Use(func(c *fiber.Ctx) error {
		store := sessions.Get(c)
		if store.Get("user") == nil {
			return c.Redirect("/login")
		}
		return c.Next()
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		store := sessions.Get(c)
		defer store.Save()
		store.Destroy()
		return c.Redirect("/")
	})
}

func _createUser(userInfo *User) {
	collection := client.Database("data").Collection("users")

	insertResult, err := collection.InsertOne(context.TODO(), userInfo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted post with ID:", insertResult.InsertedID)
}

func initUser(userInfo *User) {
	resp, err := fetch.R().
		SetAuthToken(userInfo.AccessToken).
		Get("https://discord.com/api/v8/users/@me")

	res := DiscordUserResponse{}

	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		panic(err)
	}

	userInfo.ID = res.ID
	userInfo.Username = res.Username

	collection := client.Database("data").Collection("users")

	var user User

	err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: userInfo.ID}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_createUser(userInfo)
		} else {
			log.Fatal(err)
		}
	}
}
