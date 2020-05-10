package web

import "github.com/gofiber/fiber"

func Index(c *fiber.Ctx) {
	bind := fiber.Map{
		"name": "Fiber",
	}
	if err := c.Render("index", bind); err != nil {
		c.Status(500).Send(err.Error())
	}
}
