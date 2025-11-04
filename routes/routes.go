package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rakibulbanna/go-fiber-postgres/controllers"
	"github.com/rakibulbanna/go-fiber-postgres/middleware"
)

func SetupRoutes(app *fiber.App, authController *controllers.AuthController, bookController *controllers.BookController, authMiddleware *middleware.AuthMiddleware) {
	api := app.Group("/api")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/signup", authController.SignUp)
	auth.Post("/login", authController.Login)
	auth.Post("/signin", authController.Login) // Alias for login

	// Public book routes (anyone can view)
	books := api.Group("/books")
	books.Get("/", bookController.GetBooks)
	books.Get("/:id", bookController.GetBook)

	// Protected book routes (authenticated users only)
	protectedBooks := api.Group("/books", authMiddleware.RequireAuth)
	protectedBooks.Post("/", bookController.CreateBook)
	protectedBooks.Put("/:id", bookController.UpdateBook)
	protectedBooks.Delete("/:id", bookController.DeleteBook)
}
