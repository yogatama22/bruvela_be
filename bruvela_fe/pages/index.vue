<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold text-gray-900">Dashboard</h2>
      <p class="mt-1 text-sm text-gray-500">
        <span v-if="loading">Memuat data batch aktif...</span>
        <span v-else-if="activeBatch">
          Batch #{{ activeBatch.batch_number }} - {{ activeBatch.name || 'Bruvela Bakehouse' }}
          <span v-if="activeBatch.start_date" class="text-xs text-gray-400 ml-1">
            ({{ formatDate(activeBatch.start_date) }})
          </span>
        </span>
        <span v-else>Belum ada batch aktif. <NuxtLink to="/batches" class="text-primary-600 underline">Buat batch baru</NuxtLink></span>
      </p>
    </div>

    <!-- Summary Cards -->
    <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Total Order</p>
            <p class="mt-2 text-3xl font-bold text-gray-900">{{ summary.total_orders || 0 }}</p>
            <p class="mt-1 text-xs text-gray-500">{{ summary.total_boxes || 0 }} box terjual</p>
          </div>
          <div class="p-3 bg-blue-100 rounded-lg">
            <UIcon name="i-heroicons-shopping-bag" class="w-6 h-6 text-blue-600" />
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Omzet</p>
            <p class="mt-2 text-3xl font-bold text-gray-900">{{ formatRupiahShort(summary.total_revenue) }}</p>
            <p class="mt-1 text-xs text-gray-500">{{ summary.total_orders || 0 }} pesanan</p>
          </div>
          <div class="p-3 bg-green-100 rounded-lg">
            <UIcon name="i-heroicons-banknotes" class="w-6 h-6 text-green-600" />
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Sudah Lunas</p>
            <p class="mt-2 text-3xl font-bold text-gray-900">{{ formatRupiahShort(summary.total_paid) }}</p>
            <p class="mt-1 text-xs text-amber-600">Pending: {{ formatRupiahShort(summary.total_pending) }}</p>
          </div>
          <div class="p-3 bg-purple-100 rounded-lg">
            <UIcon name="i-heroicons-credit-card" class="w-6 h-6 text-purple-600" />
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Stok Kritis</p>
            <p class="mt-2 text-3xl font-bold text-red-600">{{ summary.low_stock_count || 0 }}</p>
            <p class="mt-1 text-xs text-gray-500">bahan perlu restock</p>
          </div>
          <div class="p-3 bg-red-100 rounded-lg">
            <UIcon name="i-heroicons-exclamation-triangle" class="w-6 h-6 text-red-600" />
          </div>
        </div>
      </UCard>
    </div>

    <!-- Charts Row -->
    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900">Penjualan per Varian</h3>
            <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-arrow-path" @click="loadAll" :loading="loading" />
          </div>
        </template>
        <div v-if="loading && !salesChartSeries[0].data.length" class="flex justify-center py-12">
          <div class="text-center">
            <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
            <p class="mt-2 text-sm text-gray-500">Memuat data...</p>
          </div>
        </div>
        <ClientOnly v-else>
          <apexchart
            type="bar"
            height="300"
            :options="salesChartOptions"
            :series="salesChartSeries"
          />
        </ClientOnly>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900">Status Order</h3>
            <UBadge color="blue" variant="subtle" size="xs">{{ summary.total_orders || 0 }} total</UBadge>
          </div>
        </template>
        <div v-if="!summary.status_count || Object.keys(summary.status_count).length === 0" class="text-center py-12 text-sm text-gray-500">
          Belum ada data order
        </div>
        <div v-else class="space-y-4">
          <div
            v-for="(count, status) in summary.status_count"
            :key="status"
            class="flex items-center justify-between"
          >
            <div class="flex items-center space-x-3">
              <div :class="['w-3 h-3 rounded-full', getStatusDotColor(status)]"></div>
              <span class="text-sm font-medium text-gray-700">{{ formatStatusLabel(status) }}</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-sm font-bold text-gray-900">{{ count }}</span>
              <span class="text-xs text-gray-400">({{ getPercentage(count) }}%)</span>
            </div>
          </div>
        </div>
      </UCard>
    </div>

    <!-- Critical Stock + Recent Orders Row -->
    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900">Alert Stok Kritis</h3>
            <UBadge color="red" variant="solid">{{ summary.low_stock_count || 0 }} items</UBadge>
          </div>
        </template>
        <div v-if="!summary.low_stock_ingredients || summary.low_stock_ingredients.length === 0" class="text-center py-12">
          <UIcon name="i-heroicons-check-circle" class="w-12 h-12 text-green-400 mx-auto mb-2" />
          <p class="text-sm text-gray-500">Semua stok aman 🎉</p>
        </div>
        <div v-else class="space-y-2 max-h-80 overflow-y-auto">
          <div
            v-for="item in summary.low_stock_ingredients.slice(0, 8)"
            :key="item.id"
            class="flex items-center justify-between p-3 bg-red-50 rounded-lg"
          >
            <div>
              <p class="text-sm font-medium text-gray-900">{{ item.name }}</p>
              <p class="text-xs text-gray-500">
                Stok: {{ item.current_stock }} {{ item.use_unit }} / Min: {{ item.min_stock }} {{ item.use_unit }}
              </p>
            </div>
            <UBadge color="red" variant="solid" size="xs">RESTOCK</UBadge>
          </div>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900">Order Terbaru</h3>
            <UButton to="/orders" color="gray" variant="ghost" size="xs">
              Lihat Semua
              <UIcon name="i-heroicons-arrow-right" class="ml-1 w-4 h-4" />
            </UButton>
          </div>
        </template>
        <div v-if="loading && recentOrders.length === 0" class="flex justify-center py-12">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        </div>
        <div v-else-if="recentOrders.length === 0" class="text-center py-12 text-sm text-gray-500">
          Belum ada order
        </div>
        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead>
              <tr>
                <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase">Customer</th>
                <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase">Tanggal</th>
                <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase">Total</th>
                <th class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="order in recentOrders" :key="order.id" class="hover:bg-gray-50">
                <td class="px-3 py-3 whitespace-nowrap">
                  <div class="text-sm font-medium text-gray-900">{{ order.customer_name }}</div>
                  <div class="text-xs text-gray-500">{{ order.channel }}</div>
                </td>
                <td class="px-3 py-3 whitespace-nowrap text-xs text-gray-500">{{ formatDate(order.order_date) }}</td>
                <td class="px-3 py-3 whitespace-nowrap text-sm font-medium text-gray-900">
                  Rp {{ (order.total_bill || 0).toLocaleString('id-ID') }}
                </td>
                <td class="px-3 py-3 whitespace-nowrap">
                  <UBadge :color="getStatusColor(order.prod_status)" variant="subtle" size="xs">
                    {{ formatStatusLabel(order.prod_status) }}
                  </UBadge>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Dashboard'
})

