package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ookapratama/go-rest-api/internal/domain"
	"github.com/ookapratama/go-rest-api/internal/util"
)

// MockService implements domain.ProductService
type MockService struct{}

func (m *MockService) Create(ctx context.Context, p *domain.Product) error { p.ID = 1; return nil }
func (m *MockService) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	return &domain.Product{ID: id, Name: "Mock"}, nil
}
func (m *MockService) GetAll(ctx context.Context) ([]*domain.Product, error) { return nil, nil }
func (m *MockService) Update(ctx context.Context, p *domain.Product) error { return nil }
func (m *MockService) Delete(ctx context.Context, id int64) error { return nil }
func (m *MockService) GetEnrichedAll(ctx context.Context) ([]*domain.Product, error) { return nil, nil }

func TestProductHandler_Create(t *testing.T) {
	svc := &MockService{}
	h := NewProductHandler(svc)

	p := domain.Product{
		Name:        "New Product",
		Price:       100.0,
		Description: "New Desc",
	}
	body, _ := json.Marshal(p)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.Create(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var resp util.APIResponse
	json.NewDecoder(rr.Body).Decode(&resp)

	if !resp.Success {
		t.Errorf("Expected success response, got fail")
	}
}

func TestProductHandler_Create_Invalid(t *testing.T) {
	svc := &MockService{}
	h := NewProductHandler(svc)

	// Missing required fields
	p := domain.Product{
		Name: "",
	}
	body, _ := json.Marshal(p)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.Create(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
