package dto

type BookData struct {
	ID          string `json:"id"`
	Isbn        string `json:"isbn"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type BookStockData struct {
	Code   string `json:"code"`
	Status string `json:"status"`
}

type BookShowData struct {
	BookData
	Stocks []BookStockData `json:"stocks"`
}

type CreateBookRequest struct {
	Isbn        string `json:"isbn"  validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateBookRequest struct {
	ID          string `json:"-"`
	Isbn        string `json:"isbn"  validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}
