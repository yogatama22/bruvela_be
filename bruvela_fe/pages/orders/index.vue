<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Pemesanan</h2>
        <p class="mt-1 text-sm text-gray-500">Kelola semua order pelanggan</p>
      </div>
      <UButton to="/orders/create" icon="i-heroicons-plus" size="lg">
        Order Baru
      </UButton>
    </div>

    <UCard>
      <div class="flex flex-col sm:flex-row gap-4 mb-6">
        <div class="flex-1">
          <UInput
            v-model="search"
            icon="i-heroicons-magnifying-glass"
            placeholder="Cari customer, produk..."
            size="lg"
          />
        </div>
        <USelectMenu
          v-model="selectedBatch"
          :options="batchOptions"
          placeholder="Batch"
          size="lg"
        />
        <USelectMenu
          v-model="selectedStatus"
          :options="statusOptions"
          placeholder="Status Order"
          size="lg"
        />
        <USelectMenu
          v-model="selectedPayment"
          :options="paymentOptions"
          placeholder="Status Bayar"
          size="lg"
        />
      </div>

      <div class="overflow-x-auto">
        <!-- Loading State -->
        <div v-if="loading" class="flex flex-col items-center justify-center py-12">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mb-3"></div>
          <p class="text-sm text-gray-500">Memuat data...</p>
        </div>

        <!-- Empty State -->
        <div v-else-if="!filteredOrders || filteredOrders.length === 0" class="flex flex-col items-center justify-center py-12">
          <UIcon name="i-heroicons-inbox" class="w-12 h-12 text-gray-400 mb-3" />
          <p class="text-sm font-medium text-gray-900 mb-1">Tidak ada data order</p>
          <p class="text-xs text-gray-500">{{ search || selectedStatus || selectedPayment ? 'Coba ubah filter pencarian' : 'Belum ada order yang dibuat' }}</p>
        </div>

        <!-- Data Table -->
        <table v-else class="min-w-full divide-y divide-gray-200">
          <thead>
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Order ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Batch</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Customer</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tanggal</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Produk</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status Order</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status Bayar</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="order in filteredOrders" :key="order.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-500">
                <UTooltip v-if="order.pay_status === 'lunas'" :text="order.id" color="green" variant="solid">
                  <span class="text-green-600 font-medium">#{{ order.id?.substring(0, 8) }}</span>
                </UTooltip>
                <span v-else>#{{ order.id?.substring(0, 8) }}</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ order.batch?.batch_name || order.batch?.name || '-' }}</div>
                <div class="text-xs text-gray-500">{{ formatDate(order.batch?.start_date) }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ order.customer_name }}</div>
                <div class="text-xs text-gray-500">{{ order.channel }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDate(order.order_date) }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ order.total_product || 0 }} item</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">Rp {{ (order.total_bill || 0).toLocaleString('id-ID') }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <UDropdown :items="getOrderStatusActions(order)" :popper="{ placement: 'bottom-start' }">
                  <UButton :color="getOrderStatusColor(order.prod_status)" variant="subtle" size="xs">
                    {{ formatStatus(order.prod_status) }}
                    <template #trailing>
                      <UIcon name="i-heroicons-chevron-down" class="w-3 h-3 ml-1" />
                    </template>
                  </UButton>
                </UDropdown>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <UDropdown :items="getPaymentStatusActions(order)" :popper="{ placement: 'bottom-start' }">
                  <UButton :color="getPaymentStatusColor(order.pay_status)" variant="subtle" size="xs">
                    {{ formatStatus(order.pay_status) }}
                    <template #trailing>
                      <UIcon name="i-heroicons-chevron-down" class="w-3 h-3 ml-1" />
                    </template>
                  </UButton>
                </UDropdown>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <UDropdown :items="getOrderActions(order)">
                  <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal" />
                </UDropdown>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="flex items-center justify-between mt-6">
        <div class="text-sm text-gray-500">
          Menampilkan {{ filteredOrders.length }} dari {{ orders.length }} order
        </div>
        <div class="flex gap-2">
          <UButton color="gray" variant="outline" size="sm" icon="i-heroicons-chevron-left" disabled>
            Previous
          </UButton>
          <UButton color="gray" variant="outline" size="sm" icon-trailing="i-heroicons-chevron-right">
            Next
          </UButton>
        </div>
      </div>
    </UCard>

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
  title: 'Pemesanan'
})

const toast = useToast()
const { fetchOrders, deleteOrder, updateOrderStatus, updatePaymentStatus } = useOrders()
const { fetchShippingTypes } = useShippingTypes()
const { fetchPaymentStatuses } = usePaymentStatuses()
const { fetchOrderStatuses } = useOrderStatuses()
const { fetchBatches } = useBatches()

const search = ref('')
const selectedStatus = ref(null)
const selectedPayment = ref(null)
const selectedBatch = ref(null)
const loading = ref(false)
const orders = ref<any[]>([])
const shippingTypes = ref<any[]>([])
const paymentStatuses = ref<any[]>([])
const orderStatuses = ref<any[]>([])
const batches = ref<any[]>([])

const statusOptions = computed(() => {
  return [{ label: 'Semua', value: null }, ...orderStatuses.value.map((s: any) => ({ label: s.order_status_name, value: s.order_status_code }))]
})

