package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lucashmorais/go-interscity/models"
	"go.mongodb.org/mongo-driver/bson"
)

/*
func ResourceRoutes(app *fiber.App) {
	components := app.Group("/components")
	components.Get("/", controllers.GetResources)
	components.Get("/:uuid", controllers.GetResource)
	components.Post("/", controllers.CreateResource)
	components.Patch("/:uuid", controllers.UpdateResource)
	components.Delete("/:uuid", controllers.DeleteResource)

	components.Post("/:uuid/data", controllers.PostData)
}
*/

func GetResources(c *fiber.Ctx) error {
	filter := bson.D{{}}
	cursor, err := models.ResourceCollection.Find(c.Context(), filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err})
	}

	var users []models.Resource = make([]models.Resource, 0)

	if err := cursor.All(c.Context(), &users); err != nil {
		fmt.Println(err)
		return c.Status(500).JSON(fiber.Map{"success": false, "data": err})
	}

	return c.JSON(fiber.Map{"success": true, "data": users})
	// return c.SendString("Endpoint is working: " + c.OriginalURL())
}

func GetResource(c *fiber.Ctx) error {
	idParam := c.Params("uuid")

	filter := bson.D{{Key: "_id", Value: idParam}}
	resourceRecord := models.ResourceCollection.FindOne(c.Context(), filter)
	if resourceRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": "No resource with id: " + idParam + " was found!"})
	}

	resource := &models.Resource{}
	resourceRecord.Decode(resource)

	return c.JSON(fiber.Map{"success": true, "data": resource})
	// return c.SendString("Endpoint is working: " + c.OriginalURL())
}

func CreateResource(c *fiber.Ctx) error {
	// return c.SendString("Endpoint is working: " + c.OriginalURL())

	resultString := ""
	uuid := uuid.NewString()

	resultString += uuid

	resource := new(models.Resource)
	resource.UUID = uuid
	if err := c.BodyParser(resource); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}

	insertionResult, err := models.ResourceCollection.InsertOne(c.Context(), resource)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}

	filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
	createdRecord := models.ResourceCollection.FindOne(c.Context(), filter)

	createdResource := &models.Resource{}
	createdRecord.Decode(createdResource)

	// return c.JSON(fiber.Map{"success": true, "data": createdResource})
	return c.JSON(resource)
}

func UpdateResource(c *fiber.Ctx) error {
	idParam := c.Params("uuid")

	filter := bson.D{{Key: "_id", Value: idParam}}
	resourceRecord := models.ResourceCollection.FindOne(c.Context(), filter)
	if resourceRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": "No resource with id: " + idParam + " was found!"})
	}

	resource := new(models.Resource)
	if err := c.BodyParser(resource); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": err})
	}

	resource.UUID = idParam

	// ONE CANNOT REPLACE AN ENTRY FOR ANOTHER ONE MISSING A PRIMARY KEY
	finalResourceRecord := models.ResourceCollection.FindOneAndReplace(c.Context(), filter, resource)

	if finalResourceRecord.Err() != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "data": "[SECOND] No resource with id: " + idParam + " was found!"})
	}

	updatedResource := &models.Resource{}
	finalResourceRecord.Decode(updatedResource)

	return c.JSON(fiber.Map{"success": true, "data": updatedResource})
	// return c.SendString("Endpoint is working: " + c.OriginalURL())
}

func DeleteResource(c *fiber.Ctx) error {
	return c.SendString("Endpoint is working: " + c.OriginalURL())
}

func PostData(c *fiber.Ctx) error {
	return c.SendString("Endpoint is working: " + c.OriginalURL())
}
