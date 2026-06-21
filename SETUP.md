# Setup Guide - Bruvela Bakehouse Backend

## 📋 Checklist Struktur Proyek

Struktur proyek sudah sesuai dengan `bruvela_bakehouse_app_plan.txt` line 545:

### ✅ Struktur yang Sudah Dibuat

```
backend/
├── cmd/
│   └── main.go                     ✅
├── internal/
│   ├── handler/
│   │   ├── auth_handler.go         ✅ (auth.go)
│   │   ├── order_handler.go        ✅ (order.go)
│   │   ├── product_handler.go      ✅ (product.go)
│   │   ├── ingredient_handler.go   ✅ (ingredient.go)
│   │   ├── finance_handler.go      ✅ (finance.go)
│   │   └── dashboard_handler.go    ✅ (dashboard.go)
│   ├── repository/
│   │   ├── order_repository.go     ✅ (order.go)
│   │   ├── ingredient_repository.go ✅ (ingredient.go)
│   │   ├── product_repository.go   ✅
│   │   ├── recipe_repository.go    ✅
│   │   ├── batch_repository.go     ✅
│   │   ├── user_repository.go      ✅
│   │   └── finance_repository.go   ✅ (finance.go)
│   ├── model/
│   │   ├── order.go                ✅
│   │   ├── product.go              ✅
│   │   ├── ingredient.go           ✅
│   │   ├── user.go                 ✅
│   │   ├── batch.go                ✅
│   │   ├── customer.go             ✅
│   │   ├── recipe.go               ✅
│   │   ├── order_item.go           ✅
│   │   ├── ingredient_purchase.go  ✅
│   │   ├── stock_log.go            ✅
│   │   └── journal_entry.go        ✅ (finance.go)
│   └── middleware/
│       ├── auth.go                 ✅
│       └── cors.go                 ✅
├── pkg/
│   ├── auth/                       ✅ (JWT helper)
│   │   ├── jwt.go
│   │   └── password.go
│   └── database/                   ✅
│       └── database.go
├── migrations/                     ✅ (SQL files)
│   ├── schema.sql
│   └── seed.sql
└── config/
    └── config.go                   ✅
```

**Catatan**: Beberapa file menggunakan suffix untuk clarity (e.g., `auth_handler.go` instead of `auth.go`), yang merupakan best practice Go.

## 🚀 Langkah Setup

### 1. Install PostgreSQL

Jika belum ada, install PostgreSQL:

**macOS (Homebrew):**
```bash
brew install postgresql@15
brew services start postgresql@15
```

**Atau gunakan ServBay/Docker/PostgreSQL.app**

### 2. Buat Database

```bash
# Login ke PostgreSQL
psql -U postgres

# Buat database
CREATE DATABASE bruvela_db;

# Buat user jika perlu
CREATE USER bruvela WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE bruvela_db TO bruvela;

# Keluar
\q
```

### 3. Import Schema

```bash
# Import schema utama
psql -U postgres -d bruvela_db -f migrations/schema.sql

# Import sample data (opsional)
psql -U postgres -d bruvela_db -f migrations/seed.sql
```

### 4. Konfigurasi Environment

Edit file `.env` sesuai konfigurasi PostgreSQL Anda:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_actual_password  # <-- Ganti ini!
DB_NAME=bruvela_db
DB_SSLMODE=disable

JWT_SECRET=bruvela-super-secret-jwt-key-2024
JWT_EXPIRE_HOURS=24

SERVER_PORT=8080
GIN_MODE=debug

CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

### 5. Jalankan Aplikasi

```bash
# Cara 1: Langsung dengan go run
go run cmd/main.go

# Cara 2: Build dulu, lalu run
go build -o bin/bruvela-api cmd/main.go
./bin/bruvela-api

# Cara 3: Menggunakan Makefile
make run
```

Server akan berjalan di: **http://localhost:8080**

### 6. Test API

**Health Check:**
```bash
curl http://localhost:8080/health
```

**Login (gunakan data dari seed.sql):**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@bruvela.com",
    "password": "admin123"
  }'
```

Response akan berisi token JWT yang bisa digunakan untuk endpoint lainnya.

## 🔍 Verifikasi Build

Aplikasi sudah di-build dan tidak ada error kompilasi:

```bash
✅ go mod tidy       # Dependencies downloaded
✅ go build          # Build successful
✅ No compile errors
```

## 📊 Status Fitur

### ✅ Fitur yang Sudah Dibuat

1. **Authentication**
   - JWT login
   - Password hashing (bcrypt)
   - Protected routes

2. **Products (Menu)**
   - CRUD operations
   - Recipe management

3. **Ingredients (Bahan Baku)**
   - CRUD operations
   - Low stock alerts
   - Stock tracking

4. **Orders (Pemesanan)**
   - CRUD operations
   - Status tracking (production & payment)
   - Filtering & pagination
   - Order items support

5. **Finance (Keuangan)**
   - Journal entries
   - Batch summary
   - Income/expense tracking

6. **Dashboard**
   - Summary statistics
   - Low stock alerts
   - Order status counts

### 🔄 Fitur yang Bisa Ditambahkan Nanti

- [ ] Service layer (untuk business logic yang kompleks)
- [ ] Batch handler (CRUD batch)
- [ ] Recipe calculator
- [ ] Auto stock deduction
- [ ] Excel/PDF export
- [ ] Email notifications
- [ ] Unit tests

## 🛠️ Troubleshooting

### Error: "Failed to connect to database"

**Solusi:**
1. Pastikan PostgreSQL running: `brew services list` atau `pg_ctl status`
2. Cek password di `.env` sudah benar
3. Cek database sudah dibuat: `psql -U postgres -l`

### Error: "Port already in use"

**Solusi:**
```bash
# Cari process yang menggunakan port 8080
lsof -i :8080

# Kill process tersebut
kill -9 <PID>

# Atau ganti port di .env
SERVER_PORT=8081
```

### Error: "Table does not exist"

**Solusi:**
```bash
# Import schema lagi
psql -U postgres -d bruvela_db -f migrations/schema.sql
```

## 📚 API Endpoints

Lihat `README.md` untuk dokumentasi lengkap API endpoints.

## ✅ Kesimpulan

- ✅ Struktur proyek sesuai dengan plan
- ✅ Semua handler sudah dibuat (auth, product, ingredient, order, finance, dashboard)
- ✅ Repository layer lengkap
- ✅ Middleware (auth, cors) sudah ada
- ✅ Database schema PostgreSQL siap
- ✅ Build berhasil tanpa error
- ⚠️ Perlu setup database PostgreSQL untuk running

**Next Step:** Setup PostgreSQL dan jalankan aplikasi!
