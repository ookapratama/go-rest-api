# Makefile for Go REST API

# Load .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: run-db run-app stop migrate migrate-down supabase-migrate migrate-create help

## migrate-create: Membuat file migrasi baru (Gunakan: make migrate-create name=nama_migrasi)
migrate-create:
	docker run --rm -v $(PWD)/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq $(name)

## help: Menampilkan daftar perintah yang tersedia
help:
	@echo "Perintah yang bisa dijalankan:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run-db: Menjalankan database Postgres di Docker
run-db:
	docker-compose up -d db

## run-app: Menjalankan aplikasi Go dan Database
run-app:
	docker-compose up --build

## stop: Menghentikan semua kontainer Docker
stop:
	docker-compose down

## migrate: Menjalankan migrasi database (membuat tabel) di Docker
migrate:
	docker-compose run --rm migrate

## migrate-down: Membatalkan migrasi database (menghapus tabel)
migrate-down:
	docker-compose run --rm migrate -path=/migrations/ -database="postgres://postgres:password@db:5432/product_db?sslmode=disable" down 1

## supabase-migrate: Menjalankan migrasi ke database Supabase (Pastikan DB_URL_SUPABASE di .env sudah diisi)
supabase-migrate:
	@if [ "$(DB_URL_SUPABASE)" = "" ]; then \
		echo "Error: DB_URL_SUPABASE tidak ditemukan di .env"; \
		exit 1; \
	fi
	docker run --rm -v $(PWD)/migrations:/migrations migrate/migrate -path=/migrations/ -database "$(DB_URL_SUPABASE)" up
