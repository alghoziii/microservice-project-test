package controllers

import (
	"order-service/domain/dto"
	"order-service/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService *services.OrderService
}

func NewOrderController(orderService *services.OrderService) *OrderController {
	return &OrderController{orderService}
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req dto.OrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Data tidak valid"})
		return
	}

	order, err := c.orderService.CreateOrder(req)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"data": order})
}

func (c *OrderController) GetOrderByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	order, err := c.orderService.GetOrderByID(uint(id))

	if err != nil {
		ctx.JSON(404, gin.H{"error": "Order tidak ditemukan"})
		return
	}

	ctx.JSON(200, gin.H{"data": order})
}

func (c *OrderController) GetOrdersByProductID(ctx *gin.Context) {
	productID, _ := strconv.Atoi(ctx.Param("productId"))
	orders, err := c.orderService.GetOrdersByProductID(productID)

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Gagal mengambil data"})
		return
	}

	ctx.JSON(200, gin.H{"data": orders})
}
