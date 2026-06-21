# Bruvela Bakehouse - Backend API

Backend API untuk sistem manajemen toko Bruvela Bakehouse yang dibangun dengan **Go (Golang)** menggunakan framework **Gin** dan **PostgreSQL**.

## 🚀 Tech Stack

- **Framework**: Gin (Go web framework)
- **Database**: PostgreSQL 12+
- **ORM**: GORM
- **Authentication**: JWT (JSON Web Token)
- **Password Hashing**: bcrypt
- **Validation**: go-playground/validator
- **CORS**: gin-contrib/cors
- **Environment**: godotenv

## 📋 Prerequisites

- Go 1.21 atau lebih tinggi
- PostgreSQL 12 atau lebih tinggi
- Git

## 🛠️ Installation

### 1. Clone Repository

```bash
cd bruvela-be
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Database

Buat database PostgreSQL:

```bash
createdb bruvela_db
```

Atau menggunakan psql:

```sql
CREATE DATABASE bruvela_db;
```

### 4. Import Schema

```bash
psql -U postgres -d bruvela_db -f migrations/schema.sql
```

### 5. Setup Environment Variables

Copy file `.env.example` menjadi `.env`:

```bash
cp .env.example .env
```

Edit file `.env` sesuai konfigurasi Anda:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=bruvela_db
DB_SSLMODE=disable

JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRE_HOURS=24

SERVER_PORT=8080
GIN_MODE=debug

CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

### 6. Create Admin User

Setelah import schema, Anda perlu membuat user admin. Gunakan script berikut atau buat manual:

```sql
-- Password: admin123 (hashed)
INSERT INTO users (name, email, password, role) VALUES
('Admin', 'admin@bruvela.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin');
```

## 🏃 Running the Application

### Development Mode

```bash
go run cmd/main.go
```

### Build and Run

```bash
go build -o bruvela-api cmd/main.go
./bruvela-api
```

Server akan berjalan di `http://localhost:8080`

## 📚 API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication

#### Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@bruvela.com",
  "password": "admin123"
}
```

Response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "name": "Admin",
    "email": "admin@bruvela.com",
    "role": "admin"
  }
}
```

#### Get Current User

```http
GET /api/v1/auth/me
Authorization: Bearer {token}
```

### Products (Menu)

```http
GET    /api/v1/products           # Get all products
GET    /api/v1/products/:id       # Get product by ID
POST   /api/v1/products           # Create new product
PUT    /api/v1/products/:id       # Update product
DELETE /api/v1/products/:id       # Delete product
GET    /api/v1/products/:id/recipe # Get product recipe
```

### Ingredients (Bahan Baku)

```http
GET    /api/v1/ingredients         # Get all ingredients
GET    /api/v1/ingredients/:id     # Get ingredient by ID
GET    /api/v1/ingredients/alerts  # Get low stock ingredients
POST   /api/v1/ingredients         # Create new ingredient
PUT    /api/v1/ingredients/:id     # Update ingredient
DELETE /api/v1/ingredients/:id     # Delete ingredient
```

### Orders (Pemesanan)

```http
GET    /api/v1/orders              # Get all orders (with filters)
GET    /api/v1/orders/:id          # Get order by ID
POST   /api/v1/orders              # Create new order
PATCH  /api/v1/orders/:id/status   # Update production status
PATCH  /api/v1/orders/:id/pay      # Update payment status
DELETE /api/v1/orders/:id          # Delete order
```

Query Parameters untuk GET /api/v1/orders:

- `batch_id` - Filter by batch ID
- `pay_status` - Filter by payment status (belum_bayar, dp, lunas)
- `prod_status` - Filter by production status (baru, diproses, siap_kirim, selesai, batal)
- `limit` - Limit results (default: 10)
- `offset` - Offset for pagination (default: 0)

### Example Requests

#### Create Product

