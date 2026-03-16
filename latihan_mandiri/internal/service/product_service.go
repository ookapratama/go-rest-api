package service

import (
	"context"
	"sync"

	"github.com/ookapratama/go-rest-api/latihan_mandiri/internal/domain"
)

type productService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) domain.ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, p *domain.Product) error {
	// Di sini bisa ditambahkan logika misal: validasi stok, cek nama duplikat, dll
	return s.repo.Create(ctx, p)
}

func (s *productService) GetAll(ctx context.Context) ([]*domain.Product, error) {
	return s.repo.GetAll(ctx)
}

// GetEnrichedAll: DEMO CONCURRENCY (Sering ditanya saat interview)
func (s *productService) GetEnrichedAll(ctx context.Context) ([]*domain.Product, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup // Alat untuk menunggu semua goroutine selesai

	for _, p := range products {
		wg.Add(1) // Tambah antrian
		go func(prod *domain.Product) {
			defer wg.Done() // Kurangi antrian jika selesai
			
			// Simulasi proses tambahan (Misal: ambil rating dari API lain)
			prod.Description = "[PROMOSI] " + prod.Description
		}(p)
	}

	wg.Wait() // Tunggu sampai antrian (WaitGroup) jadi 0
	return products, nil
}

// Tambahan agar sesuai interface
func (s *productService) GetByID(ctx context.Context, id int64) (*domain.Product, error) { return nil, nil }
func (s *productService) Update(ctx context.Context, p *domain.Product) error     { return nil }
func (s *productService) Delete(ctx context.Context, id int64) error             { return nil }
