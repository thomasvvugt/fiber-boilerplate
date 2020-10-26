package routes

import (
	"github.com/gofiber/fiber/v2"
	"log"

	Controller "go-fiber-v2-boilerplate/app/controllers/web"
	"go-fiber-v2-boilerplate/app/providers"
)

func RegisterWeb(app *fiber.App) {
	// Homepage
	app.Get("/", Controller.Index)

	// Panic test route, this brings up an error
	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("Hi, I'm a panic error!")
	})

	// Make a new hash
	app.Get("/hash/*", func(c *fiber.Ctx) error {
		hash, err := providers.HashProvider().CreateHash(c.Params("*"))
		if err != nil {
			log.Fatalf("Error when creating hash: %v", err)
		}
		c.SendString(hash)
		return nil
	})

	// Auth routes
	app.Get("/login", Controller.ShowLoginForm)
	app.Post("/login", Controller.PostLoginForm)
	app.Post("/logout", Controller.PostLogoutForm)
}
