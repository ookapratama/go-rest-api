package domain

import (
	"context"
	"time"
)

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Price       float64   `json:"price" validate:"required,gt=0"`
	Description string    `json:"description" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRepository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id int64) (*Product, error)
	GetAll(ctx context.Context) ([]*Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id int64) error
}

type ProductService interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id int64) (*Product, error)
	GetAll(ctx context.Context) ([]*Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id int64) error
	// GetEnrichedAll demonstrates concurrency
	GetEnrichedAll(ctx context.Context) ([]*Product, error)
}
