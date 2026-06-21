<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Resep Produk</h2>
        <p class="mt-1 text-sm text-gray-500">{{ product?.name }} - {{ product?.code }}</p>
      </div>
      <div class="flex gap-3">
        <UButton to="/recipes" color="gray" variant="outline" icon="i-heroicons-arrow-left">
          Kembali
        </UButton>
        <UButton icon="i-heroicons-plus" @click="addRecipe">
          Tambah Bahan
        </UButton>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        <p class="mt-2 text-sm text-gray-500">Memuat data...</p>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <UCard class="lg:col-span-2">
        <template #header>
          <h3 class="text-lg font-semibold">Daftar Bahan</h3>
        </template>

        <div v-if="!recipes || recipes.length === 0" class="flex flex-col items-center justify-center py-12">
          <UIcon name="i-heroicons-inbox" class="w-12 h-12 text-gray-400 mb-3" />
          <p class="text-sm font-medium text-gray-900 mb-1">Belum ada resep</p>
          <p class="text-xs text-gray-500">Tambahkan bahan untuk produk ini</p>
        </div>

        <div v-else class="space-y-4">
          <div v-for="recipe in recipes" :key="recipe.id" class="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
            <div class="flex-1">
              <h4 class="font-medium text-gray-900">{{ recipe.ingredient?.name }}</h4>
              <p class="text-sm text-gray-500">
                {{ recipe.qty_per_box }} {{ recipe.use_unit }} per box
              </p>
              <p class="text-sm font-medium text-primary-600">
                Rp {{ (recipe.cost_per_box || 0).toLocaleString('id-ID') }} / box
              </p>
            </div>
            <div class="flex gap-2">
              <UButton color="gray" variant="ghost" size="sm" icon="i-heroicons-pencil" @click="editRecipe(recipe)" />
              <UButton color="red" variant="ghost" size="sm" icon="i-heroicons-trash" @click="deleteRecipe(recipe.id)" />
            </div>
          </div>
        </div>
      </UCard>

      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold">Ringkasan HPP</h3>
        </template>

        <div class="space-y-4">
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">Harga Jual</span>
            <span class="font-semibold text-gray-900">Rp {{ (product?.price || 0).toLocaleString('id-ID') }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">Total HPP</span>
            <span class="font-semibold text-gray-900">Rp {{ totalHPP.toLocaleString('id-ID') }}</span>
          </div>
          <div class="flex justify-between text-sm pt-4 border-t">
            <span class="text-gray-600">Margin</span>
            <span :class="margin >= 30 ? 'text-green-600' : 'text-yellow-600'" class="font-semibold">
              {{ margin.toFixed(2) }}%
            </span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-gray-600">Laba/Box</span>
            <span class="font-semibold text-primary-600">Rp {{ profit.toLocaleString('id-ID') }}</span>
          </div>
        </div>
      </UCard>
    </div>

    <UModal v-model="isAddModalOpen">
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold">Tambah Bahan</h3>
        </template>

        <form @submit.prevent="handleAddRecipe" class="space-y-4">
          <UFormGroup label="Bahan" required>
            <USelectMenu
              v-model="recipeForm.ingredient_id"
              :options="ingredientOptions"
              value-attribute="value"
              placeholder="Pilih bahan"
            />
          </UFormGroup>

          <UFormGroup label="Jumlah per Box" required>
            <UInput v-model.number="recipeForm.qty_per_box" type="number" step="0.01" placeholder="1.0" />
          </UFormGroup>

          <UFormGroup label="Satuan" required>
            <UInput v-model="recipeForm.use_unit" placeholder="gram, pcs, butir, dll" />
          </UFormGroup>

          <div class="flex justify-end gap-3 pt-4">
            <UButton type="button" color="gray" variant="outline" @click="isAddModalOpen = false">
              Batal
            </UButton>
            <UButton type="submit" :loading="submitting">
              Simpan
            </UButton>
          </div>
        </form>
      </UCard>
    </UModal>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const toast = useToast()
const { fetchProductById } = useProducts()
const { fetchRecipesByProductId, deleteRecipe: deleteRecipeAPI, createRecipe } = useRecipes()
const { fetchIngredients } = useIngredients()

useHead({
  title: 'Resep Produk'
})

const loading = ref(false)
const product = ref<any>(null)
const recipes = ref<any[]>([])
const ingredients = ref<any[]>([])
const isAddModalOpen = ref(false)
const submitting = ref(false)

const recipeForm = ref({
  ingredient_id: '' as string,
  qty_per_box: 0,
  use_unit: ''
})

const ingredientOptions = computed(() => {
  return ingredients.value.map((ing: any) => ({
    label: ing.name,
    value: ing.id
  }))
})

const totalHPP = computed(() => {
  return recipes.value.reduce((sum, recipe) => sum + (recipe.cost_per_box || 0), 0)
})

const profit = computed(() => {
  return (product.value?.price || 0) - totalHPP.value
})

const margin = computed(() => {
  if (!product.value?.price) return 0
  return (profit.value / product.value.price) * 100
})

const loadData = async () => {
  loading.value = true
  
  const [productRes, recipesRes, ingredientsRes] = await Promise.all([
    fetchProductById(route.params.id as string),
    fetchRecipesByProductId(route.params.id as string),
    fetchIngredients()
  ])
  
  loading.value = false

  if (productRes.error) {
    toast.add({
      title: 'Error',
      description: 'Gagal memuat data produk',
      color: 'red'
    })
    navigateTo('/recipes')
    return
  }

  product.value = productRes.data
  recipes.value = (recipesRes.data as any[]) || []
  ingredients.value = (ingredientsRes.data as any[]) || []
}

const addRecipe = () => {
  isAddModalOpen.value = true
}

const handleAddRecipe = async () => {
  if (!recipeForm.value.ingredient_id || !recipeForm.value.qty_per_box || !recipeForm.value.use_unit) {
    toast.add({
      title: 'Error',
      description: 'Mohon lengkapi semua field',
      color: 'red'
    })
    return
  }

  submitting.value = true
  
  const recipeData = {
    product_id: route.params.id,
    ingredient_id: recipeForm.value.ingredient_id,
    qty_per_box: Number(recipeForm.value.qty_per_box),
    use_unit: recipeForm.value.use_unit
  }

  const { error } = await createRecipe(recipeData)
  submitting.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal menambahkan bahan',
      color: 'red'
    })
    return
  }

  toast.add({
    title: 'Berhasil',
    description: 'Bahan berhasil ditambahkan',
    color: 'green'
  })

  isAddModalOpen.value = false
  recipeForm.value = {
    ingredient_id: '',
    qty_per_box: 0,
    use_unit: ''
  }
  loadData()
}

const editRecipe = (recipe: any) => {
  toast.add({
    title: 'Info',
    description: 'Fitur edit resep akan segera tersedia',
    color: 'blue'
  })
}

const deleteRecipe = async (id: string) => {
  if (!confirm('Yakin ingin menghapus bahan ini?')) return

  const { error } = await deleteRecipeAPI(id)
  
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

  loadData()
}

onMounted(() => {
  loadData()
})
</script>
