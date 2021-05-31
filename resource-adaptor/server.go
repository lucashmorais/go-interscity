package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/lucashmorais/go-interscity/database"
	"github.com/lucashmorais/go-interscity/models"
	"github.com/lucashmorais/go-interscity/routes"
)

func main() {
	// Load enviromental variables from config
	godotenv.Load("./config/config.env")

	// Create new fiber instance
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		// ServerHeader:  "Fiber",
	})

	// Connect to database
	println("Trying to connect to database now...")
	database.Connect()
	println("Just connected to DB.")
	defer database.Cancel()
	defer database.Client.Disconnect(database.Ctx)

	// Create model schemas
	println("Trying to create schemas...")
	models.CreateUserSchema()
	models.CreateResourceSchema()
	models.CreateSubscriptionSchema()
	models.CreateResourceDataSchema()
	println("Just created the schemas.")

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	routes.UserRoutes(app)
	routes.ResourceRoutes(app)
	routes.SubscriptionRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// messaging.ConnectAndSend()
	// messaging.Receive()

	app.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))
	fmt.Printf("We'll be listening at port %s\n", os.Getenv("PORT"))
}
