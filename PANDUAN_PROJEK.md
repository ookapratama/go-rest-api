# Panduan Step-by-Step Membuat REST API dengan Go (Standard Boilerplate)

Panduan ini berisi alur kerja (_workflow_) yang bisa Anda ikuti setiap kali memulai proyek REST API baru menggunakan arsitektur Service-Repository.

---

## 1. Tahap Inisialisasi Proyek

Langkah awal untuk menyiapkan struktur folder dan dependensi.

1.  **Init Module:**
    ```bash
    go mod init github.com/username/project-name
    ```
2.  **Buat Struktur Folder Dasar:**
    ```bash
    mkdir -p cmd/server internal/api/handler internal/domain internal/repository/postgres internal/service internal/util migrations
    ```
3.  **Install Library Utama:**
    - Router: `go get github.com/go-chi/chi/v5`
    - DB Driver: `go get github.com/jackc/pgx/v5`
    - Validation: `go get github.com/go-playground/validator/v10`
    - Env Loader: `go get github.com/joho/godotenv`

---

## 2. Tahap Definisi Domain (The Blueprint)

Mulailah dari **Domain**, karena ini adalah jantung dari aplikasi Anda.

1.  **Buat Model Struct:** Definisikan data yang akan disimpan (misal: `Product`, `User`). Tambahkan tag `json` dan `validate`.
2.  **Definisikan Interface:** Buat kontrak untuk Repository dan Service.
    - _Kenapa interface?_ Supaya antar layer tidak saling ketergantungan secara langsung (Decoupling) dan mudah di-test.

---

## 3. Tahap Implementasi Repository (Layer Database)

Fokus pada urusan _Query_ ke database.

1.  **Buat SQL Migration:** Tulis schema tabel di folder `migrations/`.
2.  **Implementasikan Interface Repository:** Buat struct yang menampung koneksi DB (`*pgxpool.Pool`) dan tulis fungsi-fungsi CRUD-nya.
3.  **Error Handling:** Pastikan setiap query menangani error dan mengembalikan pesan yang bermakna.

---

## 4. Tahap Implementasi Service (Layer Business Logic)

Fokus pada _Aturan Bisnis_ dan _Logic_.

1.  **Implementasikan Interface Service:** Masukkan Repository ke dalam struct Service.
2.  **Tambahkan Logic Tambahan:** Misalnya perhitungan harga, pengecekan stok, atau **Concurrency** (eksekusi paralel menggunakan Goroutines).
3.  **Unit Test (Logic):** Buat file `_test.go` untuk mengetes logika di layer ini menggunakan _Mock Repository_.

---

## 5. Tahap API Handler & Router (Layer Transport)

Fokus pada komunikasi HTTP.

1.  **Buat Utility Response:** (Opsional) Buat helper untuk standarisasi format JSON success/error.
2.  **Buat Handler:** Struct yang memanggil Service. Tugasnya:
    - Decode JSON request body.
    - Validasi input.
    - Panggil Service.
    - Balas dengan JSON response.
3.  **Setup Router:** Hubungkan URL path dengan fungsi Handler yang sesuai. Gunakan middleware standard (Logging, Recoverer).

---

## 6. Tahap Entry Point & Wiring (Main.go)

Menyatukan semua potongan puzzle.

1.  **Log Setup:** Inisialisasi logger (format JSON direkomendasikan).
2.  **Connect DB:** Baca dari `.env`, buka pool koneksi ke Postgres.
3.  **Dependency Injection:**
    - Init Repo -> Masukkan ke Service.
    - Init Service -> Masukkan ke Handler.
    - Init Handler -> Masukkan ke Router.
4.  **Graceful Shutdown:** Tambahkan mekanisme untuk menangani signal `SIGINT/SIGTERM` agar server mati dengan rapi (tutup DB pool, dll).

---

## 7. Tahap Testing & Dokumentasi

1.  **Unit Test Handler:** Tes endpoint tanpa butuh database asli (pakai Mock Service).
2.  **README:** Tulis cara menjalankan aplikasi dan daftar endpoint yang tersedia.
3.  **Dockerize (Optional):** Buat `Dockerfile` agar aplikasi mudah di-deploy di mana saja.

