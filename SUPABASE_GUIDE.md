# Panduan Koneksi Go REST API ke Supabase

Supabase menyediakan PostgreSQL managed. Karena aplikasi ini menggunakan driver `pgx/v5`, koneksinya sangat mudah.

## 1. Ambil Connection String dari Supabase

1. Masuk ke [Dashboard Supabase](https://app.supabase.com/).
2. Pilih Proyek Anda.
3. Pergi ke menu **Settings** (Ikon Gear) -> **Database**.
4. Scroll ke bawah sampai ketemu bagian **Connection string**.
5. Pilih tab **URI**.
6. Salin string tersebut. Contoh:
   `postgresql://postgres:[PASSWORD]@db.xxxx.supabase.co:5432/postgres`

## 2. Update File `.env`

Buka file `.env` di root proyek Anda dan ganti nilai `DATABASE_URL`:

```env
PORT=8080
DATABASE_URL=postgresql://postgres:PASSWORD_ANDA_DI_SINI@db.kjprgocvbdbelagskjqi.supabase.co:5432/postgres
```

_Catatan: Pastikan mengganti `[PASSWORD]` dengan password database yang Anda buat saat pertama kali setup proyek Supabase._

## 3. Jalankan Database Migration

Karena Supabase Anda masih kosong, Anda harus membuat tabel `products` di sana menggunakan Docker:

```bash
docker run --rm -v $(pwd)/migrations:/migrations migrate/migrate \
    -path=/migrations/ \
    -database "postgresql://postgres:PASSWORD_ANDA@db.xxxx.supabase.co:5432/postgres" \
    up
```

## 4. Cara Cek Koneksi di Kode Go (Tips Interview)

Di file `main.go`, kita menggunakan fungsi `Ping()` untuk memastikan aplikasi benar-benar terhubung ke Supabase sebelum server online:

```go
// Di dalam main.go
dbPool, err := pgxpool.New(context.Background(), dbURL)
if err != nil {
    log.Fatal("Gagal membuat pool:", err)
}

// Ini adalah baris kritikal untuk cek koneksi ke Supabase
if err := dbPool.Ping(context.Background()); err != nil {
    log.Fatal("Supabase Tidak Terjangkau:", err)
}
slog.Info("Koneksi Supabase Berhasil!")
```

## Tips: Troubleshooting Koneksi

- **Password Salah:** Jika muncul error `password authentication failed`, cek kembali password di .env.
- **Network/Firewall:** Supabase secara default mengizinkan koneksi dari mana saja, tapi pastikan internet Anda tidak memblokir port `5432`.
- **SSL Mode:** Jika ada masalah SSL, tambahkan `?sslmode=disable` (untuk dev) atau `?sslmode=require` (untuk prod) di ujung URL.
