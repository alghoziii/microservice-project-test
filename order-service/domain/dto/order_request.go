package dto

type OrderRequest struct {
	ProductID int `json:"productId" binding:"required"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}
