package configuration

import (
	"github.com/gofiber/fiber"
	"github.com/spf13/viper"
	"strings"

	"github.com/gofiber/template/ace"
	"github.com/gofiber/template/amber"
	"github.com/gofiber/template/django"
	"github.com/gofiber/template/handlebars"
	"github.com/gofiber/template/html"
	"github.com/gofiber/template/jet"
	"github.com/gofiber/template/mustache"
	"github.com/gofiber/template/pug"
)

func loadTemplateConfiguration() (enabled bool, views fiber.Views) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("template")
	provider.AddConfigPath("./config")

	// Set default configurations
	setDefaultTemplateConfiguration(provider)

	// Read configuration file
	_ = provider.ReadInConfig()

	// Go over the provided configuration
	switch strings.ToLower(provider.GetString("Engine")) {
		case "ace":
			views = ace.New(provider.GetString("Folder"), provider.GetString("Extension"))
		case "amber":
			views = amber.New(provider.GetString("Folder"), provider.GetString("Extension"))
		case "django":
			views = django.New(provider.GetString("Folder"), provider.GetString("Extension"))
		case "handlebars":
			views = handlebars.New(provider.GetString("Folder"), provider.GetString("Extension"))
		case "jet":
			views = jet.New(provider.GetString("Folder"), provider.GetString("Extension"))
		case "mustache":
			views = mustache.New(provider.GetString("Folder"), provider.GetString("Extension"))
		case "pug":
			views = pug.New(provider.GetString("Folder"), provider.GetString("Extension"))
		default:
			views = html.New(provider.GetString("Folder"), provider.GetString("Extension"))
	}

	// Return the configuration
	return provider.GetBool("Enabled"), views
}

// Set default configuration for the Template Middleware
func setDefaultTemplateConfiguration(provider *viper.Viper) {
	provider.SetDefault("Enabled", false)
	provider.SetDefault("Engine", nil)
	provider.SetDefault("Folder", "./views")
	provider.SetDefault("Extension", ".html")
}
