package routes

import (
	Controller "fiber-boilerplate/app/controllers/web"
	"fiber-boilerplate/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/session/v2"
	hashing "github.com/thomasvvugt/fiber-hashing"
	"log"
)

func RegisterWeb(web fiber.Router, session *session.Session, sessionLookup string, db *database.Database, hasher hashing.Driver) {
	// Homepage
	web.Get("/", Controller.Index(session, db))

	// Panic test route, this brings up an error
	web.Get("/panic", func(ctx *fiber.Ctx) error {
		panic("Hi, I'm a panic error!")
	})

	// Test to load static, compiled assets
	web.Get("/test", func(c *fiber.Ctx) error {
		return c.Render("test", fiber.Map{})
	})

	// Make a new hash
	web.Get("/hash/*", func(ctx *fiber.Ctx) error {
		hash, err := hasher.CreateHash(ctx.Params("*"))
		if err != nil {
			log.Fatalf("Error when creating hash: %v", err)
		}
		if err := ctx.SendString(hash); err != nil {
			panic(err.Error())
		}
		return err
	})

	// Auth routes
	web.Get("/login", Controller.ShowLoginForm())
	web.Post("/login", Controller.PostLoginForm(hasher, session, db))
	web.Post("/logout", Controller.PostLogoutForm(sessionLookup, session, db))
}
