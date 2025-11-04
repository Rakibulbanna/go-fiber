package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rakibulbanna/go-fiber-postgres/utils"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) RequireAuth(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is required",
		})
	}

	// Extract token from "Bearer <token>"
	token := ""
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	// Validate token
	claims, err := utils.ValidateToken(token, m.jwtSecret)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Store user info in context
	ctx.Locals("userID", claims.UserID)
	ctx.Locals("userEmail", claims.Email)

	return ctx.Next()
}
