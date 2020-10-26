package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/session/v2"

	"go-fiber-v2-boilerplate/app/configuration"
	"go-fiber-v2-boilerplate/app/models"
	"go-fiber-v2-boilerplate/app/providers"
	"go-fiber-v2-boilerplate/database"
	"go-fiber-v2-boilerplate/routes"
)

func main() {
	// Load configurations
	config, err := configuration.LoadConfigurations()
	if err != nil {
		// Error when loading the configurations
		log.Fatalf("An error occurred while loading the configurations: %v", err)
	}

	// Create a new Fiber application
	app := fiber.New(config.Fiber)

	// Use the Logger Middleware if enabled
	if config.Enabled["logger"] {
		app.Use(logger.New(config.Logger))
	}

	// Use the Recover Middleware if enabled
	if config.Enabled["recover"] {
		app.Use(recover.New(config.Recover))
	}

	// Use HTTP best practices
	app.Use(func(c *fiber.Ctx) error {
		// Suppress the `www.` at the beginning of URLs
		if config.App.SuppressWWW {
			providers.SuppressWWW(c)
		}
		// Force HTTPS protocol
		if config.App.ForceHTTPS {
			providers.ForceHTTPS(c)
		}
		// Move on the the next route
		return c.Next()
	})

	// Use the Compression Middleware if enabled
	if config.Enabled["compression"] {
		app.Use(compress.New(config.Compression))
	}

	// Use the CORS Middleware if enabled
	if config.Enabled["cors"] {
		app.Use(cors.New(config.CORS))
	}

	// Use the Helmet Middleware if enabled
	if config.Enabled["helmet"] {
		app.Use(helmet.New(config.Helmet))
	}

	// Use the Session Middleware if enabled
	if config.Enabled["session"] {
		// create session handler
		providers.SetSessionProvider(session.New(config.Session))
	}

	// Set hashing provider
	if config.Enabled["hash"] {
		providers.SetHashProvider(config.Hash)
	}

	// Connect to a database
	if config.Enabled["database"] {
		database.Connect(&config.Database)
	}

	// Run auto migrations
	database.Instance().AutoMigrate(&models.Role{})
	database.Instance().AutoMigrate(&models.User{})
	// Set CASCADE foreign key
	database.Instance().Model(&models.User{}).AddForeignKey("role_id", "roles(id)", "RESTRICT", "CASCADE")

	// Register application web routes
	routes.RegisterWeb(app)

	// Register application API routes (using the /api/v1 group)
	api := app.Group("/api").(*fiber.Group)
	apiv1 := api.Group("/v1").(*fiber.Group)
	routes.RegisterAPI(apiv1)

	// Serve public, static files
	if config.Enabled["public"] {
		app.Static(config.PublicPrefix, config.PublicRoot, config.Public)
	}

	// Custom 404-page
	app.Use(func(c *fiber.Ctx) error {

		if err := c.Render("errors/404", fiber.Map{}); err != nil {
			c.Status(500).SendString(err.Error())
		}
		return c.Status(fiber.StatusNotFound).SendString("Sorry can't find that!")
	})

	// Set configuration provider
	providers.SetConfiguration(&config)

	// Close any connections on interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		exit(&config, app, nil)
	}()

	// Start listening on the specified address
	err = app.Listen(fmt.Sprintf(":%v",config.App.Listen))
	if err != nil {
		// Exit the application
		exit(&config, app, err)
	}
}

func exit(config *configuration.Configuration, app *fiber.App, err error) {
	// Close database connection
	var dbErr error
	if config.Enabled["database"] {
		dbErr = database.Close()
		if dbErr != nil {
			fmt.Printf("Closed database: %v\n", dbErr)
		} else {
			fmt.Println("Closed database.")
		}
	}
	// Shutdown Fiber application
	var appErr error
	if err != nil {
		fmt.Printf("Shutdown Fiber application: %v", err)
		appErr = err
	} else {
		appErr = app.Shutdown()
		if appErr != nil {
			fmt.Printf("Shutdown Fiber application: %v", appErr)
		} else {
			fmt.Print("Shutdown Fiber application.")
		}
	}
	// Return with corresponding exit code
	if dbErr != nil || appErr != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
