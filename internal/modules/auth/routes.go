package auth

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(router fiber.Router, controller *Controller) {
	auth := router.Group("/auth")
	
	auth.Post("/signup", controller.SignUp)
	auth.Post("/login", controller.Login)
	auth.Post("/signin", controller.Login) // Alias for login
}

