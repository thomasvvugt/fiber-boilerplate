package web

import (
	"github.com/gofiber/fiber/v2"
	"log"

	"go-fiber-v2-boilerplate/app/providers"
)

func Index(c *fiber.Ctx) error {
	auth := providers.IsAuthenticated(c)
	// Bind data to template
	bind := fiber.Map{
		"name": "Fiber",
		"auth": auth,
	}
	if auth {
		store := providers.SessionProvider().Get(c)
		// Get User ID from session store
		userID, _ := store.Get("userid").(int64)
		user, err := FindUserByID(userID)
		if err != nil {
			log.Fatalf("Error when finding user by ID: %v", err)
			return err
		}
		bind["username"] = user.Name
	}
	// Render template
	if err := c.Render("index", bind); err != nil {
		c.Status(500).SendString(err.Error())
		return err
	}
	return nil
}
