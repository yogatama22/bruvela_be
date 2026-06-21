<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Inventory</h2>
        <p class="mt-1 text-sm text-gray-500">Kelola stok bahan baku</p>
      </div>
      <div class="flex gap-3">
        <UButton to="/inventory/purchases" color="gray" variant="outline" icon="i-heroicons-shopping-cart" size="lg">
          Pembelian Bahan
        </UButton>
        <UButton icon="i-heroicons-plus" size="lg" @click="navigateTo('/inventory/create')">
          Tambah Bahan
        </UButton>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-6 sm:grid-cols-3">
      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Total Bahan</p>
            <p class="mt-2 text-3xl font-bold text-gray-900">{{ ingredients.length }}</p>
          </div>
          <div class="p-3 bg-blue-100 rounded-lg">
            <UIcon name="i-heroicons-cube" class="w-6 h-6 text-blue-600" />
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Stok Aman</p>
            <p class="mt-2 text-3xl font-bold text-green-600">{{ safeCount }}</p>
          </div>
          <div class="p-3 bg-green-100 rounded-lg">
            <UIcon name="i-heroicons-check-circle" class="w-6 h-6 text-green-600" />
          </div>
        </div>
      </UCard>

      <UCard>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-600">Stok Kritis</p>
            <p class="mt-2 text-3xl font-bold text-red-600">{{ criticalCount }}</p>
          </div>
          <div class="p-3 bg-red-100 rounded-lg">
            <UIcon name="i-heroicons-exclamation-triangle" class="w-6 h-6 text-red-600" />
          </div>
        </div>
      </UCard>
    </div>

    <UCard>
      <div class="flex flex-col sm:flex-row gap-4 mb-6">
        <div class="flex-1">
          <UInput
            v-model="search"
            icon="i-heroicons-magnifying-glass"
            placeholder="Cari bahan..."
            size="lg"
          />
        </div>
        <USelectMenu
          v-model="selectedFilter"
          :options="filterOptions"
          placeholder="Filter Status"
          size="lg"
        />
        <USelectMenu
          v-model="selectedBatch"
          :options="batchOptions"
          value-attribute="value"
          placeholder="Filter Batch"
          size="lg"
        />
      </div>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead>
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama Bahan</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stok Saat Ini</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stok Minimum</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Satuan</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Harga/Pack</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Est. Terpakai</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="item in filteredIngredients" :key="item.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ item.name }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm" :class="item.current_stock < item.min_stock ? 'text-red-600 font-bold' : 'text-gray-900'">
                {{ item.current_stock }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ item.min_stock }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ item.use_unit }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">Rp {{ (item.price_per_pack || 0).toLocaleString('id-ID') }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm" :class="item.estimated_usage > item.current_stock ? 'text-red-600 font-bold' : 'text-gray-900'">
                {{ (item.estimated_usage || 0).toFixed(2) }} {{ item.use_unit }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <UBadge :color="getStockStatusColor(item)" variant="subtle" size="xs">
                  {{ getStockStatus(item) }}
                </UBadge>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <UDropdown :items="getIngredientActions(item)">
                  <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal" />
                </UDropdown>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Inventory'
})

const toast = useToast()
const { fetchIngredientsWithEstimation, deleteIngredient } = useIngredients()
const { fetchBatches } = useBatches()

const search = ref('')
const selectedFilter = ref('Semua')
const selectedBatch = ref('')
const filterOptions = ['Semua', 'Stok Aman', 'Stok Kritis', 'Stok Minus']
const loading = ref(false)
const ingredients = ref<any[]>([])
const batches = ref<any[]>([])

const batchOptions = computed(() => {
  return [
    { label: 'Semua Batch', value: '' },
    ...batches.value.map((batch: any) => ({
      label: `${batch.name}`,
      value: batch.id
    }))
  ]
})

const loadIngredients = async () => {
  loading.value = true
  const batchId = selectedBatch.value || undefined
  const { data, error } = await fetchIngredientsWithEstimation(batchId)
  loading.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal memuat data bahan',
      color: 'red'
    })
    return
  }

  ingredients.value = (data as any[]) || []
}

const loadBatches = async () => {
  const { data } = await fetchBatches()
  batches.value = (data as any[]) || []
  if (batches.value.length > 0) {
    selectedBatch.value = batches.value[0].id
  }
}

const filteredIngredients = computed(() => {
  let result = ingredients.value

  if (search.value) {
    result = result.filter((item: any) => 
      item.name?.toLowerCase().includes(search.value.toLowerCase())
    )
  }

  if (selectedFilter.value === 'Stok Kritis') {
    result = result.filter((item: any) => item.current_stock < item.min_stock && item.current_stock >= 0)
  } else if (selectedFilter.value === 'Stok Minus') {
    result = result.filter((item: any) => item.current_stock < 0)
  } else if (selectedFilter.value === 'Stok Aman') {
    result = result.filter((item: any) => item.current_stock >= item.min_stock)
  }

  return result
})

const criticalCount = computed(() => {
  return ingredients.value.filter((item: any) => item.current_stock < item.min_stock).length
})

const safeCount = computed(() => {
  return ingredients.value.filter((item: any) => item.current_stock >= item.min_stock).length
})

const handleDelete = async (id: string) => {
  if (!confirm('Yakin ingin menghapus bahan ini?')) return

  const { error } = await deleteIngredient(id)
  
  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal menghapus bahan',
      color: 'red'
    })
    return
  }

  toast.add({
    title: 'Berhasil',
    description: 'Bahan berhasil dihapus',
    color: 'green'
  })

  loadIngredients()
}

onMounted(() => {
  loadBatches()
  loadIngredients()
})

watch(selectedBatch, () => {
  loadIngredients()
})

const getStockStatus = (item: any) => {
  if (item.current_stock < 0) return 'MINUS'
  if (item.current_stock < item.min_stock) return 'KRITIS'
  return 'AMAN'
}

const getStockStatusColor = (item: any) => {
  if (item.current_stock < 0) return 'red'
  if (item.current_stock < item.min_stock) return 'yellow'
  return 'green'
}

const getIngredientActions = (item: any) => {
  return [
    [{
      label: 'Edit',
      icon: 'i-heroicons-pencil',
      click: () => navigateTo(`/inventory/${item.id}/edit`)
    }, {
      label: 'Adjust Stok',
      icon: 'i-heroicons-arrows-right-left',
      click: () => toast.add({
        title: 'Info',
        description: 'Fitur adjust stok akan segera tersedia',
        color: 'blue'
      })
    }],
    [{
      label: 'Hapus',
      icon: 'i-heroicons-trash',
      click: () => handleDelete(item.id)
    }]
  ]
}
</script>
