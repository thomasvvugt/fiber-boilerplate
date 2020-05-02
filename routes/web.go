package routes

import (
	"github.com/gofiber/fiber"

	"github.com/thomasvvugt/fiber-boilerplate/controllers"
)

func Register(app *fiber.App) {
	// Register routes here!

	// Homepage
	app.Get("/", controllers.Index)

	// TODO: Fix panic route, this brings up an error
	app.Get("/panic", func(c *fiber.Ctx) {
		panic("Hi, I'm a panic error!")
	})
}
