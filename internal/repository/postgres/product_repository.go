package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ookapratama/go-rest-api/internal/domain"
)

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p *domain.Product) error {
	query := `INSERT INTO products (name, price, description, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	err := r.db.QueryRow(ctx, query, p.Name, p.Price, p.Description, p.CreatedAt, p.UpdatedAt).Scan(&p.ID)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `SELECT id, name, price, description, created_at, updated_at FROM products WHERE id = $1`
	
	p := &domain.Product{}
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return p, nil
}

func (r *productRepository) GetAll(ctx context.Context) ([]*domain.Product, error) {
	query := `SELECT id, name, price, description, created_at, updated_at FROM products ORDER BY id DESC`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) Update(ctx context.Context, p *domain.Product) error {
	query := `UPDATE products SET name = $1, price = $2, description = $3, updated_at = $4 WHERE id = $5`
	
	p.UpdatedAt = time.Now()
	res, err := r.db.Exec(ctx, query, p.Name, p.Price, p.Description, p.UpdatedAt, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	
	res, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}
