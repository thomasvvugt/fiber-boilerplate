package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"strings"
)

// SuppressWWW suppresses the `www.` at the beginning of URLs
func SuppressWWW() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		hostname := utils.ImmutableString(ctx.Hostname())
		hostnameSplit := strings.Split(hostname, ".")

		if hostnameSplit[0] == "www" && len(hostnameSplit) > 1 {
			newHostname := ""
			for i := 1; i <= (len(hostnameSplit) - 1); i++ {
				if i != (len(hostnameSplit) - 1) {
					newHostname = newHostname + hostnameSplit[i] + "."
				} else {
					newHostname = newHostname + hostnameSplit[i]
				}
			}
			return ctx.Redirect(ctx.Protocol()+"://"+newHostname+ctx.OriginalURL(), 301)
		}
		return ctx.Next()
	}
}
