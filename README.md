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

## Cara Menjalankan

### 1. Inisialisasi Environment

Salin `.env.example` menjadi `.env` dan sesuaikan nilainya:

```bash
cp .env.example .env
```

### 2. Jalankan Database & Migrasi (Docker)

Gunakan Docker Compose untuk menjalankan database dan membuat tabel secara otomatis:

```bash
# Jalankan database di background (jika belum)
docker-compose up -d db

# Jalankan migrasi untuk membuat tabel
docker-compose run --rm migrate
```

### 3. Opsi: Menggunakan Supabase

Jika ingin menggunakan Supabase, ganti `DATABASE_URL` di `.env` dengan URI dari Supabase, lalu jalankan migrasi:

```bash
docker run --rm -v $(pwd)/migrations:/migrations migrate/migrate -path=/migrations/ -database "URI_SUPABASE_ANDA" up
```

### 4. Cara Membuat Migrasi Baru

Jika Anda ingin menambah tabel atau mengubah skema:

```bash
make migrate-create name=nama_migrasi
```

Lalu isi file `.up.sql` yang baru muncul di folder `migrations/`.

### 5. Jalankan Aplikasi Go

```bash
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
