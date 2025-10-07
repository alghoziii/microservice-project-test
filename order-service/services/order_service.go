package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order-service/domain/dto"
	"order-service/domain/models"
	"order-service/repositories"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Interface untuk Kafka (biar bisa di-mock)
type KafkaProducer interface {
	SendOrderCreated(orderID uint, productID int, qty int)
}

type OrderService struct {
	orderRepo      repositories.OrderRepository
	redisClient    *redis.Client
	kafkaProducer  KafkaProducer
	productBaseURL string
	httpClient     *http.Client
}

func NewOrderService(
	orderRepo repositories.OrderRepository,
	redisClient *redis.Client,
	kafkaProducer KafkaProducer,
	productBaseURL string,
	httpClient *http.Client,
) *OrderService {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &OrderService{
		orderRepo:      orderRepo,
		redisClient:    redisClient,
		kafkaProducer:  kafkaProducer,
		productBaseURL: productBaseURL,
		httpClient:     httpClient,
	}
}

func (s *OrderService) CreateOrder(req dto.OrderRequest) (*dto.OrderResponse, error) {
	// 1. Ambil data produk
	product, err := s.getProduct(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product error: %v", err)
	}

	// 2. Cek stok
	if product.Qty < req.Quantity {
		return nil, fmt.Errorf("stock kurang. tersedia: %d", product.Qty)
	}

	// 3. Simpan order
	order := &models.Order{
		ProductID:  req.ProductID,
		TotalPrice: product.Price * float64(req.Quantity),
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, err
	}

	// 4. Hapus cache & kirim event
	s.redisClient.Del(ctx, "orders:product:"+strconv.Itoa(req.ProductID))
	s.kafkaProducer.SendOrderCreated(order.ID, req.ProductID, req.Quantity)

	log.Printf("Order created: ID=%d, Product=%s", order.ID, product.Name)

	return &dto.OrderResponse{
		ID:         order.ID,
		ProductID:  order.ProductID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
	}, nil
}

func (s *OrderService) GetOrderByID(id uint) (*dto.OrderResponse, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.OrderResponse{
		ID:         order.ID,
		ProductID:  order.ProductID,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt,
	}, nil
}


// ✅ Get Orders by Product ID (with cache)
func (s *OrderService) GetOrdersByProductID(productID int) ([]dto.OrderResponse, error) {
	key := "orders:product:" + strconv.Itoa(productID)
	if cached, err := s.redisClient.Get(ctx, key).Result(); err == nil {
		var res []dto.OrderResponse
		json.Unmarshal([]byte(cached), &res)
		return res, nil
	}

	orders, err := s.orderRepo.FindByProductID(productID)
	if err != nil {
		return nil, err
	}

	var response []dto.OrderResponse
	for _, o := range orders {
		response = append(response, dto.OrderResponse{
			ID:         o.ID,
			ProductID:  o.ProductID,
			TotalPrice: o.TotalPrice,
			Status:     o.Status,
			CreatedAt:  o.CreatedAt,
		})
	}

	if data, err := json.Marshal(response); err == nil {
		s.redisClient.Set(ctx, key, string(data), 5*time.Minute)
	}

	return response, nil
}

// ✅ Helper untuk ambil data produk (via HTTP atau cache)
func (s *OrderService) getProduct(productID int) (*struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}, error) {
	key := "product:" + strconv.Itoa(productID)
	if cached, err := s.redisClient.Get(ctx, key).Result(); err == nil {
		var p struct {
			ID    uint    `json:"id"`
			Name  string  `json:"name"`
			Price float64 `json:"price"`
			Qty   int     `json:"qty"`
		}
		json.Unmarshal([]byte(cached), &p)
		return &p, nil
	}

	resp, err := s.httpClient.Get(fmt.Sprintf("%s/products/%d", s.productBaseURL, productID))
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("product tidak ditemukan")
	}
	defer resp.Body.Close()

	var p struct {
		ID    uint    `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
		Qty   int     `json:"qty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}

	if data, err := json.Marshal(p); err == nil {
		s.redisClient.Set(ctx, key, string(data), 5*time.Minute)
	}

	return &p, nil
}
