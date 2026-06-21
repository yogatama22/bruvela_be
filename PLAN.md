# Bruvela Bakehouse - Project Roadmap

> **Status terakhir diverifikasi**: Juni 2025
> **Pendekatan**: Iterative, dengan quick-win prioritized

---

## 📊 Audit Snapshot (Juni 2025)

### ✅ Yang Sudah Selesai (Live & Wired ke API)

| Komponen | Status | File / Endpoint | Catatan |
|---|---|---|---|
| Backend Go (12 handler, 11 repo, 11 model) | ✅ | `internal/**`, `cmd/main.go` | Clean architecture |
| Frontend Nuxt (24+ halaman) | ✅ | `bruvela_fe/pages/**` | Nuxt UI + Tailwind + ApexCharts |
| Auth + JWT + middleware | ✅ | `middleware/auth.ts` (BE), `plugins/auth.ts` (FE) | 7-day cookie |
| Dashboard live data | ✅ | `pages/index.vue` + `GET /dashboard/summary` | Sales chart per produk, status pie, recent orders, low stock |
| Finance page + Modal Tambah Jurnal | ✅ | `pages/finance/index.vue` + `GET/POST /finance/journal`, `GET /finance/summary` | Income/expense/balance cards + journal table |
| Inventory Create/Edit | ✅ | `pages/inventory/create.vue`, `[id]/edit.vue` | CRUD lengkap |
| Inventory Purchases page | ✅ | `pages/inventory/purchases.vue` + `GET/POST /ingredient-purchases` | Auto-add stock + auto-journal expense |
| Reports CSV export (4 endpoint) | ✅ | `pages/reports/index.vue` + `GET /reports/{orders,finance,inventory,hpp}` | Download via blob |
| Auth middleware FE | ✅ | `plugins/auth.ts` (global route middleware) | Redirect ke /login jika belum auth |
| Auto-add stock + journal saat purchase | ✅ | `internal/handler/purchase_handler.go` | Single transaction atomic |
| Recipe page list HPP | ✅ | `pages/recipes/index.vue` + `GET /recipes` | Link ke calculator |
| Orders CRUD + filter batch/status/channel | ✅ | `pages/orders/**` + `GET/POST/PATCH/DELETE /orders` | |
| Batches management | ✅ | `pages/batches/index.vue` + 8 endpoint batch | Auto-close logic |

### 🔴 Yang BELUM Selesai (Fase Lanjutan)

| Komponen | Status | Catatan |
|---|---|---|
| **Auto-deduct stok saat order → "diproses"** | 🔴 Belum | Order_handler.UpdateStatus() cuma update kolom, tidak trigger stock_log "out" + decrement ingredient.current_stock |
| **stock_logs handler/route/repository** | 🔴 Belum | Model + migration sudah ada, tapi tidak ada repository, handler, atau route di cmd/main.go |
| **Production Calculator (`/recipes/calculator`)** | 🔴 Belum | Backend endpoint `POST /recipes/calculator` belum ada. Frontend page belum ada (cuma link di recipes/index.vue) |
| **Batch Detail Page** | 🔴 Belum | Endpoint `GET /batches/:id/summary` belum ada. Halaman `pages/batches/[id].vue` belum ada |
| **Sales Chart API endpoint khusus** | 🟡 Parsial | Dashboard sudah build chart dari orders aggregation, tapi tidak ada endpoint `/dashboard/charts/sales` reusable |
| **Pinia stores full implementation** | 🔴 Stub | `stores/{orders,finance,inventory}.ts` masih kerangka kosong, action tidak implement |
| **Customer CRUD** | 🔴 Belum | Model customer ada tapi tidak ada handler/repository/routes/FE |
| **Supplier CRUD** | 🔴 Belum | Hanya field di ingredient_purchase.supplier, tidak ada entity terpisah |
| **Unit test backend** | 🔴 Belum | Tidak ada `*_test.go` |
| **Swagger API doc** | 🔴 Belum | |
| **Multi-partner modal/profit** | 🔴 Belum | Field partner sudah di journal_entry tapi belum ada UI/aggregation per partner |
| **Filter date range orders** | 🔴 Belum | Backend belum support start_date/end_date |

