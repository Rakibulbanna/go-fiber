package services

import (
	"errors"

	"github.com/rakibulbanna/go-fiber-postgres/dtos"
	"github.com/rakibulbanna/go-fiber-postgres/models"
	"github.com/rakibulbanna/go-fiber-postgres/repositories"
)

type CreateBookRequest struct {
	Author    string
	Title     string
	Publisher string
	Year      int
}

type UpdateBookRequest struct {
	Author    string
	Title     string
	Publisher string
	Year      int
}

type BookService struct {
	bookRepo *repositories.BookRepository
}

func NewBookService(bookRepo *repositories.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func (s *BookService) CreateBook(userID uint, req *CreateBookRequest) (*dtos.BookResponse, error) {
	book := &models.Book{
		UserID:    userID,
		Author:    req.Author,
		Title:     req.Title,
		Publisher: req.Publisher,
		Year:      req.Year,
	}

	if err := s.bookRepo.Create(book); err != nil {
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

func (s *BookService) GetAllBooks() ([]dtos.BookResponse, error) {
	books, err := s.bookRepo.FindAll()
	if err != nil {
		return nil, errors.New("failed to fetch books")
	}

	var response []dtos.BookResponse
	for _, book := range books {
		response = append(response, dtos.BookResponse{
			ID:        book.Id,
			UserID:    book.UserID,
			Author:    book.Author,
			Title:     book.Title,
			Publisher: book.Publisher,
			Year:      book.Year,
		})
	}

	return response, nil
}

func (s *BookService) GetBookByID(id uint) (*dtos.BookResponse, error) {
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("book not found")
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

func (s *BookService) UpdateBook(id uint, userID uint, req *UpdateBookRequest) (*dtos.BookResponse, error) {
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
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

	if err := s.bookRepo.Update(book); err != nil {
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

func (s *BookService) DeleteBook(id uint, userID uint) error {
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return errors.New("book not found")
	}

	// Check if user owns the book
	if book.UserID != userID {
		return errors.New("unauthorized: you can only delete your own books")
	}

	if err := s.bookRepo.Delete(id); err != nil {
		return errors.New("failed to delete book")
	}
	return nil
}
