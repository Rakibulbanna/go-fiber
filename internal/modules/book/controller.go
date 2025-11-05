package book

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rakibulbanna/go-fiber-postgres/dtos"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) CreateBook(ctx *fiber.Ctx) error {
	var req dtos.CreateBookRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Title == "" || req.Publisher == "" || req.Author == "" || req.Year == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title, author, publisher, and year are required",
		})
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := ctx.Locals("userID").(uint)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	book, err := c.service.CreateBook(userID, &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Book created successfully",
		"data":    book,
	})
}

func (c *Controller) GetBooks(ctx *fiber.Ctx) error {
	books, err := c.service.GetAllBooks()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": books,
	})
}

func (c *Controller) GetBook(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	if idParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	book, err := c.service.GetBookByID(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": book,
	})
}

func (c *Controller) UpdateBook(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	if idParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	var req dtos.UpdateBookRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := ctx.Locals("userID").(uint)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	book, err := c.service.UpdateBook(uint(id), userID, &req)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book updated successfully",
		"data":    book,
	})
}

func (c *Controller) DeleteBook(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	if idParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := ctx.Locals("userID").(uint)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	if err := c.service.DeleteBook(uint(id), userID); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Book deleted successfully",
	})
}

