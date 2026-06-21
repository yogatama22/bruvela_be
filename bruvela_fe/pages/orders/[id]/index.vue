<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-4">
        <UButton
          to="/orders"
          color="gray"
          variant="ghost"
          icon="i-heroicons-arrow-left"
          size="lg"
        />
        <div>
          <h2 class="text-2xl font-bold text-gray-900">Detail Order</h2>
          <p class="mt-1 text-sm text-gray-500">Order #{{ order?.id?.substring(0, 8) }}</p>
        </div>
      </div>
      <div class="flex gap-3">
        <UButton
          :to="`/orders/${route.params.id}/edit`"
          color="gray"
          variant="outline"
          icon="i-heroicons-pencil"
          size="lg"
        >
          Edit
        </UButton>
        <UButton
          @click="handlePrint"
          color="gray"
          variant="outline"
          icon="i-heroicons-printer"
          size="lg"
        >
          Cetak Nota
        </UButton>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        <p class="mt-2 text-sm text-gray-500">Memuat data...</p>
      </div>
    </div>

    <div v-else-if="order" class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <div class="lg:col-span-2 space-y-6">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900">Informasi Customer</h3>
          </template>
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="text-sm font-medium text-gray-600">Nama Customer</label>
                <p class="mt-1 text-base text-gray-900">{{ order.customer_name }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Channel</label>
                <p class="mt-1 text-base text-gray-900">{{ order.channel }}</p>
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="text-sm font-medium text-gray-600">Tanggal Order</label>
                <p class="mt-1 text-base text-gray-900">{{ formatDate(order.order_date) }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Tipe Pengiriman</label>
                <p class="mt-1 text-base text-gray-900">{{ order.shipping_type || '-' }}</p>
              </div>
            </div>
            <div v-if="order.shipping_dest">
              <label class="text-sm font-medium text-gray-600">Alamat Tujuan</label>
              <p class="mt-1 text-base text-gray-900">{{ order.shipping_dest }}</p>
            </div>
            <div v-if="order.note">
              <label class="text-sm font-medium text-gray-600">Catatan</label>
              <p class="mt-1 text-base text-gray-900">{{ order.note }}</p>
            </div>
          </div>
        </UCard>

        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900">Produk</h3>
          </template>
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead>
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Produk</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Qty</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Harga</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Subtotal</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <tr v-for="item in orderItems" :key="item.id">
                  <td class="px-4 py-3 text-sm text-gray-900">{{ item.product_name }}</td>
                  <td class="px-4 py-3 text-sm text-right text-gray-900">{{ item.qty_box }} box</td>
                  <td class="px-4 py-3 text-sm text-right text-gray-900">Rp {{ item.price_per_box?.toLocaleString('id-ID') }}</td>
                  <td class="px-4 py-3 text-sm text-right font-medium text-gray-900">Rp {{ item.subtotal?.toLocaleString('id-ID') }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </UCard>
      </div>

      <div class="space-y-6">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900">Status</h3>
          </template>
          <div class="space-y-4">
            <div>
              <label class="text-sm font-medium text-gray-600">Status Order</label>
              <div class="mt-2 flex items-center gap-3">
                <UBadge :color="getOrderStatusColor(order.prod_status)" variant="subtle" size="lg">
                  {{ formatStatus(order.prod_status) }}
                </UBadge>
                <UButton
                  v-if="order.prod_status !== 'proses' && order.prod_status !== 'selesai' && order.prod_status !== 'batal'"
                  size="xs"
                  color="yellow"
                  variant="outline"
                  @click="handleProcessOrder"
                  :loading="processing"
                >
                  Proses Order
                </UButton>
              </div>
              <div v-if="stockDeducted" class="mt-3 p-3 bg-green-50 rounded-lg">
                <p class="text-sm font-medium text-green-800 mb-1">Stok Bahan Terpotong:</p>
                <ul class="text-xs text-green-700 space-y-1">
                  <li v-for="(d, i) in stockDeducted" :key="i">
                    {{ d.ingredient_name }}: -{{ d.qty_deducted?.toFixed(3) }} (sisa: {{ d.stock_after?.toFixed(3) }})
                  </li>
                </ul>
              </div>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-600">Status Pembayaran</label>
              <div class="mt-2">
                <UBadge :color="getPaymentStatusColor(order.pay_status)" variant="subtle" size="lg">
                  {{ formatStatus(order.pay_status) }}
                </UBadge>
              </div>
            </div>
          </div>
        </UCard>

        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900">Ringkasan Pembayaran</h3>
          </template>
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Subtotal Produk</span>
              <span class="font-medium text-gray-900">Rp {{ (order.total_product || 0).toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Ongkir</span>
              <span class="font-medium text-gray-900">Rp {{ (order.shipping_cost || 0).toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Diskon</span>
              <span class="font-medium text-red-600">- Rp {{ (order.discount || 0).toLocaleString('id-ID') }}</span>
            </div>
            <div class="pt-3 border-t border-gray-200">
              <div class="flex justify-between">
                <span class="text-base font-semibold text-gray-900">Total</span>
                <span class="text-lg font-bold text-primary-600">Rp {{ (order.total_bill || 0).toLocaleString('id-ID') }}</span>
              </div>
            </div>
          </div>
        </UCard>
      </div>
    </div>

    <div v-else class="text-center py-12">
      <p class="text-gray-500">Order tidak ditemukan</p>
      <UButton to="/orders" class="mt-4">Kembali ke Daftar Order</UButton>
    </div>

    <UModal v-model="stockErrorModal">
      <UCard>
        <template #header>
          <div class="flex items-center gap-3">
            <div class="p-2 bg-red-100 rounded-lg">
              <UIcon name="i-heroicons-exclamation-triangle" class="w-6 h-6 text-red-600" />
            </div>
            <h3 class="text-lg font-semibold text-gray-900">Stok Tidak Cukup</h3>
          </div>
        </template>
        <div class="space-y-3">
          <p class="text-sm text-gray-600">Order tidak dapat diproses karena bahan baku berikut kurang:</p>
          <div class="max-h-72 overflow-y-auto rounded-lg border border-gray-200">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Bahan</th>
                  <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Butuh</th>
                  <th class="px-4 py-2 text-right text-xs font-medium text-gray-500 uppercase">Stok</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <tr v-for="(item, i) in stockErrorItems" :key="i">
                  <td class="px-4 py-2 text-sm font-medium text-gray-900">{{ item.name }}</td>
                  <td class="px-4 py-2 text-sm text-right text-red-600 font-bold">{{ item.needed }}</td>
                  <td class="px-4 py-2 text-sm text-right text-gray-500">{{ item.stock }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <template #footer>
          <div class="flex justify-end">
            <UButton color="red" @click="stockErrorModal = false">Tutup</UButton>
          </div>
        </template>
      </UCard>
    </UModal>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Detail Order'
})

const route = useRoute()
const toast = useToast()
const { fetchOrderById, updateOrderStatus } = useOrders()

const loading = ref(true)
const order = ref<any>(null)
const orderItems = ref<any[]>([])
const processing = ref(false)
const stockDeducted = ref<any[] | null>(null)
const stockErrorModal = ref(false)
const stockErrorItems = ref<{name: string, needed: string, stock: string}[]>([])

const loadOrder = async () => {
  loading.value = true
  const { data, error } = await fetchOrderById(route.params.id as string)
  loading.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal memuat detail order',
      color: 'red'
    })
    return
  }

  order.value = data
  orderItems.value = data?.items || []
}

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

const formatStatus = (status: string) => {
  if (!status) return '-'
  const statusMap: Record<string, string> = {
    'baru': 'Baru',
    'proses': 'Proses',
    'siap_kirim': 'Siap Kirim',
    'selesai': 'Selesai',
    'batal': 'Batal',
    'belum_bayar': 'Belum Bayar',
    'dp': 'DP',
    'lunas': 'Lunas'
  }
  return statusMap[status] || status
}

const getOrderStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    'baru': 'blue',
    'proses': 'yellow',
    'siap_kirim': 'purple',
    'selesai': 'green',
    'batal': 'red'
  }
  return colors[status] || 'gray'
}

const getPaymentStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    'belum_bayar': 'red',
    'dp': 'yellow',
    'lunas': 'green'
  }
  return colors[status] || 'gray'
}

