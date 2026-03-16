# Persiapan Interview: Go REST API

Dokumen ini berisi poin-poin kunci yang sering ditanyakan saat interview teknis untuk posisi Go Backend Developer berdasarkan proyek ini.

## 1. Arsitektur Proyek

**T: Kenapa menggunakan struktur folder seperti ini?**

- **J:** Kita menggunakan arsitektur berlapis (Clean Architecture/Hexagonal-ish). Tujuannya adalah **Separation of Concerns**. Layer Handler tidak perlu tahu cara kerja database, dan layer Repository tidak perlu tahu cara kerja HTTP.

**T: Apa fungsi folder `internal`?**

- **J:** Folder `internal` di Go diproteksi oleh compiler. Kode di dalamnya hanya bisa diimport oleh modul yang berada di root yang sama. Ini memastikan library luar tidak bisa mengakses logika internal aplikasi kita secara sembarangan.

## 2. Dependency Injection (DI)

**T: Bagaimana Anda melakukan Dependency Injection di proyek ini?**

- **J:** Kita melakukannya secara manual di `main.go`. Repository diinject ke Service, dan Service diinject ke Handler melalui konstruktor (`NewService`, `NewHandler`).
- **Manfaat:** Memudahkan **Unit Testing** (kita bisa memasukkan "Mock" object) dan membuat kode lebih modular.

## 3. Interface

**T: Mengapa Repository dan Service menggunakan Interface?**

- **J:** Agar antar layer tidak tergantung pada implementasi konkret (_Decoupling_). Jika suatu saat kita ingin ganti Database dari Postgres ke MongoDB, kita cukup buat implementasi Repository baru tanpa harus mengubah kode di layer Service.

## 4. Concurrency (Goroutines)

**T: Di bagian mana Anda menggunakan Concurrency?**

- **J:** Di `product_service.go`, pada fungsi `GetEnrichedAll`. Kita menggunakan **Goroutines** untuk mengambil data tambahan secara paralel untuk setiap produk.
- **T: Bagaimana Anda memastikan semua Goroutine selesai?**
- **J:** Menggunakan `sync.WaitGroup`. Kita `Add(1)` sebelum memulai goroutine dan memanggil `Done()` saat selesai. Lalu `Wait()` di akhir untuk menunggu semuanya rampung.

## 5. Error Handling & Logging

**T: Bagaimana strategi Error Handling Anda?**

- **J:** Kita menangkap error sedini mungkin. Di layer Repository, error database dibungkus. Di layer Handler, kita menentukan HTTP Status Code yang sesuai (400, 404, 500) berdasarkan jenis error yang diterima dari Service.
- **Logging:** Kita menggunakan `slog` (Structured Logging) agar log mudah dibaca oleh sistem monitoring (seperti ELK atau Datadog) dalam format JSON.

## 6. Database & Migration

**T: Kenapa menggunakan Migration Tool daripada membuat tabel manual?**

- **J:** Agar perubahan skema database tercatat dalam history (Version Control). Anggota tim lain bisa menyamakan struktur database mereka hanya dengan menjalankan perintah migrate.

---

## Tip Interview:

Jika ditanya _"Apa yang bisa ditingkatkan dari proyek ini?"_, Anda bisa menjawab:

1. Menambahkan **Middleware** untuk Auth (JWT).
2. Menambahkan **Redis** untuk Caching.
3. Menambahkan **Swagger/OpenAPI** untuk dokumentasi API.
4. Menambahkan **Integration Test** yang langsung menembak database asli.
