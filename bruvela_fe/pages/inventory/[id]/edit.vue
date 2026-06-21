<template>
  <div class="space-y-6">
    <div class="flex items-center space-x-4">
      <UButton to="/inventory" color="gray" variant="ghost" icon="i-heroicons-arrow-left" size="lg" />
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Edit Bahan</h2>
        <p class="mt-1 text-sm text-gray-500">{{ ingredient?.name || 'Memuat...' }}</p>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
    </div>

    <UCard v-else-if="ingredient">
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <UFormGroup label="Nama Bahan" required>
            <UInput v-model="form.name" />
          </UFormGroup>

          <UFormGroup label="Satuan Beli (pack)" required>
            <UInput v-model="form.pack_unit" />
          </UFormGroup>

          <UFormGroup label="Isi per Pack" required>
            <UInput v-model.number="form.qty_per_pack" type="number" step="0.001" />
          </UFormGroup>

          <UFormGroup label="Satuan Pakai" required>
            <UInput v-model="form.use_unit" />
          </UFormGroup>

          <UFormGroup label="Harga per Pack (Rp)" required>
            <UInput v-model.number="form.price_per_pack" type="number" min="0" />
          </UFormGroup>

          <UFormGroup label="Stok Minimum">
            <UInput v-model.number="form.min_stock" type="number" step="0.001" />
          </UFormGroup>

          <UFormGroup label="Stok Saat Ini" help="Untuk menambah stok via pembelian, gunakan menu Pembelian">
            <UInput v-model.number="form.current_stock" type="number" step="0.001" />
          </UFormGroup>
        </div>

        <div class="flex justify-end gap-3 pt-4 border-t">
          <UButton type="button" color="gray" variant="outline" @click="navigateTo('/inventory')">
            Batal
          </UButton>
          <UButton type="submit" :loading="submitting" :disabled="!isValid">
            Simpan Perubahan
          </UButton>
        </div>
      </form>
    </UCard>

    <div v-else class="text-center py-12 text-gray-500">
      Bahan tidak ditemukan
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({ title: 'Edit Bahan' })

const route = useRoute()
const toast = useToast()
const { fetchIngredientById, updateIngredient } = useIngredients()

const loading = ref(true)
const submitting = ref(false)
const ingredient = ref<any>(null)

const form = ref({
  name: '',
  pack_unit: '',
  qty_per_pack: 0,
  use_unit: '',
  price_per_pack: 0,
  min_stock: 0,
  current_stock: 0
})

const isValid = computed(() => {
  return (
    form.value.name &&
    form.value.pack_unit &&
    form.value.qty_per_pack > 0 &&
    form.value.use_unit &&
    form.value.price_per_pack > 0
  )
})

const loadData = async () => {
  loading.value = true
  const id = route.params.id as string
  const { data, error } = await fetchIngredientById(id)
  loading.value = false

  if (error || !data) {
    toast.add({ title: 'Error', description: 'Bahan tidak ditemukan', color: 'red' })
    navigateTo('/inventory')
    return
  }

  ingredient.value = data
  const d = data as any
  form.value = {
    name: d.name,
    pack_unit: d.pack_unit,
    qty_per_pack: Number(d.qty_per_pack),
    use_unit: d.use_unit,
    price_per_pack: Number(d.price_per_pack),
    min_stock: Number(d.min_stock) || 0,
    current_stock: Number(d.current_stock) || 0
  }
}

const handleSubmit = async () => {
  if (!isValid.value) {
    toast.add({ title: 'Error', description: 'Lengkapi semua field wajib', color: 'red' })
    return
  }
  submitting.value = true
  const id = route.params.id as string
  const { error } = await updateIngredient(id, { ...form.value })
  submitting.value = false

  if (error) {
    toast.add({ title: 'Gagal', description: 'Gagal memperbarui bahan', color: 'red' })
    return
  }
  toast.add({ title: 'Berhasil', description: 'Bahan berhasil diperbarui', color: 'green' })
  navigateTo('/inventory')
}

onMounted(() => {
  loadData()
})
</script>
