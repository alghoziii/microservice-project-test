package models

import "time"

type Order struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	ProductID  int       `gorm:"not null" json:"productId"`
	TotalPrice float64   `gorm:"type:decimal(10,2)" json:"totalPrice"`
	Status     string    `gorm:"type:varchar(50); default:'pending'" json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (Order) TableName() string {
	return "orders"
}