const handlePrint = () => {
  window.print()
}

const handleProcessOrder = async () => {
  processing.value = true
  const { data, error } = await updateOrderStatus(route.params.id as string, 'proses')
  processing.value = false

  if (error) {
    const errorMsg = typeof error === 'object' && error.error ? error.error : 'Gagal memproses order'
    if (errorMsg.includes('stok tidak cukup')) {
      const items = errorMsg.match(/\[([^\]]+)\]/)
      if (items) {
        stockErrorItems.value = items[1].split(', ').map((item: string) => {
          const match = item.match(/(.+): butuh ([\d.]+), stok ([\d.]+)/)
          if (match) {
            return { name: match[1].trim(), needed: match[2], stock: match[3] }
          }
          return { name: item, needed: '-', stock: '-' }
        })
        stockErrorModal.value = true
        return
      }
    }
    toast.add({
      title: 'Error',
      description: errorMsg,
      color: 'red',
      timeout: 8000
    })
    return
  }

  if (data?.stock_deducted) {
    stockDeducted.value = data.stock_deducted
    toast.add({
      title: 'Order Diproses',
      description: 'Stok bahan berhasil dipotong',
      color: 'green',
      timeout: 8000
    })
  } else {
    toast.add({
      title: 'Berhasil',
      description: 'Status order berhasil diupdate',
      color: 'green'
    })
  }

  loadOrder()
}

onMounted(() => {
  loadOrder()
})
</script>
