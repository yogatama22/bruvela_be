# Integrasi Frontend-Backend Bruvela Bakehouse

## Status Integrasi

### ✅ Sudah Terintegrasi

#### 0. Authentication
- **Login** (`/login`) - JWT-based authentication
- **Auto-protect** - Semua routes kecuali login
- **User info** - Display di header
- **Logout** - Clear token & redirect

**Test Account:**
- Email: `admin@bruvela.com`
- Password: `admin123`

**API Endpoints:**
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/me` - Get current user

#### 1. Orders (Pemesanan)
- **List Orders** (`/orders`) - Fetch, filter, search, delete
- **Create Order** (`/orders/create`) - Create order dengan multiple items
- **View Order** (`/orders/[id]`) - Detail order lengkap
- **Edit Order** (`/orders/[id]/edit`) - Update order (coming soon)

**API Endpoints:**
- `GET /api/v1/orders` - List orders dengan filter
- `POST /api/v1/orders` - Create order baru
- `GET /api/v1/orders/:id` - Get order detail
- `PUT /api/v1/orders/:id` - Update order
- `PATCH /api/v1/orders/:id/status` - Update status order
- `PATCH /api/v1/orders/:id/pay` - Update status pembayaran
- `DELETE /api/v1/orders/:id` - Delete order

#### 2. Inventory (Stok Bahan)
- **List Ingredients** (`/inventory`) - Fetch, filter, search, delete
- **Summary Cards** - Total bahan, stok aman, stok kritis (real-time)

**API Endpoints:**
- `GET /api/v1/ingredients` - List semua bahan
- `GET /api/v1/ingredients/:id` - Get detail bahan
- `GET /api/v1/ingredients/alerts` - Get bahan dengan stok kritis
- `POST /api/v1/ingredients` - Create bahan baru
- `PUT /api/v1/ingredients/:id` - Update bahan
- `DELETE /api/v1/ingredients/:id` - Delete bahan

### 🔄 Belum Terintegrasi (Masih Dummy Data)

#### 3. Recipes & Menu
- List produk dengan HPP
- Kalkulator produksi

**API Endpoints yang perlu digunakan:**
- `GET /api/v1/products`
- `GET /api/v1/products/:id/recipe`
- `POST /api/v1/recipes/calculator`

#### 4. Finance (Keuangan)
- Summary keuangan
- Laporan per batch
- Jurnal keuangan

**API Endpoints yang perlu digunakan:**
- `GET /api/v1/batches`
- `GET /api/v1/journal`
- `GET /api/v1/finance/summary`
- `GET /api/v1/finance/profit-loss`

#### 5. Dashboard
- Summary cards
- Charts (sales, revenue)

**API Endpoints yang perlu digunakan:**
- `GET /api/v1/dashboard/summary`
- `GET /api/v1/dashboard/charts/sales`
- `GET /api/v1/dashboard/charts/finance`

## Composables yang Tersedia

### 0. `useAuth()`
Authentication & user management

```typescript
const { 
  login,
  logout,
  getMe,
  isAuthenticated,
  user,
  token
} = useAuth()

// Login
const { data, error } = await login('admin@bruvela.com', 'admin123')

// Logout
logout()

// Check auth
if (isAuthenticated.value) {
  console.log('User:', user.value)
}
```

### 1. `useApi()`
Base API client dengan authentication

```typescript
const api = useApi()
const { data, error } = await api.get('/endpoint')
const { data, error } = await api.post('/endpoint', body)
```

### 2. `useOrders()`
```typescript
const { 
  fetchOrders,
  fetchOrderById,
  createOrder,
  updateOrder,
  updateOrderStatus,
  updatePaymentStatus,
  deleteOrder
} = useOrders()
```

### 3. `useProducts()`
```typescript
const {
  fetchProducts,
  fetchProductById,
  fetchProductRecipe,
  createProduct,
  updateProduct,
  deleteProduct
} = useProducts()
```

### 4. `useIngredients()`
```typescript
const {
  fetchIngredients,
  fetchIngredientById,
  fetchCriticalStock,
  createIngredient,
  updateIngredient,
  deleteIngredient,
  fetchPurchases,
  createPurchase
} = useIngredients()
```

### 5. `useFinance()`
```typescript
const {
  fetchJournal,
  createJournalEntry,
  fetchSummary,
  fetchProfitLoss,
  fetchBatches,
  createBatch,
  closeBatch
} = useFinance()
```

### 6. `useDashboard()`
```typescript
const {
  fetchSummary,
  fetchSalesChart,
  fetchFinanceChart
} = useDashboard()
```

## Konfigurasi Backend

Edit file `.env`:

```env
NUXT_PUBLIC_API_BASE=http://localhost:8080/api/v1
```

Untuk production:
```env
NUXT_PUBLIC_API_BASE=https://api.bruvela.com/api/v1
```

## Format Data Backend

### Order
```typescript
{
  id: string (UUID)
  customer_name: string
  order_date: string (ISO date)
  channel: 'whatsapp' | 'instagram' | 'offline' | 'titip_teman'
  shipping_type: 'dalam_kota' | 'luar_kota' | 'jnt'
  shipping_dest: string
  shipping_cost: number
  discount: number
  total_product: number
  total_bill: number
  pay_status: 'belum_bayar' | 'dp' | 'lunas'
  prod_status: 'baru' | 'diproses' | 'siap_kirim' | 'selesai' | 'batal'
  note: string
  items: OrderItem[]
}
```

### Ingredient
```typescript
{
  id: string (UUID)
  name: string
  pack_unit: string
  qty_per_pack: number
  use_unit: string
  price_per_pack: number
  price_per_use: number
  min_stock: number
  current_stock: number
}
```

## Testing

1. **Start Backend**
```bash
cd /Users/yogatama.egiantoro/Documents/Development/bruvela-be
make run
```

2. **Start Frontend**
```bash
cd bruvela_fe
npm run dev
```

3. **Login**
- Buka http://localhost:3000
- Auto-redirect ke `/login`
- Login: `admin@bruvela.com` / `admin123`
- Redirect ke dashboard

4. **Test Endpoints**
- Dashboard: http://localhost:3000/
- Orders: http://localhost:3000/orders
- Inventory: http://localhost:3000/inventory
- Recipes: http://localhost:3000/recipes
- Finance: http://localhost:3000/finance

## Error Handling

Semua composables mengembalikan format:
```typescript
{
  data: any | null,
  error: any | null
}
```

Contoh penggunaan:
```typescript
const { data, error } = await fetchOrders()

if (error) {
  toast.add({
    title: 'Error',
    description: 'Gagal memuat data',
    color: 'red'
  })
  return
}

// Use data
orders.value = data || []
```

## Next Steps

1. ✅ Authentication - DONE
2. ✅ Integrasi Orders - DONE
3. ✅ Integrasi Inventory - DONE
4. ⏳ Integrasi Recipes
5. ⏳ Integrasi Finance
6. ⏳ Integrasi Dashboard
7. ⏳ Add Edit Order page
8. ⏳ Add Create/Edit Ingredient pages
9. ⏳ Add Register/Forgot Password pages
