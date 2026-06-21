<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Production Calculator</h2>
        <p class="mt-1 text-sm text-gray-500">Hitung kebutuhan bahan baku untuk produksi</p>
      </div>
      <UButton to="/recipes" color="gray" variant="ghost" icon="i-heroicons-arrow-left" size="lg">
        Kembali
      </UButton>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-1 space-y-6">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900">Pilih Produk</h3>
          </template>
          <div class="space-y-4">
            <div v-for="(item, index) in productionItems" :key="index" class="flex items-end gap-3">
              <div class="flex-1">
                <label class="text-xs font-medium text-gray-500">Produk</label>
                <USelectMenu
                  v-model="item.productId"
                  :options="productOptions"
                  value-attribute="value"
                  option-attribute="label"
                  placeholder="Pilih produk"
                  size="md"
                />
              </div>
              <div class="w-24">
                <label class="text-xs font-medium text-gray-500">Qty Box</label>
                <UInput v-model.number="item.qtyBox" type="number" placeholder="0" size="md" />
              </div>
              <UButton
                @click="removeItem(index)"
                color="red"
                variant="ghost"
                icon="i-heroicons-trash"
                size="md"
              />
            </div>
            <UButton
              @click="addItem"
              color="gray"
              variant="outline"
              icon="i-heroicons-plus"
              block
            >
              Tambah Produk
            </UButton>
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-600">Total HPP Estimasi</p>
              <p class="text-2xl font-bold text-primary-600">Rp {{ (calcResult?.total_hpp || 0).toLocaleString('id-ID', { maximumFractionDigits: 0 }) }}</p>
            </div>
            <UButton
              @click="calculate"
              :loading="calculating"
              :disabled="!hasValidItems"
              size="lg"
            >
              Hitung
            </UButton>
          </div>
        </UCard>
      </div>

      <div class="lg:col-span-2 space-y-6">
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-900">Kebutuhan Bahan Baku</h3>
              <div v-if="calcResult" class="flex gap-2">
                <UBadge color="green" variant="subtle" size="xs">Cukup: {{ statusCount.cukup }}</UBadge>
                <UBadge color="yellow" variant="subtle" size="xs">Pas-pasan: {{ statusCount['pas-pasan'] }}</UBadge>
                <UBadge color="red" variant="subtle" size="xs">Kurang: {{ statusCount.kurang }}</UBadge>
              </div>
            </div>
          </template>

          <div v-if="calculating" class="flex flex-col items-center justify-center py-12">
            <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mb-3"></div>
            <p class="text-sm text-gray-500">Menghitung...</p>
          </div>

          <div v-else-if="!calcResult" class="flex flex-col items-center justify-center py-12">
            <UIcon name="i-heroicons-calculator" class="w-12 h-12 text-gray-400 mb-3" />
            <p class="text-sm font-medium text-gray-900 mb-1">Belum ada hasil</p>
            <p class="text-xs text-gray-500">Pilih produk dan klik "Hitung" untuk melihat kebutuhan bahan</p>
          </div>

          <div v-else-if="calcResult.ingredients.length === 0" class="flex flex-col items-center justify-center py-12">
            <UIcon name="i-heroicons-inbox" class="w-12 h-12 text-gray-400 mb-3" />
            <p class="text-sm text-gray-500">Tidak ada bahan yang dibutuhkan</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead>
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Bahan</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Butuh</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Stok Saat Ini</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <tr v-for="ing in calcResult.ingredients" :key="ing.id" class="hover:bg-gray-50">
                  <td class="px-4 py-3 text-sm font-medium text-gray-900">{{ ing.name }}</td>
                  <td class="px-4 py-3 text-sm text-right text-gray-900">{{ ing.needed.toFixed(3) }}</td>
                  <td class="px-4 py-3 text-sm text-right text-gray-500">{{ ing.current_stock.toFixed(3) }}</td>
                  <td class="px-4 py-3">
                    <div class="flex items-center gap-2">
                      <span class="w-2.5 h-2.5 rounded-full" :class="getStatusDot(ing.status)"></span>
                      <span class="text-xs font-medium" :class="getStatusText(ing.status)">{{ getStatusLabel(ing.status) }}</span>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </UCard>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Production Calculator'
})

const toast = useToast()
const { fetchProducts } = useProducts()
const { calculateProduction } = useRecipes()

const products = ref<any[]>([])
const calculating = ref(false)
const calcResult = ref<any>(null)

const productionItems = ref<{ productId: string; qtyBox: number }[]>([
  { productId: '', qtyBox: 1 }
])

const productOptions = computed(() => {
  return products.value.map((p: any) => ({ label: p.name, value: p.id }))
})

const hasValidItems = computed(() => {
  return productionItems.value.some(item => item.productId && item.qtyBox > 0)
})

const statusCount = computed(() => {
  if (!calcResult.value?.ingredients) return { cukup: 0, 'pas-pasan': 0, kurang: 0 }
  return calcResult.value.ingredients.reduce((acc: any, ing: any) => {
    acc[ing.status] = (acc[ing.status] || 0) + 1
    return acc
  }, { cukup: 0, 'pas-pasan': 0, kurang: 0 })
})

const addItem = () => {
  productionItems.value.push({ productId: '', qtyBox: 1 })
}

const removeItem = (index: number) => {
  if (productionItems.value.length > 1) {
    productionItems.value.splice(index, 1)
  }
}

const calculate = async () => {
  const validItems = productionItems.value
    .filter(item => item.productId && item.qtyBox > 0)
    .map(item => ({ product_id: item.productId, qty_box: item.qtyBox }))

  if (validItems.length === 0) {
    toast.add({ title: 'Error', description: 'Pilih minimal 1 produk', color: 'red' })
    return
  }

  calculating.value = true
  const { data, error } = await calculateProduction(validItems)
  calculating.value = false

  if (error) {
    toast.add({ title: 'Error', description: 'Gagal menghitung', color: 'red' })
    return
  }

  calcResult.value = data
}

const getStatusDot = (status: string) => {
  if (status === 'cukup') return 'bg-green-500'
  if (status === 'pas-pasan') return 'bg-yellow-500'
  return 'bg-red-500'
}

const getStatusText = (status: string) => {
  if (status === 'cukup') return 'text-green-700'
  if (status === 'pas-pasan') return 'text-yellow-700'
  return 'text-red-700'
}

const getStatusLabel = (status: string) => {
  if (status === 'cukup') return 'Cukup'
  if (status === 'pas-pasan') return 'Pas-pasan'
  return 'Kurang'
}

const loadProducts = async () => {
  const { data } = await fetchProducts()
  products.value = data || []
}

onMounted(() => {
  loadProducts()
})
</script>
