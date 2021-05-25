package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lucashmorais/go-interscity/controllers"
)

func SubscriptionRoutes(app *fiber.App) {
	components := app.Group("/subscriptions")

	components.Post("/", controllers.CreateSubscription)
	components.Put("/", controllers.UpdateSubscription)

	components.Get("/:id", controllers.GetSubscription)
	components.Get("/", controllers.GetSubscriptionByResourceUUID)

	components.Delete("/:uuid", controllers.DeleteSubscription)
}
