package middleware

import "github.com/gofiber/fiber/v2"

// ForceHTTPS forces HTTPS protocol if not forwarded using a reverse proxy
func ForceHTTPS() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Get("X-Forwarded-Proto") != "https" && ctx.Protocol() == "http" {
			return ctx.Redirect("https://"+ctx.Hostname()+ctx.OriginalURL(), 308)
		}
		return ctx.Next()
	}
}
