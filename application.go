package main

import (
	"github.com/thomasvvugt/fiber-boilerplate/config"
	"github.com/thomasvvugt/fiber-boilerplate/database"
	"github.com/thomasvvugt/fiber-boilerplate/models"
	"github.com/thomasvvugt/fiber-boilerplate/routes"

	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gofiber/compression"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"
	"github.com/gofiber/template"

	"github.com/spf13/viper"
)

func main() {
	// Load in configurations using spf13/viper
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	var config configuration.Configuration

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			// Maybe wrong file permissions or charset error
			log.Fatalf("unable to load configuration file, %v", err)
		}
	}

	// Set default configurations
	setDefaultConfig()

	// Config file found and successfully parsed
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	// Parse configuration
	config = parseConfig(config)

	// Initialize Fiber
	fiberConfig := convertFiberConfig(config.Fiber)
	app := fiber.New(&fiberConfig)

	// Use logger if enabled
	if config.Logger.Enabled {
		loggerConfig := convertLoggerConfig(config.Logger)
		app.Use(logger.New(loggerConfig))
	}

	// Use Panic Recover (Internal Server Errors) if enabled
	if config.Recover.Enabled {
		recoverConfig := convertRecoverConfig(config.Recover)
		app.Use(recover.New(recoverConfig))
	}

	// Use HTTP best practices
	app.Use(func(c *fiber.Ctx) {
		// Suppress the `www.` at the beginning of URLs
		if config.App.SuppressWWW {
			suppressWWW(c)
		}

		// Force HTTPS protocol
		if config.App.ForceHTTPS {
			forceHTTPS(c)
		}

		// Move on the the next route
		c.Next()
	})

	// Use compression if enabled
	if config.Compression.Enabled {
		compressionConfig := convertCompressionConfig(config.Compression)
		app.Use(compression.New(compressionConfig))
	}

	// Use Cross-Origin Resource Sharing (CORS) if enabled
	if config.CORS.Enabled {
		corsConfig := convertCORSConfig(config.CORS)
		app.Use(cors.New(corsConfig))
	}

	// Set Helmet security settings
	if config.Helmet.Enabled {
		helmetConfig := convertHelmetConfig(config.Helmet)
		app.Use(helmet.New(helmetConfig))
	}

	// Connect to a database and add models
	database.Connect(config.Database)

	// Run auto migrations
	database.Instance().AutoMigrate(&models.User{})

	// Register application routes before serving static file
	routes.Register(app)

	// Serve static, public files
	if config.Public.Enabled {
		staticConfig := convertPublicConfig(config.Public)
		app.Static("/", "./public", staticConfig)
	}

	// Custom 404-page
	app.Use(func(c *fiber.Ctx) {
		c.SendStatus(404)
		if err := c.Render("errors/404", fiber.Map{}); err != nil {
			c.Status(500).Send(err.Error())
		}
	})

	// Close Database connection on Signal Interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		database.Close()
		fmt.Println("Exiting application")
		os.Exit(1)
	}()

	// Start listening
	err = app.Listen(config.App.ListenAddress)
	if err != nil {
		log.Fatalf("Error when start listening: %v", err)
	}
}

// Set default configurations
func setDefaultConfig() {
	viper.SetDefault("App.ListenAddress", 8080)
	viper.SetDefault("App.SuppressWWW", true)

	viper.SetDefault("Fiber.ServerHeader", "Fiber v" + fiber.Version)
	viper.SetDefault("Fiber.BodyLimit", 4 * 1024 * 1024)
	viper.SetDefault("Fiber.Concurrency", 256 * 1024)
	viper.SetDefault("Fiber.TemplateExtension", "html")

	viper.SetDefault("Logger.Enabled", true)
	viper.SetDefault("Logger.Format", "${time} - ${ip} - ${method} ${path}\t${ua}\n")
	viper.SetDefault("Logger.TimeFormat", "15:04:05")

	viper.SetDefault("Recover.Enabled", true)
	viper.SetDefault("Recover.Log", true)

	viper.SetDefault("Compression.Enabled", true)

	viper.SetDefault("CORS.Enabled", true)
	viper.SetDefault("CORS.AllowOrigins", "*")
	viper.SetDefault("CORS.AllowMethods", "GET,POST,HEAD,PUT,DELETE,PATCH")

	viper.SetDefault("Helmet.Enabled", true)
	viper.SetDefault("Helmet.XSSProtection", "1; mode=block")
	viper.SetDefault("Helmet.ContentTypeNosniff", "nosniff")
	viper.SetDefault("Helmet.XFrameOptions", "SAMEORIGIN")

	viper.SetDefault("Public.Enabled", true)
	viper.SetDefault("Public.Index", "index.html")
}

