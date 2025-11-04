package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rakibulbanna/go-fiber-postgres/models"
	"github.com/rakibulbanna/go-fiber-postgres/storage"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {
	var books []models.Book
	err := r.DB.Find(&books).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": books})
}
func (r *Repository) GetBook(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
		return nil
	}
	fmt.Println("id____: ", id)
	bookModel := &models.Book{}
	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get book"})
		return err
	}
	c.Status(fiber.StatusOK).JSON(fiber.Map{"data": bookModel})
	return nil
}
func (r *Repository) UpdateBook(c *fiber.Ctx) error {
	var book models.Book
	err := r.DB.Where("id = ?", c.Params("id")).First(&book).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": book})
}
func (r *Repository) DeleteBook(c *fiber.Ctx) error {
	bookModel := models.Book{}
	id := c.Params("id")
	if id == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID is required"})
		return nil
	}
	err := r.DB.Delete(&bookModel, id)
	if err.Error != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete book"})
		return nil
	}
	c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Book deleted successfully"})
	return nil
}
func (r *Repository) CreateBook(c *fiber.Ctx) error {
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	err := r.DB.Create(&book).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Book created successfully"})
}
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Get("/get_books", r.GetBooks)
	api.Get("/get_book/:id", r.GetBook)
	api.Put("/update_book/:id", r.UpdateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
}
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	// Run database migrations
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Error migrating database")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()

	r.SetupRoutes(app)

	app.Listen(":8080")
}
