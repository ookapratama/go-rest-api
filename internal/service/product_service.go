package service

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/ookapratama/go-rest-api/internal/domain"
)

type productService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) domain.ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, p *domain.Product) error {
	return s.repo.Create(ctx, p)
}

func (s *productService) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) GetAll(ctx context.Context) ([]*domain.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *productService) Update(ctx context.Context, p *domain.Product) error {
	return s.repo.Update(ctx, p)
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// GetEnrichedAll demonstrates CONCURRENCY
// It simulates fetching extra data for each product from an external source in parallel.
func (s *productService) GetEnrichedAll(ctx context.Context) ([]*domain.Product, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	// We use a channel or a mutex-protected slice to handle concurrent updates.
	// Since we are updating the existing slice elements, and each goroutine 
	// updates a unique index, we don't strictly need a mutex, but it's good practice
	// if we were appending.
	
	for _, p := range products {
		wg.Add(1)
		go func(prod *domain.Product) {
			defer wg.Done()
			
			// Simulate a slow external call (e.g., getting stock stats or something)
			// In a real app, this could be a request to another microservice.
			enrichProduct(prod)
			
			slog.Info("Product enriched concurrently", "id", prod.ID)
		}(p)
	}

	wg.Wait()
	return products, nil
}

func enrichProduct(p *domain.Product) {
	// Simulate processing time
	time.Sleep(10 * time.Millisecond)
	p.Description = "[Enriched] " + p.Description
}
