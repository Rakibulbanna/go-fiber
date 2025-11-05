package book

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rakibulbanna/go-fiber-postgres/middleware"
)

func SetupRoutes(router fiber.Router, controller *Controller, authMiddleware *middleware.AuthMiddleware) {
	books := router.Group("/books")
	
	// Public routes
	books.Get("/", controller.GetBooks)
	books.Get("/:id", controller.GetBook)

	// Protected routes
	protectedBooks := router.Group("/books", authMiddleware.RequireAuth)
	protectedBooks.Post("/", controller.CreateBook)
	protectedBooks.Put("/:id", controller.UpdateBook)
	protectedBooks.Delete("/:id", controller.DeleteBook)
}

