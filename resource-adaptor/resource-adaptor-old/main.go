package main

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/lucashmorais/go_fiber/database"
)

func basicTest(c *fiber.Ctx) error {
	c.Accepts("application/json")
	return c.SendString("Hello, World!")
}

func createResource(c *fiber.Ctx) error {
	c.Accepts("application/json")

	var r database.Resource
	var dbResource database.DBResource
	var byteData []byte

	if err := c.BodyParser(&r); err != nil {
		return err
	}

	byteData, err := json.Marshal(r.Data)
	dbResource.Data = string(byteData)
	dbResource.UUID = uuid.NewString()

	if err != nil {
		return err
	}

	db := database.DBConnResources
	db.Create(&dbResource)
	return c.JSON(dbResource)
}

func updateResource(c *fiber.Ctx) error {
	c.Accepts("application/json")
	return c.SendString("Hello, World!")
}

func postData(c *fiber.Ctx) error {
	c.Accepts("application/json")
	return c.SendString("Hello, World!")
}

func subscribeToResource(c *fiber.Ctx) error {
	c.Accepts("application/json")
	return c.SendString("Hello, World!")
}

func updateSubscription(c *fiber.Ctx) error {
	c.Accepts("application/json")
	return c.SendString("Hello, World!")
}

func getSubscription(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

//TODO: If no query string is found, this should probably fail
func filterSubscriptionsByUUID(c *fiber.Ctx) error {
	return c.SendString(c.Query("uuid"))
}

func deleteSubscription(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func initDatabase() {
	var err1, err2 error

	database.DBConnResources, err1 = gorm.Open("sqlite3", "resources.db")
	database.DBConnSubscriptions, err2 = gorm.Open("sqlite3", "subscriptions.db")

	if err1 != nil || err2 != nil {
		panic("Failed to connect to the Database")
	}

	fmt.Println("Database connection successfully established")

	database.DBConnResources.AutoMigrate(&database.DBResource{})
	database.DBConnSubscriptions.AutoMigrate(&database.Subscription{})
	fmt.Println("DB auto-migration was set up")
}

func setupRouter(app *fiber.App) {
	app.Get("/", basicTest)

	app.Post("/components", createResource)
	app.Put("/components/:uuid", updateResource)
	app.Post("/components/:uuid/data", postData)
	app.Post("/subscriptions", subscribeToResource)
	app.Put("/subscriptions/:uuid", updateSubscription)
	app.Get("/subscriptions/:id", getSubscription)

	// app.Get("/subscriptions\\?uuid=:uuid", filterSubscriptionsByUUID)
	app.Get("/subscriptions", filterSubscriptionsByUUID)

	app.Delete("/subscriptions/:id", deleteSubscription)
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		// StrictRouting: true,
		// ServerHeader:  "Fiber",
	})

	initDatabase()

	// This is only executed at the end of the program,
	// for closing the DB connections
	defer database.DBConnResources.Close()
	defer database.DBConnSubscriptions.Close()

	setupRouter(app)

	app.Listen(":4123")
}
