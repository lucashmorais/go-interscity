package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lucashmorais/go-interscity/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSubscriptions(c *fiber.Ctx) error {
	return c.SendString("Endpoint is working: " + c.OriginalURL())
}

func GetSubscription(c *fiber.Ctx) error {
	idParam := c.Params("id")
	subscriptionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": idParam + " is not a valid id!"})
	}

	filer := bson.D{{Key: "_id", Value: subscriptionID}}
	subscriptionRecord := models.SubscriptionCollection.FindOne(c.Context(), filer)
	if subscriptionRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": "No subscription with id: " + idParam + " was found!"})
	}

	subscription := &models.Subscription{}
	subscriptionRecord.Decode(subscription)

	return c.JSON(fiber.Map{"success": true, "data": subscription})
	// return c.SendString("Endpoint is working: " + c.OriginalURL())
}

func GetSubscriptionByResourceUUID(c *fiber.Ctx) error {
	uuidParam := c.Query("uuid")

	filer := bson.D{{Key: "subscription.uuid", Value: uuidParam}}
	subscriptionRecord := models.SubscriptionCollection.FindOne(c.Context(), filer)
	if subscriptionRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": "No subscription was found whose resource had uuid: " + uuidParam})
	}

	subscription := &models.Subscription{}
	subscriptionRecord.Decode(subscription)

	return c.JSON(fiber.Map{"success": true, "data": subscription})
	// return c.SendString("Endpoint is working: " + c.OriginalURL())
}

func CreateSubscription(c *fiber.Ctx) error {
	// return c.SendString("Endpoint is working: " + c.OriginalURL())

	resultString := ""
	uuid := uuid.NewString()

	resultString += uuid

	subscription := new(models.Subscription)
	if err := c.BodyParser(subscription); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}

	insertionResult, err := models.SubscriptionCollection.InsertOne(c.Context(), subscription)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}

	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := models.SubscriptionCollection.FindOne(c.Context(), filter)

	createdSubscription := &models.Subscription{}
	createdRecord.Decode(createdSubscription)

	return c.JSON(fiber.Map{"id": insertionResult.InsertedID, "success": true, "data": createdSubscription})
}

func UpdateSubscription(c *fiber.Ctx) error {
	return c.SendString("Endpoint is working: " + c.OriginalURL())
}

func DeleteSubscription(c *fiber.Ctx) error {
	return c.SendString("Endpoint is working: " + c.OriginalURL())
}
