package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rakibulbanna/go-fiber-postgres/config"
	"github.com/rakibulbanna/go-fiber-postgres/controllers"
	"github.com/rakibulbanna/go-fiber-postgres/middleware"
	"github.com/rakibulbanna/go-fiber-postgres/models"
	"github.com/rakibulbanna/go-fiber-postgres/repositories"
	"github.com/rakibulbanna/go-fiber-postgres/routes"
	"github.com/rakibulbanna/go-fiber-postgres/services"
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

	// Run migrations
	if err := models.RunMigrations(db); err != nil {
		log.Fatal("Error running migrations: ", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	bookRepo := repositories.NewBookRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	bookService := services.NewBookService(bookRepo, userRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	bookController := controllers.NewBookController(bookService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

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
	routes.SetupRoutes(app, authController, bookController, authMiddleware)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
