package services

import (
	"encoding/json"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"order-service/domain/dto"
	"order-service/domain/models"
	"testing"
)

// fake repo & kafka
type fakeRepo struct{ called bool }

func (r *fakeRepo) Create(o *models.Order) error                    { r.called = true; o.ID = 1; return nil }
func (r *fakeRepo) FindByID(id uint) (*models.Order, error)         { return nil, nil }
func (r *fakeRepo) FindByProductID(pid int) ([]models.Order, error) { return nil, nil }

type fakeKafka struct {
	orderID        uint
	productID, qty int
}

func (k *fakeKafka) SendOrderCreated(orderID uint, productID int, qty int) {
	k.orderID, k.productID, k.qty = orderID, productID, qty
}

// adapter biar cocok ke tipe KafkaProducer di service
type kafkaProducerAdapter fakeKafka

func (k *kafkaProducerAdapter) SendOrderCreated(orderID uint, productID int, qty int) {
	(*fakeKafka)(k).SendOrderCreated(orderID, productID, qty)
}

func TestCreateOrder_Success(t *testing.T) {
	// mock product-service HTTP
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": 1, "name": "A", "price": 10000, "qty": 10,
		})
	}))
	defer ts.Close()

	// mock redis
	rdb, rmock := redismock.NewClientMock()
	rmock.ExpectGet("product:1").RedisNil()
	rmock.ExpectDel("orders:product:1").SetVal(1)

	repo := &fakeRepo{}
	kf := &fakeKafka{}
	svc := NewOrderService(repo, rdb, (*kafkaProducerAdapter)(kf), ts.URL, ts.Client())

	res, err := svc.CreateOrder(dto.OrderRequest{ProductID: 1, Quantity: 2})
	assert.NoError(t, err)
	assert.Equal(t, float64(20000), res.TotalPrice)
	assert.True(t, repo.called)
	assert.Equal(t, 1, kf.productID)
	assert.Equal(t, 2, kf.qty)
	assert.NoError(t, rmock.ExpectationsWereMet())
}
