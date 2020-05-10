package routes

import (
	Controller "github.com/thomasvvugt/fiber-boilerplate/app/controllers/web"

	"github.com/gofiber/fiber"
)

func RegisterWeb(app *fiber.App) {
	// Register routes here!

	// Homepage
	app.Get("/", Controller.Index)

	// Panic test route, this brings up an error
	app.Get("/panic", func(c *fiber.Ctx) {
		panic("Hi, I'm a panic error!")
	})
}
