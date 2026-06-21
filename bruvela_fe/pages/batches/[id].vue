<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-4">
        <UButton to="/batches" color="gray" variant="ghost" icon="i-heroicons-arrow-left" size="lg" />
        <div>
          <h2 class="text-2xl font-bold text-gray-900">Batch Detail</h2>
          <p class="mt-1 text-sm text-gray-500" v-if="batch">Batch #{{ batch.batch_number }} - {{ batch.name || 'Tanpa Nama' }}</p>
        </div>
      </div>
      <div class="flex gap-3" v-if="batch">
        <UBadge :color="batch.status === 'open' ? 'green' : 'gray'" variant="solid" size="lg">
          {{ batch.status === 'open' ? 'AKTIF' : 'TUTUP' }}
        </UBadge>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        <p class="mt-2 text-sm text-gray-500">Memuat data...</p>
      </div>
    </div>

    <div v-else-if="summary" class="space-y-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <UCard>
          <div class="flex items-center gap-3">
            <div class="p-3 bg-blue-50 rounded-lg">
              <UIcon name="i-heroicons-shopping-bag" class="w-6 h-6 text-blue-600" />
            </div>
            <div>
              <p class="text-sm text-gray-600">Total Order</p>
              <p class="text-2xl font-bold text-gray-900">{{ summary.orders_total }}</p>
            </div>
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center gap-3">
            <div class="p-3 bg-green-50 rounded-lg">
              <UIcon name="i-heroicons-banknotes" class="w-6 h-6 text-green-600" />
            </div>
            <div>
              <p class="text-sm text-gray-600">Revenue</p>
              <p class="text-2xl font-bold text-gray-900">Rp {{ (summary.revenue || 0).toLocaleString('id-ID') }}</p>
            </div>
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center gap-3">
            <div class="p-3 rounded-lg" :class="summary.gross_profit >= 0 ? 'bg-primary-50' : 'bg-red-50'">
              <UIcon name="i-heroicons-chart-bar" class="w-6 h-6" :class="summary.gross_profit >= 0 ? 'text-primary-600' : 'text-red-600'" />
            </div>
            <div>
              <p class="text-sm text-gray-600">Gross Profit</p>
              <p class="text-2xl font-bold" :class="summary.gross_profit >= 0 ? 'text-primary-600' : 'text-red-600'">
                Rp {{ (summary.gross_profit || 0).toLocaleString('id-ID', { maximumFractionDigits: 0 }) }}
              </p>
            </div>
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center gap-3">
            <div class="p-3 bg-purple-50 rounded-lg">
              <UIcon name="i-heroicons-percent" class="w-6 h-6 text-purple-600" />
            </div>
            <div>
              <p class="text-sm text-gray-600">Margin</p>
              <p class="text-2xl font-bold text-gray-900">{{ (summary.margin_pct || 0).toFixed(1) }}%</p>
            </div>
          </div>
        </UCard>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900">Ringkasan Keuangan</h3>
          </template>
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Total dibayar (Lunas)</span>
              <span class="font-medium text-green-600">Rp {{ (summary.total_paid || 0).toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Total HPP</span>
              <span class="font-medium text-gray-900">Rp {{ (summary.hpp_total || 0).toLocaleString('id-ID', { maximumFractionDigits: 0 }) }}</span>
            </div>
            <div class="flex justify-between text-sm" v-if="summary.finance">
              <span class="text-gray-600">Pemasukan (Journal)</span>
              <span class="font-medium text-green-600">Rp {{ (summary.finance.total_income || 0).toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between text-sm" v-if="summary.finance">
              <span class="text-gray-600">Pengeluaran (Journal)</span>
              <span class="font-medium text-red-600">Rp {{ (summary.finance.total_expense || 0).toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between text-sm" v-if="summary.finance">
              <span class="text-gray-600">Saldo Journal</span>
              <span class="font-bold text-gray-900">Rp {{ (summary.finance.balance || 0).toLocaleString('id-ID') }}</span>
            </div>
          </div>
        </UCard>

        <div class="lg:col-span-2">
          <UCard>
            <template #header>
              <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold text-gray-900">Order dalam Batch</h3>
                <div class="flex gap-2">
                  <UButton
                    :to="`/reports?batch_id=${batch?.id}`"
                    color="gray"
                    variant="outline"
                    size="xs"
                    icon="i-heroicons-document-chart-bar"
                  >
                    Laporan
                  </UButton>
                  <UButton
                    :to="`/finance?batch_id=${batch?.id}`"
                    color="gray"
                    variant="outline"
                    size="xs"
                    icon="i-heroicons-banknotes"
                  >
                    Finance
                  </UButton>
                </div>
              </div>
            </template>
            <div class="overflow-x-auto" v-if="summary.orders && summary.orders.length > 0">
              <table class="min-w-full divide-y divide-gray-200">
                <thead>
                  <tr>
                    <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tanggal</th>
                    <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Customer</th>
                    <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                    <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-200">
                  <tr v-for="order in summary.orders" :key="order.id" class="hover:bg-gray-50 cursor-pointer" @click="navigateTo(`/orders/${order.id}`)">
                    <td class="px-4 py-3 text-sm text-gray-500">{{ formatDate(order.order_date) }}</td>
                    <td class="px-4 py-3 text-sm font-medium text-gray-900">{{ order.customer_name || '-' }}</td>
                    <td class="px-4 py-3">
                      <UBadge :color="getOrderStatusColor(order.prod_status)" variant="subtle" size="xs">
                        {{ formatStatus(order.prod_status) }}
                      </UBadge>
                    </td>
                    <td class="px-4 py-3 text-sm text-right font-medium text-gray-900">Rp {{ (order.total_bill || 0).toLocaleString('id-ID') }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-else class="text-center py-8">
              <p class="text-sm text-gray-500">Belum ada order dalam batch ini</p>
            </div>
          </UCard>
        </div>
      </div>
    </div>

    <div v-else class="text-center py-12">
      <p class="text-gray-500">Batch tidak ditemukan</p>
      <UButton to="/batches" class="mt-4">Kembali ke Daftar Batch</UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Batch Detail'
})

const route = useRoute()
const toast = useToast()
const { fetchBatchSummary } = useBatches()

const loading = ref(true)
const batch = ref<any>(null)
const summary = ref<any>(null)

const loadSummary = async () => {
  loading.value = true
  const { data, error } = await fetchBatchSummary(route.params.id as string)
  loading.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal memuat detail batch',
      color: 'red'
    })
    return
  }

  summary.value = data
  batch.value = data?.batch
}

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
    'batal': 'Batal'
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

onMounted(() => {
  loadSummary()
})
</script>
