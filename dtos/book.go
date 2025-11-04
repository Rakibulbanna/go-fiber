package dtos

type CreateBookRequest struct {
	Author    string `json:"author" validate:"required"`
	Title     string `json:"title" validate:"required"`
	Publisher string `json:"publisher" validate:"required"`
	Year      int    `json:"year" validate:"required,min=1000,max=9999"`
}

type UpdateBookRequest struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
	Year      int    `json:"year" validate:"omitempty,min=1000,max=9999"`
}

type BookResponse struct {
	ID        uint          `json:"id"`
	UserID    uint          `json:"user_id"`
	Author    string        `json:"author"`
	Title     string        `json:"title"`
	Publisher string        `json:"publisher"`
	Year      int           `json:"year"`
	User      *UserResponse `json:"user"`
}
