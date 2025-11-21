package main

import (
	"context"
	"fmt"
	"sync"

	"elproxy.cloud/elproxy-server/common/model"
	"go-micro.dev/v5"
)

type ProductsService struct {
	mu       sync.Mutex
	products map[string]*model.Product
}

func (s *ProductsService) Get(ctx context.Context, req *model.GetProductRequest, rsp *model.GetProductResponse) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[req.ID]
	if !exists {
		return fmt.Errorf("product with ID %s not found", req.ID)
	}

	*rsp = model.GetProductResponse{Product: *product}
	return nil
}

func main() {
	svc := micro.NewService(
		micro.Name("products"),
		micro.Version("1.0.0"),
	)
	svc.Init()

	productsService := newProductsService()
	svc.Handle(productsService)

	if err := svc.Run(); err != nil {
		fmt.Println("Error running service:", err)
	}
}

func newProductsService() *ProductsService {
	products := []*model.Product{
		{ID: "101", Name: "Laptop", Price: 999.99},
		{ID: "102", Name: "Smartphone", Price: 499.99},
	}

	productMap := make(map[string]*model.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	return &ProductsService{
		products: productMap,
	}
}
