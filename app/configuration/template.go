package configuration

import (
	"strings"

	"github.com/gofiber/template"

	"github.com/spf13/viper"
)

func loadTemplateConfiguration() (enabled bool, engine func(raw string, bind interface{}) (out string, err error)) {
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
		case "mustache":
			engine = template.Mustache()
		case "amber":
			engine = template.Amber()
		case "handlebars":
			engine = template.Handlebars()
		case "pug":
			engine = template.Pug()
		default:
			engine = nil
	}

	// Return the configuration
	return provider.GetBool("Enabled"), engine
}

// Set default configuration for the Template Middleware
func setDefaultTemplateConfiguration(provider *viper.Viper) {
	provider.SetDefault("Enabled", false)
	provider.SetDefault("Engine", nil)
}
