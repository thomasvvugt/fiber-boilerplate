package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"

	"go-fiber-v2-boilerplate/app/providers"
)

func ShowLoginForm(c *fiber.Ctx) error {
	if err := c.Render("login", fiber.Map{}); err != nil {
		c.Status(500).SendString(err.Error())
		return err
	}
	return nil
}

func PostLoginForm(c *fiber.Ctx) error {
	username := c.FormValue("username")
	// Find user
	user, err := FindUserByUsername(username)
	if err != nil {
		log.Fatalf("Error when finding user: %v", err)
	}
	// Check if password matches hash
	if providers.HashProvider() != nil {
		password := c.FormValue("password")
		match, err := providers.HashProvider().MatchHash(password, user.Password)
		if err != nil {
			log.Fatalf("Error when matching hash for password: %v", err)
		}
		if match {
			store := providers.SessionProvider().Get(c)
			defer store.Save()
			// Set the user ID in the session store
			store.Set("userid", user.ID)
			fmt.Printf("User set in session store with ID: %v\n", user.ID)
			c.SendString("You should be logged in successfully!")
		} else {
			c.SendString("The entered details do not match our records.")
		}
	} else {
		panic("Hash provider was not set")
	}
	return nil
}

func PostLogoutForm(c *fiber.Ctx) error {
	if providers.IsAuthenticated(c) {
		store := providers.SessionProvider().Get(c)
		store.Delete("userid")
		store.Save()
		// Check if cookie needs to be unset
		config := providers.GetConfiguration()
		lookup := config.Session.Lookup
		split := strings.Split(lookup, ":")
		if strings.ToLower(split[0]) == "cookie" {
			// Unset cookie on client-side
			c.Set("Set-Cookie", split[1] + "=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; HttpOnly")
			c.SendString("You are now logged out.")
			return nil
		}
	}
	return nil
}
