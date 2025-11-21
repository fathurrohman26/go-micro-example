package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"elproxy.cloud/elproxy-server/common/model"
	"go-micro.dev/v5"
	"go-micro.dev/v5/client"
)

type OrdersService struct {
	mu      sync.Mutex
	counter int
	orders  map[string]*model.Order
	client  client.Client
}

func (s *OrdersService) Create(ctx context.Context, req *model.CreateOrderRequest, rsp *model.CreateOrderResponse) error {
	userReq := s.client.NewRequest("users", "UsersService.Get", &model.GetUserRequest{ID: req.UserID})
	userRsp := &model.GetUserResponse{}
	if err := s.client.Call(ctx, userReq, userRsp); err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	prodReq := s.client.NewRequest("products", "ProductsService.Get", &model.GetProductRequest{ID: req.ProductID})
	prodRsp := &model.GetProductResponse{}
	if err := s.client.Call(ctx, prodReq, prodRsp); err != nil {
		return fmt.Errorf("failed to get product: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.counter++
	orderID := fmt.Sprintf("order_%d", s.counter)
	order := &model.Order{
		ID:        orderID,
		UserID:    req.UserID,
		ProductID: req.ProductID,
		Amount:    req.Amount,
		CreatedAt: time.Now(),
	}
	s.orders[orderID] = order
	rsp.Order = order
	return nil
}

func main() {
	svc := micro.NewService(
		micro.Name("orders"),
		micro.Version("1.0.0"),
	)
	svc.Init()

	orderService := newOrdersService(svc.Client())
	svc.Handle(orderService)

	if err := svc.Run(); err != nil {
		fmt.Println("Error running service:", err)
	}
}

func newOrdersService(c client.Client) *OrdersService {
	return &OrdersService{
		orders:  make(map[string]*model.Order),
		counter: 0,
		client:  c,
	}
}
