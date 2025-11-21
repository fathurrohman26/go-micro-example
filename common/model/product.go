package model

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type GetProductRequest struct {
	ID string `json:"id"`
}

type GetProductResponse struct {
	Product Product `json:"product"`
}
