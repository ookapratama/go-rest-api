package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ookapratama/go-rest-api/latihan_mandiri/internal/domain"
)

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, p *domain.Product) error {
	// Query SQL dengan RETURNING agar kita dapat ID yang baru digenerate oleh DB
	query := `INSERT INTO products (name, price, description, created_at, updated_at) 
			  VALUES ($1, $2, $3, NOW(), NOW()) 
			  RETURNING id, created_at, updated_at`

	// Scan hasil returning ke dalam struct
	err := r.db.QueryRow(ctx, query, p.Name, p.Price, p.Description).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("error creating product: %v", err)
	}
	return nil
}

func (r *productRepository) GetAll(ctx context.Context) ([]*domain.Product, error) {
	query := `SELECT id, name, price, description, created_at, updated_at FROM products`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// TODO: Anda bisa mencoba melengkapi GetByID, Update, dan Delete untuk latihan mandiri
func (r *productRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) { return nil, nil }
func (r *productRepository) Update(ctx context.Context, p *domain.Product) error     { return nil }
func (r *productRepository) Delete(ctx context.Context, id int64) error             { return nil }
