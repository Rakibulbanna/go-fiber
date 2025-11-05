package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rakibulbanna/go-fiber-postgres/config"
	authModule "github.com/rakibulbanna/go-fiber-postgres/internal/modules/auth"
	bookModule "github.com/rakibulbanna/go-fiber-postgres/internal/modules/book"
	"github.com/rakibulbanna/go-fiber-postgres/middleware"
	"github.com/rakibulbanna/go-fiber-postgres/storage"
)

func main() {
	// Load configuration
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env.dev"
	}
	cfg := config.LoadConfig(envFile)

	// Database connection
	dbConfig := storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(dbConfig)
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Initialize modules
	authService := authModule.NewService(db, cfg.JWTSecret)
	authController := authModule.NewController(authService)

	bookService := bookModule.NewService(db)
	bookController := bookModule.NewController(bookService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())

	// Setup routes
	api := app.Group("/api")
	authModule.SetupRoutes(api, authController)
	bookModule.SetupRoutes(api, bookController, authMiddleware)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
