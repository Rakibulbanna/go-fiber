package repositories

import (
	"github.com/rakibulbanna/go-fiber-postgres/models"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) FindAll() ([]models.Book, error) {
	var books []models.Book
	err := r.db.Preload("User").Find(&books).Error
	return books, err
}

// FindAllWithUsers performs a JOIN query to fetch books with user data in a single query
func (r *BookRepository) FindAllWithUsers() ([]models.Book, error) {
	var books []models.Book
	err := r.db.Joins("User").Find(&books).Error
	return books, err
}

func (r *BookRepository) FindByID(id uint) (*models.Book, error) {
	var book models.Book
	err := r.db.Joins("User").First(&book, id).Error
	return &book, err
}

func (r *BookRepository) Update(book *models.Book) error {
	return r.db.Save(book).Error
}

func (r *BookRepository) Delete(id uint) error {
	return r.db.Delete(&models.Book{}, id).Error
}
