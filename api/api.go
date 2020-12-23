package api

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	userPack "hoopback.schwa.tech/user"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/oklog/ulid/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	app      *fiber.App
	sessions *session.Store
	client   *mongo.Client
	t        = time.Now()
	entropy  = ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
)

type webhookRequest struct {
	Destination     string   `json:"destination" bson:"destination" form:"destination" validate:"required"`
	Name            string   `json:"name" bson:"name" form:"name" validate:"required"`
	Transformations []string `json:"transformations" bson:"transformations" form:"transformations" validate:"required,min=1"`
	Website         bool     `form:"web"`
}

type errorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func validateStruct(webhook interface{}) []*errorResponse {
	var errors []*errorResponse
	validate := validator.New()
	err := validate.Struct(webhook)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element errorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

// Setup API routes and such
func Setup(a *fiber.App, s *session.Store, c *mongo.Client) {
	app = a
	client = c
	sessions = s

	api := app.Group("/api")

	v1 := api.Group("/v1")

	v1.Post("/webhooks/create", func(c *fiber.Ctx) error {
		store, _ := sessions.Get(c)
		usersCollection := client.Database("data").Collection("users")
		var user userPack.User
		err := usersCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return fiber.NewError(fiber.StatusNotFound, "User not found! Are they registered?")
			}
			log.Println(err)
		}

		var body webhookRequest
		if err := c.BodyParser(&body); err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		errors := validateStruct(body)
		if errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errors)
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

		newWebhook := userPack.Webhook{
			Destination:     body.Destination,
			Name:            body.Name,
			Transformations: body.Transformations,
			ID:              id,
			Type:            "basic",
			Method:          "post",
			Status:          "active",
		}
		user.Webhooks[newWebhook.ID] = newWebhook
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "webhooks", Value: user.Webhooks}}}}

		usersCollection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}, update)

		if body.Website == true {
			return c.Redirect("/webhooks/success?id=" + newWebhook.ID)
		}

		return c.JSON(&fiber.Map{
			"success": true,
			"webhook": newWebhook,
		})
	})

	v1.Post("/webhooks/edit", editWebhook)
	v1.Patch("/webhooks/edit", editWebhook)

	v1.Post("/webhooks/delete", deleteWebhook)
	v1.Delete("/webhooks/delete", deleteWebhook)
}

type webhookEditRequest struct {
	ID              string   `json:"id" bson:"_id" form:"id" validate:"required"`
	Destination     string   `json:"destination" bson:"destination" form:"destination" validate:"required,url"`
	Name            string   `json:"name" bson:"name" form:"name" validate:"required"`
	Transformations []string `json:"transformations" bson:"transformations" form:"transformations" validate:"required,min=1"`
	Status          string   `json:"status" bson:"status" form:"status" validate:"required,oneof=active unavailable"`
	Website         bool     `form:"web"`
}

func editWebhook(c *fiber.Ctx) error {
	store, _ := sessions.Get(c)
	usersCollection := client.Database("data").Collection("users")
	var user userPack.User
	err := usersCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "User not found! Are they registered?")
		}
		log.Println(err)
	}

	var body webhookEditRequest
	if err := c.BodyParser(&body); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	errors := validateStruct(body)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	updatedWebhook := userPack.Webhook{
		Destination:     body.Destination,
		Name:            body.Name,
		Transformations: body.Transformations,
		ID:              body.ID,
		Type:            "basic",
		Method:          "post",
		Status:          body.Status,
	}
	user.Webhooks[updatedWebhook.ID] = updatedWebhook
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "webhooks", Value: user.Webhooks}}}}

	usersCollection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}, update)

	if body.Website == true {
		return c.Redirect("/webhooks/edit/" + updatedWebhook.ID)
	}

	return c.JSON(&fiber.Map{
		"success": true,
		"webhook": updatedWebhook,
	})
}

type webhookDeleteRequest struct {
	ID      string `json:"id" bson:"_id" form:"id" validate:"required"`
	Website bool   `form:"web"`
}

func deleteWebhook(c *fiber.Ctx) error {
	var body webhookDeleteRequest
	if err := c.BodyParser(&body); err != nil {
		log.Println(err)
		return err
	}

	if body.ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "A webhook was not given in the request!",
		})
	}

	id := body.ID

	store, _ := sessions.Get(c)
	usersCollection := client.Database("data").Collection("users")
	var user userPack.User
	err := usersCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "User not found! Are they registered?")
		}
		log.Println(err)
	}

	_, ok := user.Webhooks[id]
	if ok == false {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No webhook found with that ID!",
		})
	}

	delete(user.Webhooks, id)

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "webhooks", Value: user.Webhooks}}}}

	usersCollection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: store.Get("user")}}, update)

	if body.Website == true {
		return c.Redirect("/home")
	}

	return c.JSON(&fiber.Map{
		"success": true,
	})
}