```http
POST /api/v1/products
Authorization: Bearer {token}
Content-Type: application/json

{
  "code": "BRV-CLASSIC",
  "name": "Bruv Classic",
  "price": 20000,
  "pcs_per_box": 6,
  "status": "active"
}
```

#### Create Order

```http
POST /api/v1/orders
Authorization: Bearer {token}
Content-Type: application/json

{
  "batch_id": "uuid-batch-id",
  "customer_name": "John Doe",
  "order_date": "2024-01-15",
  "channel": "whatsapp",
  "shipping_type": "dalam_kota",
  "shipping_dest": "Jakarta Selatan",
  "shipping_cost": 15000,
  "discount": 0,
  "items": [
    {
      "product_id": "uuid-product-id",
      "qty_box": 2,
      "price_per_box": 20000
    }
  ]
}
```

## 🗂️ Project Structure

```
bruvela-be/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   └── config.go               # Configuration management
├── internal/
│   ├── handler/                # HTTP handlers/controllers
│   │   ├── auth_handler.go
│   │   ├── product_handler.go
│   │   ├── order_handler.go
│   │   └── ingredient_handler.go
│   ├── middleware/             # HTTP middleware
│   │   ├── auth.go
│   │   └── cors.go
│   ├── model/                  # Database models
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── ingredient.go
│   │   ├── recipe.go
│   │   ├── order.go
│   │   └── ...
│   └── repository/             # Data access layer
│       ├── user_repository.go
│       ├── product_repository.go
│       └── ...
├── pkg/
│   ├── auth/                   # Authentication utilities
│   │   ├── jwt.go
│   │   └── password.go
│   └── database/               # Database connection
│       └── database.go
├── migrations/
│   └── schema.sql              # PostgreSQL schema
├── .env.example                # Environment variables template
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## 🔐 Security

- Passwords di-hash menggunakan bcrypt
- JWT untuk autentikasi stateless
- CORS middleware untuk mengatur akses
- Input validation menggunakan go-playground/validator
- Prepared statements via GORM untuk mencegah SQL injection

## 📊 Database Schema

Database menggunakan UUID sebagai primary key dan memiliki 11 tabel utama:

1. **users** - User accounts
2. **products** - Product/menu master
3. **ingredients** - Ingredient master
4. **recipes** - Recipe composition
5. **customers** - Customer data
6. **batches** - Production batches
7. **orders** - Customer orders
8. **order_items** - Order line items
9. **ingredient_purchases** - Purchase transactions
10. **stock_logs** - Stock movement audit trail
11. **journal_entries** - Financial journal

Lihat file `migrations/schema.sql` untuk detail lengkap.

## 🧪 Testing

```bash
go test ./...
```

## 📦 Build for Production

```bash
# Build binary
go build -o bruvela-api cmd/main.go

# Set production environment
export GIN_MODE=release

# Run
./bruvela-api
```

## 🚀 Deployment

### Railway.app (Recommended)

1. Push code ke GitHub
2. Connect repository di Railway
3. Set environment variables
4. Deploy

### Docker (Optional)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o bruvela-api cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bruvela-api .
EXPOSE 8080
CMD ["./bruvela-api"]
```

## 📝 TODO / Future Enhancements

- [ ] Batch handler implementation
- [ ] Recipe calculator endpoint
- [ ] Dashboard analytics endpoints
- [ ] Excel/PDF export functionality
- [ ] Stock auto-deduction on order processing
- [ ] Email/WhatsApp notifications
- [ ] Unit tests
- [ ] API documentation with Swagger
- [ ] Rate limiting
- [ ] Logging middleware

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📄 License

This project is private and proprietary to Bruvela Bakehouse.

## 👥 Contact

- Developer: [Your Name]
- Email: [your.email@example.com]
- Project Link: [https://github.com/yourusername/bruvela-be](https://github.com/yourusername/bruvela-be)

---

**Note**: Jangan lupa untuk mengganti password default admin setelah deployment pertama!
