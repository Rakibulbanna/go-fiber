package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rakibulbanna/go-fiber-postgres/dtos"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) SignUp(ctx *fiber.Ctx) error {
	var req dtos.SignUpRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email, password, and name are required",
		})
	}

	if len(req.Password) < 6 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters",
		})
	}

	response, err := c.service.SignUp(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"data":    response,
	})
}

func (c *Controller) Login(ctx *fiber.Ctx) error {
	var req dtos.LoginRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	response, err := c.service.Login(&req)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"data":    response,
	})
}

