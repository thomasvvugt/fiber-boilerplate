package web

import (
	"fiber-boilerplate/database"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/session/v2"
	hashing "github.com/thomasvvugt/fiber-hashing"
	"log"
	"strings"
)

func IsAuthenticated(session *session.Session, ctx *fiber.Ctx) (authenticated bool) {
	store := session.Get(ctx)
	// Get User ID from session store
	userID, correct := store.Get("userid").(int64)
	if !correct {
		userID = 0
	}
	auth := false
	if userID > 0 {
		auth = true
	}
	return auth
}

func ShowLoginForm() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Render("login", fiber.Map{})
		if err != nil {
			if err2 := ctx.Status(500).SendString(err.Error()); err2 != nil {
				panic(err2.Error())
			}
		}
		return err
	}
}

func PostLoginForm(hasher hashing.Driver, session *session.Session, db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		username := ctx.FormValue("username")
		// Find user
		user, err := FindUserByUsername(db, username)
		if err != nil {
			log.Fatalf("Error when finding user: %v", err)
		}

		// Check if password matches hash
		if hasher != nil {
			password := ctx.FormValue("password")
			match, err := hasher.MatchHash(password, user.Password)
			if err != nil {
				log.Fatalf("Error when matching hash for password: %v", err)
			}
			if match {
				store := session.Get(ctx)
				defer store.Save()
				// Set the user ID in the session store
				store.Set("userid", user.ID)
				fmt.Printf("User set in session store with ID: %v\n", user.ID)
				if err := ctx.SendString("You should be logged in successfully!"); err != nil {
					panic(err.Error())
				}
			} else {
				if err := ctx.SendString("The entered details do not match our records."); err != nil {
					panic(err.Error())
				}
			}
		} else {
			panic("Hash provider was not set")
		}
		return nil
	}
}

func PostLogoutForm(sessionLookup string, session *session.Session, db *database.Database) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if IsAuthenticated(session, ctx) {
			store := session.Get(ctx)
			store.Delete("userid")
			if err := store.Save(); err != nil {
				panic(err.Error())
			}
			// Check if cookie needs to be unset
			split := strings.Split(sessionLookup, ":")
			if strings.ToLower(split[0]) == "cookie" {
				// Unset cookie on client-side
				ctx.Set("Set-Cookie", split[1] + "=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/; HttpOnly")
				if err := ctx.SendString("You are now logged out."); err != nil {
					panic(err.Error())
				}
				return nil
			}
			return nil
		}
		// TODO: Redirect?
		return nil
	}
}
