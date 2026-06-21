<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Resep & Menu</h2>
        <p class="mt-1 text-sm text-gray-500">Kelola menu produk dan HPP</p>
      </div>
      <div class="flex gap-3">
        <UButton to="/recipes/calculator" color="gray" variant="outline" icon="i-heroicons-calculator" size="lg">
          Kalkulator Produksi
        </UButton>
        <UButton icon="i-heroicons-plus" size="lg" @click="addMenu">
          Tambah Menu
        </UButton>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        <p class="mt-2 text-sm text-gray-500">Memuat data...</p>
      </div>
    </div>

    <div v-else-if="!products || products.length === 0" class="flex flex-col items-center justify-center py-12">
      <UIcon name="i-heroicons-inbox" class="w-12 h-12 text-gray-400 mb-3" />
      <p class="text-sm font-medium text-gray-900 mb-1">Tidak ada data produk</p>
      <p class="text-xs text-gray-500">Belum ada produk dengan resep yang dibuat</p>
    </div>

    <div v-else class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
      <UCard v-for="product in products" :key="product.code">
        <div class="space-y-4">
          <div class="flex items-start justify-between">
            <div>
              <h3 class="text-lg font-semibold text-gray-900">{{ product.name }}</h3>
              <p class="text-sm text-gray-500">{{ product.code }}</p>
            </div>
            <UBadge :color="product.status === 'active' ? 'green' : 'gray'" variant="subtle" size="xs">
              {{ product.status === 'active' ? 'Aktif' : 'Nonaktif' }}
            </UBadge>
          </div>

          <div class="space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Harga Jual</span>
              <span class="font-semibold text-gray-900">Rp {{ product.price.toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">HPP</span>
              <span class="font-semibold text-gray-900">Rp {{ product.hpp.toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Margin</span>
              <span :class="product.margin >= 30 ? 'text-green-600' : 'text-yellow-600'" class="font-semibold">
                {{ product.margin?.toFixed(2) }}%
              </span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Laba/Box</span>
              <span class="font-semibold text-primary-600">Rp {{ product.profit.toLocaleString('id-ID') }}</span>
            </div>
          </div>

          <div class="pt-4 border-t border-gray-200">
            <div class="flex gap-2">
              <UButton color="gray" variant="outline" size="sm" class="flex-1" @click="viewRecipe(product)">
                <UIcon name="i-heroicons-eye" class="w-4 h-4 mr-1" />
                Lihat Resep
              </UButton>
              <UButton color="gray" variant="outline" size="sm" class="flex-2" @click="editProduct(product)">
                <UIcon name="i-heroicons-pencil" class="w-4 h-4 mr-1" />
                Edit
              </UButton>
            </div>
          </div>
        </div>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Resep & Menu'
})

const toast = useToast()
const { fetchRecipes } = useRecipes()

const products = ref<any[]>([])
const loading = ref(false)

const loadRecipes = async () => {
  loading.value = true
  const { data, error } = await fetchRecipes()
  loading.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal memuat data resep',
      color: 'red'
    })
    return
  }

  products.value = (data as any[]) || []
}

onMounted(() => {
  loadRecipes()
})

const addMenu = () => {
  navigateTo('/products/create')
}

const viewRecipe = (product: any) => {
  navigateTo(`/products/${product.id}/recipe`)
}

const editProduct = (product: any) => {
  navigateTo(`/products/${product.id}/edit`)
}
</script>
