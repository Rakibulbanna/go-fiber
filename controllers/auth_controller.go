package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rakibulbanna/go-fiber-postgres/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) SignUp(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

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

	signUpReq := &services.SignUpRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	response, err := c.authService.SignUp(signUpReq)
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

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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

	loginReq := &services.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	response, err := c.authService.Login(loginReq)
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