---

## 🎯 Rencana Lanjutan (4 Fase)

### 📦 FASE 2: Auto-Deduct Stok + Stock Log (PRIORITAS TINGGI)
**Goal**: Fitur pembeda dari Excel — stok benar-benar bergerak otomatis

#### 2.1 Backend - stock_logs infrastructure
- **File baru**: `internal/repository/stock_log_repository.go`
  - `Create(tx *gorm.DB, log *model.StockLog) error`
  - `FindAll(filters map[string]interface{}) ([]*model.StockLog, error)`
  - `FindByIngredient(ingredientID uuid.UUID) ([]*model.StockLog, error)`
  - `FindByBatch(batchID uuid.UUID) ([]*model.StockLog, error)`
- **File baru**: `internal/handler/stock_log_handler.go`
  - `GET /api/v1/stock-logs?batch_id=&ingredient_id=&log_type=&start_date=&end_date=`
- **Wire ke cmd/main.go**: Tambah `stockLogRepo` + `stockLogHandler` + route group

#### 2.2 Backend - Auto-deduct stok saat order diproses
- **File edit**: `internal/handler/order_handler.go` — `UpdateStatus()`
  - Detect transisi `*→diproses`: load order items, hitung kebutuhan bahan via `recipe × qty_box`
  - Validasi stok cukup (return 422 jika insufficient)
  - Single transaction:
    1. Decrement `ingredients.current_stock` per bahan
    2. Insert `stock_logs` (log_type="out", reference_type="order", reference_id=orderID)
- **Update cmd/main.go**: Inject `recipeRepo` + `ingredientRepo` ke `NewOrderHandler`

#### 2.3 Backend - Auto-add stock_log saat purchase (sudah ada auto-add stock + journal, tinggal tambah stock_log)
- **File edit**: `internal/handler/purchase_handler.go` — `Create()`
  - Tambah insert `stock_logs` (log_type="in", reference_type="purchase") dalam transaction

#### 2.4 Frontend - Halaman Stock Logs
- **File baru**: `bruvela_fe/pages/inventory/stock-logs.vue`
  - Filter by ingredient, batch, log_type, date range
  - Tabel: tanggal, bahan, tipe, qty, stock_before/after, reference, note
- **Update composable**: `useStockLogs()` di `bruvela_fe/composables/useStockLogs.ts`
- **Update menu**: Sidebar link "Stock Movement" di `layouts/default.vue`

#### 2.5 Frontend - Notifikasi success saat order diproses
- **File edit**: `bruvela_fe/pages/orders/[id]/index.vue`
  - Tampilkan ringkasan bahan yang terpotong jika status berubah ke diproses

**Estimasi total: 16-22 jam**

---

### 📦 FASE 3: Production Calculator + Batch Detail
**Goal**: Tool analisis yang bikin app ini "production-ready"

#### 3.1 Backend - Production Calculator
- **File edit**: `internal/handler/recipe_handler.go` — tambah `CalculateProduction()`
  - Input: `{ items: [{product_id, qty_box}], batch_id? }`
  - Output: `{ ingredients: [{id, name, needed, current_stock, status}], total_hpp }`
  - Logic: aggregate `recipe.qty_use × qty_box` per ingredient, bandingkan dengan current_stock
- **Wire route**: `POST /api/v1/recipes/calculator` (protected)

#### 3.2 Frontend - Calculator UI
- **File baru**: `bruvela_fe/pages/recipes/calculator.vue`
  - Form interaktif: pilih multiple produk + qty box
  - Real-time preview kebutuhan bahan
  - Visual indicator: 🟢 cukup, 🟡 pas-pasan, 🔴 kurang
  - Tombol "Simpan sebagai Plan Produksi" (optional future)

