package dto

type CreateBookRequest struct {
	Title  string `json:"title" binding:"required,min=3,max=255"`
	Author string `json:"author" binding:"required,min=3"`
	ISBN   string `json:"isbn"`
	Year   int    `json:"year" binding:"required,min=1900,max=2026"`
	Stock  int    `json:"stock" binding:"required,min=0"`
}

type UpdateBookRequest struct {
	Title     string `json:"title" binding:"omitempty,min=3,max=255"`
	Author    string `json:"author" binding:"omitempty,min=3"`
	ISBN      string `json:"isbn"`
	Year      int    `json:"year" binding:"omitempty,min=1900,max=2026"`
	Stock     int    `json:"stock" binding:"omitempty,min=0"`
	Available int    `json:"available" binding:"omitempty,min=0"`
}

type BookResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	ISBN      string `json:"isbn"`
	Year      int    `json:"year"`
	Stock     int    `json:"stock"`
	Available int    `json:"available"`
}
