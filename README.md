# Go REST API CRUD Template

Projek ini adalah contoh referensi pembuatan REST API menggunakan Go dengan PostgreSQL, menerapkan Service/Repository pattern, validasi, penanganan error, logging, concurrency, dan unit testing.

## Fitur

- **Clean Architecture**: Pemisahan tanggung jawab antara Domain, Handler, Service, dan Repository.
- **PostgreSQL**: Menggunakan `pgxpool` untuk koneksi database yang efisien.
- **Validasi Data**: Menggunakan `go-playground/validator`.
- **Concurrency**: Contoh eksekusi konkuren pada endpoint `GET /products` (metode `GetEnrichedAll`).
- **Logging**: Implementasi `slog` (JSON output) untuk mencatat error dan aktivitas aplikasi.
- **Graceful Shutdown**: Menangani signal interrupt untuk menutup koneksi database dan server dengan rapi.
- **Unit Testing**: Pengujian untuk unit logic (service) dan endpoint (handler).

## Struktur Folder

```text
.
├── cmd/server/          # Entry point aplikasi (main.go)
├── internal/
│   ├── api/             # HTTP Handler, Router, Middleware
│   ├── domain/          # Model Data & Interfaces (Kontrak)
│   ├── repository/      # Implementasi Database (Postgres)
│   ├── service/         # Business Logic & Concurrency
│   └── util/            # Helper (Response JSON, Validator)
├── migrations/          # File SQL untuk skema database
└── README.md
```

## Persyaratan

- Go 1.21+
- PostgreSQL
- Docker (Opsional, untuk menjalankan DB dengan mudan)

## Cara Menjalankan

### 1. Persiapkan Database

Jika menggunakan Docker:

```bash
docker run --name postgres-db -e POSTGRES_PASSWORD=password -e POSTGRES_DB=product_db -p 5432:5432 -d postgres
```

Lalu jalankan query di `migrations/001_create_products_table.up.sql` untuk membuat tabel.

### 2. Konfigurasi Environment

Salin `.env.example` menjadi `.env` dan sesuaikan nilainya:

```bash
cp .env.example .env
```

Isi `.env`:

```env
PORT=8080
DATABASE_URL=postgres://postgres:password@localhost:5432/product_db?sslmode=disable
```

### 3. Jalankan Aplikasi

```bash
go mod tidy
go run cmd/server/main.go
```

## API Endpoints

| Method | Endpoint                | Deskripsi                                              |
| ------ | ----------------------- | ------------------------------------------------------ |
| POST   | `/api/v1/products`      | Menambah produk baru                                   |
| GET    | `/api/v1/products`      | Mengambil semua produk (dengan concurrency enrichment) |
| GET    | `/api/v1/products/{id}` | Mengambil detail produk                                |
| PUT    | `/api/v1/products/{id}` | Memperbarui data produk                                |
| DELETE | `/api/v1/products/{id}` | Menghapus produk                                       |

### Contoh Request (Create Product)

```json
{
  "name": "Macbook Pro M3",
  "price": 2500.0,
  "description": "Powerful laptop for developers"
}
```

## Menjalankan Unit Test

```bash
go test ./internal/service/...
go test ./internal/api/handler/...
# Atau semua test
go test ./...
```
