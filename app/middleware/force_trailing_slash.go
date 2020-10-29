package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"regexp"
	"strings"
)

// ForceTrailingSlash will add a trailing slash (`/`) if this is not present is the client's request.
// This also takes the ability for the client to request a file extension into account.
func ForceTrailingSlash() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		originalUrl := utils.ImmutableString(ctx.OriginalURL())

		// Check if the client is requesting a file extension
		extMatch, _ := regexp.MatchString("\\.[a-zA-Z0-9]+$", originalUrl)

		if !strings.HasSuffix(originalUrl, "/") && !extMatch {
			return ctx.Redirect(originalUrl + "/")
		}
		return ctx.Next()
	}
}