// Parse the configuration
func parseConfig(config configuration.Configuration) configuration.Configuration {
	// Set custom HTTP Server header
	if config.Fiber.ServerHeader == "fiber" {
		config.Fiber.ServerHeader = "Fiber v" + fiber.Version
	}

	// Set custom timeouts
	if config.Fiber.ReadTimeoutSeconds == 0 {
		config.Fiber.ReadTimeout = time.Duration(0)
	} else {
		config.Fiber.ReadTimeout = time.Second * time.Duration(config.Fiber.ReadTimeoutSeconds)
	}
	if config.Fiber.WriteTimeoutSeconds == 0 {
		config.Fiber.WriteTimeout = time.Duration(0)
	} else {
		config.Fiber.WriteTimeout = time.Second * time.Duration(config.Fiber.WriteTimeoutSeconds)
	}
	if config.Fiber.IdleTimeoutSeconds == 0 {
		config.Fiber.IdleTimeout = time.Duration(0)
	} else {
		config.Fiber.IdleTimeout = time.Second * time.Duration(config.Fiber.IdleTimeoutSeconds)
	}

	// Set templating engine
	switch strings.ToLower(config.Fiber.TemplateEngineName) {
	case "mustache":
		config.Fiber.TemplateEngine = template.Mustache()
	case "amber":
		config.Fiber.TemplateEngine = template.Amber()
	case "handlebars":
		config.Fiber.TemplateEngine = template.Handlebars()
	case "pug":
		config.Fiber.TemplateEngine = template.Pug()
	default:
		config.Fiber.TemplateEngine = nil
	}

	// Set Logger format presets
	switch strings.ToLower(config.Logger.TimeFormat) {
	case "ansic":
		config.Logger.TimeFormat = time.ANSIC
	case "unixdate":
		config.Logger.TimeFormat = time.UnixDate
	case "rubydate":
		config.Logger.TimeFormat = time.RubyDate
	case "rfc822":
		config.Logger.TimeFormat = time.RFC822
	case "rfc822z":
		config.Logger.TimeFormat = time.RFC822Z
	case "rfc850":
		config.Logger.TimeFormat = time.RFC850
	case "rfc1123":
		config.Logger.TimeFormat = time.RFC1123
	case "rfc1123z":
		config.Logger.TimeFormat = time.RFC1123Z
	case "rfc3339":
		config.Logger.TimeFormat = time.RFC3339
	case "rfc3339nano":
		config.Logger.TimeFormat = time.RFC3339Nano
	case "kitchen":
		config.Logger.TimeFormat = time.Kitchen
	case "stamp":
		config.Logger.TimeFormat = time.Stamp
	case "stampmilli":
		config.Logger.TimeFormat = time.StampMilli
	case "stampmicro":
		config.Logger.TimeFormat = time.StampMicro
	case "stampnano":
		config.Logger.TimeFormat = time.StampNano
	}

	return config
}

