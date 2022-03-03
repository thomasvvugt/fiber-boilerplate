package web

import (
	"fiber-boilerplate/database"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func Index(session *session.Store, db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth := IsAuthenticated(session, ctx)

		// Bind data to template
		bind := fiber.Map{
			"name": "Fiber",
			"auth": auth,
		}

		if auth {
			store, err := session.Get(ctx)
			if err != nil {
				panic(err)
			}
			// Get User ID from session store
			userID := store.Get("userid").(int64)
			user, err := FindUserByID(db, userID)
			if err != nil {
				log.Fatalf("Error when finding user by ID: %v", err)
			}
			bind["username"] = user.Name
		}

		// Render template
		err := ctx.Render("index", bind)
		if err != nil {
			err2 := ctx.Status(500).SendString(err.Error())
			if err2 != nil {
				panic(err2.Error())
			}
		}
		return err
	}
}
