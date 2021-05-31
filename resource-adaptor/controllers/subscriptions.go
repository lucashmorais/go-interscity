package controllers

import (
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/lucashmorais/go-interscity/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSubscriptions(c *fiber.Ctx) error {
	filter := bson.D{{}}
	cursor, err := models.SubscriptionCollection.Find(c.Context(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err})
	}

	var users []models.Subscription = make([]models.Subscription, 0)

	if err := cursor.All(c.Context(), &users); err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err})
	}

	return c.JSON(fiber.Map{"success": true, "data": users})
	// return c.SendString("Endpoint is working: " + c.OriginalURL())
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
}

func CreateSubscription(c *fiber.Ctx) error {
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

func coreMergeStructs(source reflect.Value, target reflect.Value) {
	for i := 0; i < source.NumField(); i++ {
		target.FieldByName(source.Type().Field(i).Name).Set(source.Field(i))
	}
}

func UpdateSubscription(c *fiber.Ctx) error {
	idParam := c.Params("id")
	subscriptionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": idParam + " is not a valid id!"})
	}

	// subscription := new(models.Subscription)
	var subscription models.Subscription
	if err := c.BodyParser(&subscription); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}

	filter := bson.D{{Key: "_id", Value: subscriptionID}}

	subscriptionRecord := models.SubscriptionCollection.FindOne(c.Context(), filter)
	if subscriptionRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": "No subscription with id: " + idParam + " was found!"})
	}

	subscriptionRecordStruct := &models.Subscription{}
	subscriptionRecord.Decode(subscriptionRecordStruct)

	coreMergeStructs(reflect.ValueOf(subscription), reflect.Indirect(reflect.ValueOf(subscriptionRecordStruct)))

	finalSubscriptionRecord := models.SubscriptionCollection.FindOneAndReplace(c.Context(), filter, subscriptionRecordStruct)

	if finalSubscriptionRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": "No subscription with id: " + idParam + " was found!"})
	}

	updatedSubscription := &models.Subscription{}
	finalSubscriptionRecord.Decode(updatedSubscription)

	return c.JSON(fiber.Map{"success": true, "data": updatedSubscription})
	// return c.SendString("Endpoint is working: " + c.OriginalURL())
}

func DeleteSubscription(c *fiber.Ctx) error {
	idParam := c.Params("id")
	subscriptionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": idParam + " is not a valid id!"})
	}

	filter := bson.D{{Key: "_id", Value: subscriptionID}}
	subscriptionRecord := models.SubscriptionCollection.FindOneAndDelete(c.Context(), filter)
	if subscriptionRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": "No subscription with id: " + idParam + " was found!"})
	}

	return c.JSON(fiber.Map{"success": true, "data": "Subscription was deleted!"})
	// return c.SendString("Endpoint is working: " + c.OriginalURL())
}
