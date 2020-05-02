package controllers

import "github.com/gofiber/fiber"

func Index(c *fiber.Ctx) {
	if err := c.Render("index", fiber.Map{"question": true}); err != nil {
		c.Status(500).Send(err.Error())
	}
}
