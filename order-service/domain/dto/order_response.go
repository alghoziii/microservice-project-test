package dto

import "time"

type OrderResponse struct {
	ID         uint      `json:"id"`
	ProductID  int       `json:"productId"`
	TotalPrice float64   `json:"totalPrice"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}