const paymentOptions = computed(() => {
  return [{ label: 'Semua', value: null }, ...paymentStatuses.value.map((s: any) => ({ label: s.status_name, value: s.status_code }))]
})

const batchOptions = computed(() => {
  return [{ label: 'Semua', value: null }, ...batches.value.map((b: any) => ({ label: b.batch_name || b.name, value: b.id }))]
})

const loadOrders = async () => {
  loading.value = true
  const filters: any = {}

  if (selectedStatus.value?.value) {
    filters.prodStatus = selectedStatus.value.value
  }
  if (selectedPayment.value?.value) {
    filters.payStatus = selectedPayment.value.value
  }
  if (selectedBatch.value?.value) {
    filters.batchId = selectedBatch.value.value
  }
  if (search.value) {
    filters.search = search.value
  }

  const { data, error } = await fetchOrders(filters)
  loading.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal memuat data order',
      color: 'red'
    })
    return
  }

  orders.value = data?.data || []
}

const filteredOrders = computed(() => {
  let result = orders.value

  // Filter by search (client-side)
  if (search.value) {
    result = result.filter((order: any) =>
      order.customer_name?.toLowerCase().includes(search.value.toLowerCase()) ||
      order.note?.toLowerCase().includes(search.value.toLowerCase())
    )
  }

  return result
})

const handleDelete = async (orderId: string) => {
  if (!confirm('Yakin ingin menghapus order ini?')) return

  const { error } = await deleteOrder(orderId)
  
  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal menghapus order',
      color: 'red'
    })
    return
  }

  toast.add({
    title: 'Berhasil',
    description: 'Order berhasil dihapus',
    color: 'green'
  })

  loadOrders()
}

onMounted(() => {
  loadOrders()
  loadShippingTypes()
  loadPaymentStatuses()
  loadOrderStatuses()
  loadBatches()
})

const loadShippingTypes = async () => {
  const { data, error } = await fetchShippingTypes()
  if (error) {
    console.error('Failed to load shipping types:', error)
    return
  }
  shippingTypes.value = data || []
}

const loadPaymentStatuses = async () => {
  const { data, error } = await fetchPaymentStatuses()
  if (error) {
    console.error('Failed to load payment statuses:', error)
    return
  }
  paymentStatuses.value = data || []
}

const loadOrderStatuses = async () => {
  const { data, error } = await fetchOrderStatuses()
  if (error) {
    console.error('Failed to load order statuses:', error)
    return
  }
  orderStatuses.value = data || []
}

const loadBatches = async () => {
  const { data, error } = await fetchBatches()
  if (error) {
    console.error('Failed to load batches:', error)
    return
  }
  batches.value = data?.data || data || []
  if (batches.value.length > 0) {
    selectedBatch.value = { label: batches.value[0].batch_name || batches.value[0].name, value: batches.value[0].id }
  }
}

watch([selectedStatus, selectedPayment, selectedBatch], () => {
  loadOrders()
})

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
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

const getOrderActions = (order: any) => {
  return [
    [{
      label: 'Lihat Detail',
      icon: 'i-heroicons-eye',
      click: () => navigateTo(`/orders/${order.id}`)
    }],
    [{
      label: 'Invoice',
      icon: 'i-heroicons-document-text',
      click: () => navigateTo(`/orders/${order.id}/invoice`)
    }],
    [{
      label: 'Edit Order',
      icon: 'i-heroicons-pencil',
      click: () => navigateTo(`/orders/${order.id}/edit`)
    }],
    [{
      label: 'Hapus',
      icon: 'i-heroicons-trash',
      click: () => handleDelete(order.id)
    }]
  ]
}

const getOrderStatusActions = (order: any) => {
  return [orderStatuses.value.map((status: any) => ({
    label: status.order_status_name,
    click: () => handleUpdateOrderStatus(order.id, status.order_status_code)
  }))]
}

const getPaymentStatusActions = (order: any) => {
  return [paymentStatuses.value.map((status: any) => ({
    label: status.status_name,
    click: () => handleUpdatePaymentStatus(order.id, status.status_code)
  }))]
}

const stockErrorModal = ref(false)
const stockErrorItems = ref<{name: string, needed: string, stock: string}[]>([])

const handleUpdateOrderStatus = async (orderId: string, status: string) => {
  const { data, error } = await updateOrderStatus(orderId, status)
  
  if (error) {
    if (status === 'proses' && error.error && error.error.includes('stok tidak cukup')) {
      const errorMsg = error.error
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
      description: 'Gagal update status order',
      color: 'red'
    })
    return
  }

  if (status === 'proses' && data?.stock_deducted) {
    const deducted = data.stock_deducted.map((d: any) => 
      `${d.ingredient_name}: -${d.qty_deducted?.toFixed(3)} (sisa: ${d.stock_after?.toFixed(3)})`
    ).join(', ')
    toast.add({
      title: 'Stok Berhasil Dipotong',
      description: deducted,
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

  loadOrders()
}

const handleUpdatePaymentStatus = async (orderId: string, payStatus: string) => {
  const { error } = await updatePaymentStatus(orderId, payStatus)
  
  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal update status pembayaran',
      color: 'red'
    })
    return
  }

  toast.add({
    title: 'Berhasil',
    description: 'Status pembayaran berhasil diupdate',
    color: 'green'
  })

  loadOrders()
}
</script>