// Convert Fiber configuration to fiber.Config
func convertFiberConfig(config configuration.FiberConfiguration) fiber.Settings {
	return fiber.Settings{
		Prefork:                   config.Prefork,
		StrictRouting:             config.StrictRouting,
		CaseSensitive:             config.CaseSensitive,
		ServerHeader:              config.ServerHeader,
		Immutable:                 config.Immutable,
		ETag:                      config.ETag,
		BodyLimit:                 config.BodyLimit,
		Concurrency:               config.Concurrency,
		DisableKeepalive:          config.DisableKeepalive,
		DisableDefaultDate:        config.DisableDefaultDate,
		DisableDefaultContentType: config.DisableDefaultContentType,
		DisableStartupMessage:     config.DisableStartupMessage,
		TemplateFolder:            config.TemplateFolder,
		TemplateEngine:            config.TemplateEngine,
		TemplateExtension:         config.TemplateExtension,
		ReadTimeout:               config.ReadTimeout,
		WriteTimeout:              config.WriteTimeout,
		IdleTimeout:               config.IdleTimeout,
	}
}

// Convert Logger configuration to logger.Config
func convertLoggerConfig(config configuration.LoggerConfiguration) logger.Config {
	return logger.Config{
		Format:     config.Format,
		TimeFormat: config.TimeFormat,
	}
}

// Convert Recover configuration to recover.Config
func convertRecoverConfig(config configuration.RecoverConfiguration) recover.Config {
	return recover.Config{
		Handler: func(c *fiber.Ctx, err error) {
			c.SendStatus(500)
			data := fiber.Map{"error" : err.Error()}
			if err := c.Render("errors/500", data); err != nil {
				c.Status(500).Send(err.Error())
			}
		},
		Log:     config.Log,
	}
}

// Suppress `www.` at the beginning of URLs
func suppressWWW(c *fiber.Ctx) {
	hostnameSplit := strings.Split(c.Hostname(), ".")
	if hostnameSplit[0] == "www" && len(hostnameSplit) > 1 {
		newHostname := ""
		for i := 1; i <= (len(hostnameSplit) - 1); i++ {
			if i != (len(hostnameSplit) - 1) {
				newHostname = newHostname + hostnameSplit[i] + "."
			} else {
				newHostname = newHostname + hostnameSplit[i]
			}
		}
		c.Redirect(c.Protocol() + "://" + newHostname + c.OriginalURL(), 301)
	}
}

// Force the use of HTTPS
func forceHTTPS(c *fiber.Ctx) {
	if c.Protocol() == "http" {
		c.Redirect("https://" + c.Hostname() + c.OriginalURL(), 308)
	}
}

// Convert Compression configuration to compression.Config
func convertCompressionConfig(config configuration.CompressionConfiguration) compression.Config {
	return compression.Config{
		Level:  config.Level,
	}
}

// Convert CORS configuration to cors.Config
func convertCORSConfig(config configuration.CORSConfiguration) cors.Config {
	return cors.Config{
		AllowOrigins:     strings.Split(config.AllowOrigins, ","),
		AllowMethods:     strings.Split(config.AllowMethods, ","),
		AllowHeaders:     strings.Split(config.AllowHeaders, ","),
		AllowCredentials: config.AllowCredentials,
		ExposeHeaders:    strings.Split(config.ExposeHeaders, ","),
		MaxAge:           config.MaxAge,
	}
}

// Convert Helmet configuration to helmet.Config
func convertHelmetConfig(config configuration.HelmetConfiguration) helmet.Config {
	return helmet.Config{
		XSSProtection:         config.XSSProtection,
		ContentTypeNosniff:    config.ContentTypeNosniff,
		XFrameOptions:         config.XFrameOptions,
		HSTSMaxAge:            config.HSTSMaxAge,
		HSTSExcludeSubdomains: config.HSTSExcludeSubdomains,
		ContentSecurityPolicy: config.ContentSecurityPolicy,
		CSPReportOnly:         config.CSPReportOnly,
		HSTSPreloadEnabled:    config.HSTSPreloadEnabled,
		ReferrerPolicy:        config.ReferrerPolicy,
	}
}

// Convert Public configuration to fiber.Static
func convertPublicConfig(config configuration.PublicConfiguration) fiber.Static {
	return fiber.Static{
		Compress:  config.Compress,
		ByteRange: config.ByteRange,
		Browse:    config.Browse,
		Index:     config.Index,
	}
}