const toast = useToast()
const { fetchActiveBatch } = useBatches()
const { fetchSummary } = useDashboard()
const { fetchOrders } = useOrders()

const loading = ref(false)
const activeBatch = ref<any>(null)
const summary = ref<any>({
  total_orders: 0,
  total_boxes: 0,
  total_revenue: 0,
  total_paid: 0,
  total_pending: 0,
  low_stock_count: 0,
  status_count: {},
  low_stock_ingredients: []
})
const recentOrders = ref<any[]>([])

// === Sales Chart (by product variant) ===
const salesChartSeries = ref([{ name: 'Penjualan', data: [] as number[] }])
const salesChartOptions = reactive({
  chart: {
    type: 'bar',
    toolbar: { show: false },
    fontFamily: 'Inter, sans-serif'
  },
  plotOptions: {
    bar: { borderRadius: 8, horizontal: false, columnWidth: '60%' }
  },
  dataLabels: { enabled: false },
  colors: ['#4f46e5'],
  xaxis: {
    categories: [] as string[],
    labels: { style: { fontSize: '12px', fontFamily: 'Inter, sans-serif' } }
  },
  yaxis: {
    title: {
      text: 'Box Terjual',
      style: { fontSize: '12px', fontFamily: 'Inter, sans-serif' }
    }
  },
  grid: { borderColor: '#f3f4f6' }
})

