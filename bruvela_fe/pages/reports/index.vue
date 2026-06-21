<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Laporan</h2>
        <p class="mt-1 text-sm text-gray-500">
          Export data ke CSV
          <span v-if="activeBatch" class="ml-1 text-gray-400">
            — Batch #{{ activeBatch.batch_number }}
          </span>
        </p>
      </div>
      <UButton color="gray" variant="ghost" icon="i-heroicons-arrow-path" size="sm" @click="loadBatch" :loading="loading" />
    </div>

    <div v-if="!loading && !activeBatch" class="text-center py-16 bg-white rounded-lg border border-dashed border-gray-300">
      <UIcon name="i-heroicons-document-chart-bar" class="w-12 h-12 text-gray-300 mx-auto mb-3" />
      <p class="text-sm font-medium text-gray-900 mb-1">Belum ada batch aktif</p>
      <p class="text-xs text-gray-500 mb-4">Buat batch terlebih dahulu untuk mulai generate laporan</p>
      <UButton to="/batches" icon="i-heroicons-plus">Buat Batch</UButton>
    </div>

    <div v-else class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
      <UCard v-for="report in reports" :key="report.key">
        <div class="space-y-4">
          <div :class="['flex items-center justify-center w-12 h-12 rounded-lg', report.bgColor]">
            <UIcon :name="report.icon" :class="['w-6 h-6', report.iconColor]" />
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">{{ report.title }}</h3>
            <p class="mt-1 text-sm text-gray-500">{{ report.description }}</p>
          </div>
          <UButton
            :color="report.color"
            variant="outline"
            block
            :loading="downloading === report.key"
            :disabled="!!downloading"
            @click="downloadReport(report)"
          >
            <UIcon name="i-heroicons-arrow-down-tray" class="mr-2" />
            {{ downloading === report.key ? 'Mengunduh...' : report.buttonLabel }}
          </UButton>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({ title: 'Laporan' })

const toast = useToast()
const { fetchActiveBatch } = useBatches()

const loading = ref(false)
const downloading = ref<string | null>(null)
const activeBatch = ref<any>(null)
const token = useCookie('auth_token')
const config = useRuntimeConfig()

const reports = [
  {
    key: 'orders',
    title: 'Laporan Pemesanan',
    description: 'Export data order ke CSV',
    icon: 'i-heroicons-shopping-bag',
    color: 'blue',
    bgColor: 'bg-blue-100',
    iconColor: 'text-blue-600',
    buttonLabel: 'Download CSV',
    requiresBatch: true,
    endpoint: '/reports/orders'
  },
  {
    key: 'finance',
    title: 'Laporan Keuangan',
    description: 'Export jurnal keuangan',
    icon: 'i-heroicons-banknotes',
    color: 'green',
    bgColor: 'bg-green-100',
    iconColor: 'text-green-600',
    buttonLabel: 'Download CSV',
    requiresBatch: true,
    endpoint: '/reports/finance'
  },
  {
    key: 'inventory',
    title: 'Laporan Inventory',
    description: 'Export data stok bahan',
    icon: 'i-heroicons-cube',
    color: 'purple',
    bgColor: 'bg-purple-100',
    iconColor: 'text-purple-600',
    buttonLabel: 'Download CSV',
    requiresBatch: false,
    endpoint: '/reports/inventory'
  },
  {
    key: 'hpp',
    title: 'Laporan HPP',
    description: 'Export analisis HPP per produk',
    icon: 'i-heroicons-book-open',
    color: 'amber',
    bgColor: 'bg-amber-100',
    iconColor: 'text-amber-600',
    buttonLabel: 'Download CSV',
    requiresBatch: false,
    endpoint: '/reports/hpp'
  }
]

const loadBatch = async () => {
  loading.value = true
  try {
    const { data, error } = await fetchActiveBatch()
    if (error || !data) {
      activeBatch.value = null
    } else {
      activeBatch.value = data
    }
  } finally {
    loading.value = false
  }
}

const downloadReport = async (report: any) => {
  if (report.requiresBatch && !activeBatch.value) {
    toast.add({
      title: 'Info',
      description: 'Belum ada batch aktif. Buat batch terlebih dahulu.',
      color: 'yellow'
    })
    return
  }

  downloading.value = report.key
  try {
    let url = `${config.public.apiBase}${report.endpoint}`
    if (report.requiresBatch) {
      url += `?batch_id=${activeBatch.value.id}`
    }

    const response: any = await $fetch(url, {
      headers: { Authorization: `Bearer ${token.value}` },
      responseType: 'blob'
    })

    // Create blob link and trigger download
    const blob = response instanceof Blob ? response : new Blob([response])
    const downloadUrl = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = downloadUrl
    a.download = `${report.key}_${new Date().toISOString().split('T')[0]}.csv`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(downloadUrl)

    toast.add({
      title: 'Berhasil',
      description: `${report.title} berhasil diunduh`,
      color: 'green'
    })
  } catch (err: any) {
    console.error('Download error:', err)
    toast.add({
      title: 'Gagal',
      description: 'Gagal mengunduh laporan. Pastikan backend berjalan.',
      color: 'red'
    })
  } finally {
    downloading.value = null
  }
}

onMounted(() => {
  loadBatch()
})
</script>
