package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lucashmorais/go-interscity/controllers"
)

// UserRoutes | @desc: "/api/v1/users" api route
/*
func UserRoutes(app *fiber.App) {
	users := app.Group("/api/v1/users")
	users.Get("/", controllers.GetUsers)
	users.Get("/:id", controllers.GetUser)
	users.Post("/", controllers.CreateUser)
	users.Patch("/:id", controllers.UpdateUser)
	users.Delete("/:id", controllers.DeleteUser)
}
*/

func ResourceRoutes(app *fiber.App) {
	components := app.Group("/resources")
	components.Get("/", controllers.GetResources)
	components.Get("/:uuid", controllers.GetResource)
	components.Post("/", controllers.CreateResource)
	components.Put("/:uuid", controllers.UpdateResource)
	components.Delete("/:uuid", controllers.DeleteResource)

	components.Post("/:uuid/data", controllers.PostData)
}
