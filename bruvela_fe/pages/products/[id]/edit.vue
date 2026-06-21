<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Edit Produk</h2>
        <p class="mt-1 text-sm text-gray-500">Ubah informasi produk</p>
      </div>
      <UButton to="/recipes" color="gray" variant="outline" icon="i-heroicons-arrow-left">
        Kembali
      </UButton>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        <p class="mt-2 text-sm text-gray-500">Memuat data...</p>
      </div>
    </div>

    <UCard v-else>
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <UFormGroup label="Kode Produk" required>
            <UInput v-model="form.code" placeholder="BRV-CLASSIC" />
          </UFormGroup>

          <UFormGroup label="Nama Produk" required>
            <UInput v-model="form.name" placeholder="Bruv Classic" />
          </UFormGroup>

          <UFormGroup label="Harga Jual (Rp)" required>
            <UInput v-model.number="form.price" type="number" placeholder="20000" />
          </UFormGroup>

          <UFormGroup label="Pcs Per Box">
            <UInput v-model.number="form.pcs_per_box" type="number" placeholder="1" />
          </UFormGroup>

          <UFormGroup label="Status">
            <USelectMenu v-model="form.status" :options="statusOptions" />
          </UFormGroup>
        </div>

        <div class="flex justify-end gap-3 pt-4 border-t">
          <UButton type="button" color="gray" variant="outline" @click="navigateTo('/recipes')">
            Batal
          </UButton>
          <UButton type="submit" :loading="submitting">
            Simpan Perubahan
          </UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const toast = useToast()
const { fetchProductById, updateProduct } = useProducts()

useHead({
  title: 'Edit Produk'
})

const loading = ref(false)
const submitting = ref(false)
const form = ref({
  code: '',
  name: '',
  price: 0,
  pcs_per_box: 1,
  status: 'active'
})

const statusOptions = ['active', 'inactive']

const loadProduct = async () => {
  loading.value = true
  const { data, error } = await fetchProductById(route.params.id as string)
  loading.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal memuat data produk',
      color: 'red'
    })
    navigateTo('/recipes')
    return
  }

  const product = data as any
  form.value = {
    code: product.code,
    name: product.name,
    price: product.price,
    pcs_per_box: product.pcs_per_box,
    status: product.status
  }
}

const handleSubmit = async () => {
  if (!form.value.code || !form.value.name || !form.value.price) {
    toast.add({
      title: 'Error',
      description: 'Mohon lengkapi semua field yang wajib diisi',
      color: 'red'
    })
    return
  }

  submitting.value = true
  const { error } = await updateProduct(route.params.id as string, form.value)
  submitting.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal mengupdate produk',
      color: 'red'
    })
    return
  }

  toast.add({
    title: 'Berhasil',
    description: 'Produk berhasil diupdate',
    color: 'green'
  })

  navigateTo('/recipes')
}

onMounted(() => {
  loadProduct()
})
</script>
