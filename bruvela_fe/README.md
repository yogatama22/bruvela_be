# Bruvela Bakehouse - Frontend

Aplikasi manajemen toko Bruvela Bakehouse menggunakan Nuxt.js 3 dengan UI modern dan responsif.

## Tech Stack

- **Nuxt.js 3** - Framework Vue.js dengan SSR
- **Nuxt UI** - Komponen UI modern berbasis Tailwind CSS
- **Pinia** - State management
- **ApexCharts** - Library untuk grafik dan visualisasi data
- **TypeScript** - Type safety
- **Tailwind CSS** - Utility-first CSS framework

## Fitur

### 0. Authentication
- Login dengan JWT token
- Protected routes
- Auto-redirect jika belum login
- User info di header
- Logout functionality

**Test Account:**
- Email: `admin@bruvela.com`
- Password: `admin123`

### 1. Dashboard
- Summary cards (Total Order, Omzet, Lunas, Estimasi Laba)
- Grafik penjualan per varian (Bar Chart)
- Trend pendapatan per batch (Area Chart)
- Status order real-time
- Alert stok kritis
- Tabel order terbaru

### 2. Pemesanan
- List semua order dengan filter dan search
- Buat order baru dengan multiple produk
- Tracking status order dan pembayaran
- Detail order per customer

### 3. Resep & Menu
- Katalog produk dengan HPP dan margin
- Analisis profitabilitas per produk
- Kalkulator produksi (coming soon)

### 4. Inventory
- Monitoring stok bahan baku real-time
- Alert stok kritis dan minus
- Filter berdasarkan status stok
- Pencatatan pembelian bahan

### 5. Keuangan
- Summary keuangan (Modal, Pendapatan, Pengeluaran, Laba)
- Laporan per batch
- Jurnal keuangan
- Grafik trend keuangan

### 6. Laporan
- Export data ke Excel dan PDF
- Laporan pemesanan, keuangan, inventory, HPP

## Setup

### Install Dependencies

```bash
npm install
```

### Development Server

```bash
npm run dev
```

Aplikasi akan berjalan di `http://localhost:3000`

**Login:**
- Email: `admin@bruvela.com`
- Password: `admin123`

### Build untuk Production

```bash
npm run build
```

### Preview Production Build

```bash
npm run preview
```

## Konfigurasi

Edit file `.env` untuk konfigurasi API backend:

```env
NUXT_PUBLIC_API_BASE=http://localhost:8080/api/v1
```

## Struktur Folder

```
bruvela_fe/
├── assets/
│   └── css/
│       └── main.css          # Global styles & Tailwind
├── layouts/
│   └── default.vue           # Layout utama dengan sidebar
├── pages/
│   ├── index.vue             # Dashboard
│   ├── orders/
│   │   ├── index.vue         # List orders
│   │   └── create.vue        # Buat order baru
│   ├── recipes/
│   │   └── index.vue         # Menu & HPP
│   ├── inventory/
│   │   └── index.vue         # Stok bahan
│   ├── finance/
│   │   └── index.vue         # Keuangan
│   └── reports/
│       └── index.vue         # Laporan
├── stores/
│   ├── orders.ts             # Pinia store untuk orders
│   ├── inventory.ts          # Pinia store untuk inventory
│   └── finance.ts            # Pinia store untuk finance
├── app.vue
├── nuxt.config.ts
└── package.json
```

## Fitur UI/UX

- **Smooth Page Transitions** - Transisi halus antar halaman
- **Responsive Design** - Tampilan optimal di desktop, tablet, dan mobile
- **Modern Dashboard** - Inspirasi dari template Nuxt Charts
- **Inter Font** - Typography yang clean dan modern
- **Color-coded Status** - Status visual yang jelas dengan warna
- **Sidebar Navigation** - Navigasi yang mudah dengan collapsible sidebar

## Data Dummy

Aplikasi menggunakan data dummy berdasarkan data aktual Bruvela Bakehouse Batch 3:
- 72 total order
- 119 box terjual
- Rp 2,8M omzet
- 8 bahan dengan stok kritis
- 9 varian produk

## Integrasi Backend

Untuk menghubungkan dengan backend Golang, update fungsi di Pinia stores untuk memanggil API endpoints:

```typescript
// stores/orders.ts
async fetchOrders() {
  const config = useRuntimeConfig()
  const { data } = await useFetch(`${config.public.apiBase}/orders`)
  this.orders = data.value
}
```

## Deployment

Aplikasi dapat di-deploy ke:
- **Vercel** (Recommended) - Zero config deployment
- **Netlify** - Static site hosting
- **Railway** - Full-stack deployment

```bash
# Deploy ke Vercel
npm run build
vercel --prod
```

## License

Private - Bruvela Bakehouse © 2024
