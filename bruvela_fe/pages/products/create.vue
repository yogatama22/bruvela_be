<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Tambah Menu Baru</h2>
        <p class="mt-1 text-sm text-gray-500">Buat produk menu baru</p>
      </div>
      <UButton to="/recipes" color="gray" variant="outline" icon="i-heroicons-arrow-left">
        Kembali
      </UButton>
    </div>

    <UCard>
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
          <UButton type="submit" :loading="loading">
            Simpan Produk
          </UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Tambah Menu Baru'
})

const toast = useToast()
const { createProduct } = useProducts()

const loading = ref(false)
const form = ref({
  code: '',
  name: '',
  price: 0,
  pcs_per_box: 1,
  status: 'active'
})

const statusOptions = ['active', 'inactive']

const handleSubmit = async () => {
  if (!form.value.code || !form.value.name || !form.value.price) {
    toast.add({
      title: 'Error',
      description: 'Mohon lengkapi semua field yang wajib diisi',
      color: 'red'
    })
    return
  }

  loading.value = true
  const { error } = await createProduct(form.value)
  loading.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal membuat produk',
      color: 'red'
    })
    return
  }

  toast.add({
    title: 'Berhasil',
    description: 'Produk berhasil dibuat',
    color: 'green'
  })

  navigateTo('/recipes')
}
</script>
