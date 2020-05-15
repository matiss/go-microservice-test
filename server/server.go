package server

import (
	"context"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/limiter"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"

	"github.com/matiss/go-microservice-test/services"
)

func Run(config *services.ConfigService, currencyService *services.CurrencyService) {
	ctx := context.Background()

	s := fiber.New()

	// Set prefork
	s.Settings.Prefork = false

	// Disable startup message
	s.Settings.DisableStartupMessage = true

	// Recover middleware
	s.Use(recover.New(recover.Config{
		// Config is optional
		Handler: func(c *fiber.Ctx, err error) {
			c.SendString(err.Error())
			c.SendStatus(500)
		},
	}))

	// Create a rate limiter struct.
	rateLimiter := limiter.Config{
		Timeout: 1,
		Max:     config.HTTP.RateLimit,
	}
	s.Use(limiter.New(rateLimiter))

	// CORS middleware
	corsConfig := cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"content-type", "Authorization"},
	}
	s.Use(cors.New(corsConfig))

	// Logger middleware
	s.Use(logger.New())

	currencyHandler := NewCurrencyHandler(ctx, currencyService)

	// Currency Latest handler
	s.Get("/currency/latest", currencyHandler.Latest)

	// Currency by symbol handler
	s.Get("/currency/:id", currencyHandler.BySymbol)

	// Root handler
	s.Get("/", RootHandler)

	// Handle robots.txt file
	s.Get("/robots.txt", RobotsTXTHandler)

	// Start server
	s.Listen(config.HTTP.Address)
}
