package dto

type CreateBookRequest struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	ISBN   string `json:"isbn"`
	Stock  int `json:"stock"`
}

type UpdateBookRequest struct{
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
	Stock  int `json:"stock"`
}

type BookResponse struct {
	ID  uint `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
	Stock  int `json:"stock"`
}
