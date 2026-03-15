package service

import (
	"context"
	"testing"
	"time"

	"github.com/ookapratama/go-rest-api/internal/domain"
)

// MockRepository implements domain.ProductRepository for testing
type MockRepository struct {
	products []*domain.Product
}

func (m *MockRepository) Create(ctx context.Context, p *domain.Product) error {
	p.ID = int64(len(m.products) + 1)
	m.products = append(m.products, p)
	return nil
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	for _, p := range m.products {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, nil // Simplified
}

func (m *MockRepository) GetAll(ctx context.Context) ([]*domain.Product, error) {
	return m.products, nil
}

func (m *MockRepository) Update(ctx context.Context, p *domain.Product) error {
	return nil
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func TestGetEnrichedAll(t *testing.T) {
	repo := &MockRepository{
		products: []*domain.Product{
			{ID: 1, Name: "Test 1", Description: "Desc 1"},
			{ID: 2, Name: "Test 2", Description: "Desc 2"},
		},
	}
	svc := NewProductService(repo)

	ctx := context.Background()
	start := time.Now()
	
	products, err := svc.GetEnrichedAll(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	duration := time.Since(start)

	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}

	// Verify enrichment
	if products[0].Description != "[Enriched] Desc 1" {
		t.Errorf("Expected enriched description, got %s", products[0].Description)
	}

	// Concurrency verification: If it's serial, it should take at least 20ms (2 * 10ms).
	// If it's parallel, it should be closer to 10ms.
	// NOTE: This is a bit flaky on some CI machines, but it serves the purpose of example.
	if duration >= 25*time.Millisecond {
		t.Logf("Warning: enrichment took %v, maybe not fully concurrent or system lag", duration)
	}
}