// === Helpers ===
const formatRupiahShort = (value: number) => {
  if (!value) return 'Rp 0'
  if (value >= 1_000_000) return `Rp ${(value / 1_000_000).toFixed(1)}M`
  if (value >= 1_000) return `Rp ${(value / 1_000).toFixed(0)}k`
  return `Rp ${value.toLocaleString('id-ID')}`
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      year: 'numeric'
    })
  } catch {
    return dateStr
  }
}

const formatStatusLabel = (status: string) => {
  const map: Record<string, string> = {
    baru: 'Baru',
    proses: 'Proses',
    siap_kirim: 'Siap Kirim',
    selesai: 'Selesai',
    batal: 'Batal'
  }
  return map[status] || status
}

const getStatusColor = (status: string) => {
  const map: Record<string, string> = {
    baru: 'blue',
    proses: 'yellow',
    siap_kirim: 'purple',
    selesai: 'green',
    batal: 'red'
  }
  return map[status] || 'gray'
}

const getStatusDotColor = (status: string) => {
  const map: Record<string, string> = {
    baru: 'bg-blue-500',
    proses: 'bg-yellow-500',
    siap_kirim: 'bg-purple-500',
    selesai: 'bg-green-500',
    batal: 'bg-red-500'
  }
  return map[status] || 'bg-gray-500'
}

const getPercentage = (count: number) => {
  const total = summary.value.total_orders || 0
  if (total === 0) return 0
  return Math.round((count / total) * 100)
}

// === Data Loading ===
const loadAll = async () => {
  loading.value = true
  try {
    // 1. Get active batch first
    const batchRes = await fetchActiveBatch()
    if (batchRes.error || !batchRes.data) {
      activeBatch.value = null
      summary.value = {
        total_orders: 0,
        total_boxes: 0,
        total_revenue: 0,
        total_paid: 0,
        total_pending: 0,
        low_stock_count: 0,
        status_count: {},
        low_stock_ingredients: []
      }
      recentOrders.value = []
      return
    }
    activeBatch.value = batchRes.data
    const batchID = (batchRes.data as any).id

    // 2. Parallel fetch: dashboard summary + recent orders
    const [summaryRes, ordersRes] = await Promise.all([
      fetchSummary(undefined, batchID),
      fetchOrders({ batchId: batchID })
    ])

    if (summaryRes.error) {
      toast.add({
        title: 'Error',
        description: 'Gagal memuat summary dashboard',
        color: 'red'
      })
    } else if (summaryRes.data) {
      const data = summaryRes.data as any
      // Compute total boxes from orders
      const allOrders = (ordersRes.data as any)?.data || []
      const totalBoxes = allOrders.reduce((sum: number, o: any) => {
        return sum + ((o.Items || o.items || []).reduce((s: number, i: any) => s + (i.QtyBox || i.qty_box || 0), 0))
      }, 0)
      summary.value = {
        ...data,
        total_boxes: totalBoxes
      }
    }

    if (ordersRes.error) {
      toast.add({
        title: 'Error',
        description: 'Gagal memuat data order',
        color: 'red'
      })
    } else {
      const allOrders = (ordersRes.data as any)?.data || []
      recentOrders.value = allOrders.slice(0, 5)

      // Build sales chart by aggregating order items per product
      const productMap = new Map<string, { name: string; qty: number }>()
      allOrders.forEach((order: any) => {
        ;(order.Items || order.items || []).forEach((item: any) => {
          const key = item.product_id || item.ProductID
          const name = item.product_name || item.ProductName || 'Unknown'
          const qty = item.qty_box || item.QtyBox || 0
          if (productMap.has(key)) {
            productMap.get(key)!.qty += qty
          } else {
            productMap.set(key, { name, qty })
          }
        })
      })
      const sorted = Array.from(productMap.values()).sort((a, b) => b.qty - a.qty).slice(0, 8)
      salesChartSeries.value = [{ name: 'Penjualan', data: sorted.map(p => p.qty) }]
      salesChartOptions.xaxis.categories = sorted.map(p => p.name)
    }
  } catch (err) {
    console.error('Dashboard load error:', err)
    toast.add({
      title: 'Error',
      description: 'Terjadi kesalahan saat memuat dashboard',
      color: 'red'
    })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadAll()
})
</script>