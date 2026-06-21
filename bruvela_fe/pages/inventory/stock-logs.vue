<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Stock Movement</h2>
        <p class="mt-1 text-sm text-gray-500">Riwayat pergerakan stok bahan baku</p>
      </div>
    </div>

    <UCard>
      <div class="flex flex-col sm:flex-row gap-4 mb-6">
        <div class="flex-1">
          <USelectMenu
            v-model="selectedIngredient"
            :options="ingredientOptions"
            value-attribute="value"
            placeholder="Semua Bahan"
            size="lg"
          />
        </div>
        <USelectMenu
          v-model="selectedBatch"
          :options="batchOptions"
          value-attribute="value"
            placeholder="Semua Batch"
          size="lg"
        />
        <USelectMenu
          v-model="selectedLogType"
          :options="logTypeOptions"
          value-attribute="value"
          placeholder="Semua Tipe"
          size="lg"
        />
      </div>

      <div class="overflow-x-auto">
        <div v-if="loading" class="flex flex-col items-center justify-center py-12">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mb-3"></div>
          <p class="text-sm text-gray-500">Memuat data...</p>
        </div>

        <div v-else-if="!stockLogs || stockLogs.length === 0" class="flex flex-col items-center justify-center py-12">
          <UIcon name="i-heroicons-inbox" class="w-12 h-12 text-gray-400 mb-3" />
          <p class="text-sm font-medium text-gray-900 mb-1">Tidak ada data stock log</p>
          <p class="text-xs text-gray-500">Coba ubah filter pencarian</p>
        </div>

        <table v-else class="min-w-full divide-y divide-gray-200">
          <thead>
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tanggal</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Bahan</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Batch</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tipe</th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Qty</th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Stok Sebelum</th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Stok Sesudah</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Referensi</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Catatan</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="log in stockLogs" :key="log.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDateTime(log.created_at) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ log.ingredient?.name || '-' }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ log.batch?.name || log.batch?.batch_name || '-' }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <UBadge :color="log.log_type === 'in' ? 'green' : 'red'" variant="subtle" size="xs">
                  {{ log.log_type === 'in' ? 'MASUK' : 'KELUAR' }}
                </UBadge>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{{ log.qty?.toFixed(3) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-500">{{ log.stock_before?.toFixed(3) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-right text-gray-900">{{ log.stock_after?.toFixed(3) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <UBadge v-if="log.reference_type" :color="log.reference_type === 'order' ? 'blue' : 'purple'" variant="subtle" size="xs">
                  {{ log.reference_type }}
                </UBadge>
                <span v-else>-</span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ log.note || '-' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Stock Movement'
})

const toast = useToast()
const { fetchStockLogs } = useStockLogs()
const { fetchIngredients } = useIngredients()
const { fetchBatches } = useBatches()

const loading = ref(false)
const stockLogs = ref<any[]>([])
const ingredients = ref<any[]>([])
const batches = ref<any[]>([])

const selectedIngredient = ref('')
const selectedBatch = ref('')
const selectedLogType = ref('')

const logTypeOptions = [
  { label: 'Semua Tipe', value: '' },
  { label: 'Masuk', value: 'in' },
  { label: 'Keluar', value: 'out' }
]

const ingredientOptions = computed(() => {
  return [
    { label: 'Semua Bahan', value: '' },
    ...ingredients.value.map((i: any) => ({ label: i.name, value: i.id }))
  ]
})

const batchOptions = computed(() => {
  return [
    { label: 'Semua Batch', value: '' },
    ...batches.value.map((b: any) => ({ label: b.batch_name || b.name, value: b.id }))
  ]
})

const loadStockLogs = async () => {
  loading.value = true
  const filters: any = {}
  if (selectedIngredient.value) filters.ingredientId = selectedIngredient.value
  if (selectedBatch.value) filters.batchId = selectedBatch.value
  if (selectedLogType.value) filters.logType = selectedLogType.value

  const { data, error } = await fetchStockLogs(filters)
  loading.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal memuat data stock log',
      color: 'red'
    })
    return
  }

  stockLogs.value = data || []
}

const loadIngredients = async () => {
  const { data } = await fetchIngredients()
  ingredients.value = data || []
}

const loadBatches = async () => {
  const { data } = await fetchBatches()
  batches.value = data?.data || data || []
}

const formatDateTime = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  loadIngredients()
  loadBatches()
  loadStockLogs()
})

watch([selectedIngredient, selectedBatch, selectedLogType], () => {
  loadStockLogs()
})
</script>