#### 3.3 Backend - Batch Summary endpoint
- **File baru**: `internal/handler/batch_summary_handler.go` atau tambah method di batch_handler
  - `GET /api/v1/batches/:id/summary`
  - Response: `{ batch, orders_total, revenue, total_paid, hpp_total, gross_profit, margin_pct, stock_usage }`

#### 3.4 Frontend - Batch Detail Page
- **File baru**: `bruvela_fe/pages/batches/[id].vue`
  - Header info batch + status (active/closed)
  - 4 cards: orders, revenue, profit, margin
  - Tabel order pada batch tersebut
  - Link ke laporan & finance
- **File edit**: `bruvela_fe/pages/batches/index.vue` — tambah tombol "Detail" per row

**Estimasi total: 14-20 jam**

---

### 📦 FASE 4: Polish & Portfolio
**Goal**: Siap showcase & deploy

| Task | Detail | Estimasi |
|---|---|---|
| Unit test backend | `*_test.go` untuk order, finance, recipe handler (table-driven) | 8-10 jam |
| API documentation | Swagger pakai `swaggo/swag` + integrasi ke `cmd/main.go` | 4-5 jam |
| Pinia stores full | Pindah dari stub ke proper state management dengan composable integration | 5-7 jam |
| Error handling global | Toast/error boundary di FE, panic recovery middleware di BE | 3-4 jam |
| Customer CRUD | Model + repo + handler + FE list/create/edit | 8-10 jam |
| Supplier CRUD | Sama seperti customer | 6-8 jam |
| Multi-partner analytics | Agregasi modal & profit per partner (Aul/Dhavinna) di finance page | 5-6 jam |
| Date range filter orders | Backend support start_date/end_date + FE date picker | 3-4 jam |
| Deploy Railway.app + Vercel | Setup env + build command + domain | 4-6 jam |
| README updated | Screenshot + fitur list + API doc link + ERD | 3-4 jam |

**Estimasi total: 49-64 jam (~2 minggu part-time)**

---

## 🚀 Quick-Win Prioritas (Mulai dari Sini)

Jika hanya punya waktu 1 minggu, eksekusi 5 ini dulu:

1. 🔴 **Backend auto-deduct stok di order UpdateStatus** (6-8 jam)
   - Impact paling tinggi: fitur pembeda dari Excel, owner langsung lihat
2. 🔴 **Backend stock_logs handler + route + repository** (3-4 jam)
   - Audit trail pergerakan stok, modal audit
3. 🔴 **Frontend stock logs page** (4-5 jam)
   - Visibilitas history stok
4. 🔴 **Production Calculator (BE + FE)** (5-6 jam)
   - Tool interaktif untuk owner & tim
5. 🟡 **Pinia stores full implementation** (4-5 jam)
   - Clean architecture, mudah maintain

**Total quick-win: ~22-28 jam = 1 minggu part-time**

---

## 📝 Catatan & Prinsip Kerja

- **Commit per fase**: Setiap fase di-commit terpisah dengan pesan Conventional Commits
- **Test manual**: Gunakan Postman/Insomnia atau curl untuk test endpoint backend sebelum integrasi FE
- **Migration safety**: Selalu backup database sebelum perubahan schema
- **Backward compatibility**: Jangan hapus endpoint lama, tambah baru dengan versioning (`/api/v2/...`)
- **Security**: Selalu validasi user_id via JWT middleware, jangan percaya client input

---

## 🔗 Referensi

- **API Base URL**: `http://localhost:8080/api/v1`
- **Frontend URL**: `http://localhost:3000`
- **Database**: PostgreSQL dengan UUID + GORM AutoMigrate
- **Login default**: `admin@bruvela.com` / `admin123`
- **ERD lengkap**: Lihat `bruvela_bakehouse_app_plan.txt`
- **Setup guide**: Lihat `SETUP.md`

---

_Last updated: Juni 2025 oleh audit session_