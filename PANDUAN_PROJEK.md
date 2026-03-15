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

## Tips Belajar:

- **Jangan Hafalkan Code:** Hafalkan **Alur Masuk Data-nya**: `Request -> Router -> Handler -> Service -> Repository -> Database`.
- **Gunakan Interface:** Jika Anda ingin mengganti database dari Postgres ke MongoDB, Anda hanya perlu membuat implementasi Repository baru tanpa menyentuh layer Service/Handler.
- **Log Everything Error:** Selalu log error di layer Handler atau Main agar gampang di-debug.