---

## 8. Database Migrations (Membuat Tabel)

Setelah Docker Postgres berjalan, Anda perlu membuat tabel. Kita menggunakan tool `golang-migrate`.

### A. Jika Menggunakan Local Docker (Rekomendasi)

Kita sudah mensetup container `migrate` di `docker-compose.yml`. Jalankan perintah ini di terminal:

```bash
docker-compose run --rm migrate
```

**Tujuan Perintah:**

- `docker-compose run`: Menjalankan service tertentu secara mandiri.
- `--rm`: Menghapus container setelah selesai (supaya tidak menumpuk sampah container).
- `migrate`: Nama service di `docker-compose.yml` yang berisi konfigurasi migrasi.
- **Hasilnya:** Semua file `.sql` di folder `/migrations` akan dieksekusi ke database Postgres Anda untuk membuat tabel.

### B. Opsi: Menggunakan Supabase PostgreSQL

Jika Anda lebih suka menggunakan database cloud dari **Supabase**, ikuti langkah ini:

1.  **Dapatkan Connection String:**
    - Buka dashboard Supabase.
    - Settings -> Database.
    - Cari **Connection String (URI)**. Pilih mode `Transaction` atau `Session`.
2.  **Edit File `.env`:**
    Ganti `DATABASE_URL` dengan string dari Supabase.
    - Contoh: `DATABASE_URL=postgres://postgres:[PASSWORD]@db.xxxx.supabase.co:5432/postgres`
3.  **Jalankan Migrasi ke Supabase:**
    Gunakan Docker untuk mendorong tabel ke Supabase (tidak perlu install tool tambahan di laptop):
    ```bash
    docker run --rm -v $(pwd)/migrations:/migrations migrate/migrate -path=/migrations/ -database "ISI_DENGAN_URL_SUPABASE_ANDA" up
    ```

### C. Cara Membuat File Migrasi Baru

Setiap kali Anda ingin menambah tabel baru atau mengubah struktur tabel, Anda harus membuat file migrasi baru.

**Gunakan perintah ini (via Makefile):**

```bash
make migrate-create name=keterangan_migrasi
```

_(Ganti `keterangan_migrasi` sesuai kebutuhan, misal: `create_users_table` atau `add_column_to_products`)_

**Penjelasan:**

- Perintah ini akan menghasilkan 2 file di folder `/migrations`:
  1.  `XXX_keterangan.up.sql`: Isi dengan kode SQL untuk **menambah/mengubah** (misal: `CREATE TABLE`).
  2.  `XXX_keterangan.down.sql`: Isi dengan kode SQL untuk **membatalkan** (misal: `DROP TABLE`).
- **Nomor Urut:** Tool akan otomatis memberi nomor urut (seperti `001`, `002`) agar urutan eksekusi tetap benar.

---

## Ringkasan Perintah Penting (Cheat Sheet)

| Perintah                          | Tujuan                                            |
| :-------------------------------- | :------------------------------------------------ |
| `make migrate-create name=xxx`    | **Membuat** file migrasi baru (.sql).             |
| `docker-compose up -d db`         | Menjalankan database Postgres di background.      |
| `docker-compose run --rm migrate` | **Eksekusi** migrasi ke database (Membuat tabel). |
| `go run cmd/server/main.go`       | Menjalankan aplikasi Go Anda.                     |
| `docker-compose logs -f db`       | Melihat log error database jika koneksi gagal.    |

---

## Tips Belajar:

- **Jangan Hafalkan Code:** Hafalkan **Alur Masuk Data-nya**: `Request -> Router -> Handler -> Service -> Repository -> Database`.
- **Gunakan Interface:** Jika Anda ingin mengganti database dari Postgres ke MongoDB, Anda hanya perlu membuat implementasi Repository baru tanpa menyentuh layer Service/Handler.
- **Log Everything Error:** Selalu log error di layer Handler atau Main agar gampang di-debug.
