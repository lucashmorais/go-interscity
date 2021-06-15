package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/lucashmorais/go-interscity/database"
	"github.com/lucashmorais/go-interscity/models"
	"github.com/lucashmorais/go-interscity/routes"
)

func PrecisionLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// start timer
		start := time.Now()

		// next routes
		err := c.Next()

		// stop timer
		stop := time.Now()

		// Print request duration
		// fmt.Printf("app;dur=%v\n", stop.Sub(start).String())

		// Print final time

		client_address := c.Context().RemoteAddr().String()
		fmt.Printf("[%s]: %v %v %v \n", client_address, stop.Sub(start).Microseconds(), start.Format(time.RFC3339Nano), stop.Format(time.RFC3339Nano))

		// return stack error if exist
		return err
	}
}

func main() {
	// Load enviromental variables from config
	godotenv.Load("./config/config.env")

	// Create new fiber instance
	app := fiber.New(fiber.Config{
		// Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		// Concurrency:   4 * 1024 * 1024,
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
	// app.Use(logger.New())
	// app.Use(cors.New())
	app.Use(PrecisionLogger())

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
