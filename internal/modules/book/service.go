package book

import (
	"errors"

	"github.com/rakibulbanna/go-fiber-postgres/dtos"
	"github.com/rakibulbanna/go-fiber-postgres/models"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateBook(userID uint, req *dtos.CreateBookRequest) (*dtos.BookResponse, error) {
	book := &models.Book{
		UserID:    userID,
		Author:    req.Author,
		Title:     req.Title,
		Publisher: req.Publisher,
		Year:      req.Year,
	}

	if err := s.db.Create(book).Error; err != nil {
		return nil, errors.New("failed to create book")
	}

	return &dtos.BookResponse{
		ID:        book.Id,
		UserID:    book.UserID,
		Author:    book.Author,
		Title:     book.Title,
		Publisher: book.Publisher,
		Year:      book.Year,
	}, nil
}

func (s *Service) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	if err := s.db.Joins("User").Find(&books).Error; err != nil {
		return nil, errors.New("failed to fetch books")
	}
	return books, nil
}

func (s *Service) GetBookByID(id uint) (*models.Book, error) {
	var book models.Book
	if err := s.db.Joins("User").First(&book, id).Error; err != nil {
		return nil, errors.New("book not found")
	}
	return &book, nil
}

func (s *Service) UpdateBook(id uint, userID uint, req *dtos.UpdateBookRequest) (*dtos.BookResponse, error) {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		return nil, errors.New("book not found")
	}

	// Check if user owns the book
	if book.UserID != userID {
		return nil, errors.New("unauthorized: you can only update your own books")
	}

	if req.Author != "" {
		book.Author = req.Author
	}
	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Publisher != "" {
		book.Publisher = req.Publisher
	}
	if req.Year != 0 {
		book.Year = req.Year
	}

	if err := s.db.Save(&book).Error; err != nil {
		return nil, errors.New("failed to update book")
	}

	return &dtos.BookResponse{
		ID:        book.Id,
		UserID:    book.UserID,
		Author:    book.Author,
		Title:     book.Title,
		Publisher: book.Publisher,
		Year:      book.Year,
	}, nil
}

func (s *Service) DeleteBook(id uint, userID uint) error {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		return errors.New("book not found")
	}

	// Check if user owns the book
	if book.UserID != userID {
		return errors.New("unauthorized: you can only delete your own books")
	}

	if err := s.db.Delete(&book).Error; err != nil {
		return errors.New("failed to delete book")
	}
	return nil
}

