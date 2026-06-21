<template>
  <div class="space-y-6">
    <div class="flex items-center space-x-4">
      <UButton to="/inventory" color="gray" variant="ghost" icon="i-heroicons-arrow-left" size="lg" />
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Tambah Bahan Baru</h2>
        <p class="mt-1 text-sm text-gray-500">Daftarkan bahan baku baru ke master inventory</p>
      </div>
    </div>

    <UCard>
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
          <UFormGroup label="Nama Bahan" required>
            <UInput v-model="form.name" placeholder="cth: Tepung Terigu" />
          </UFormGroup>

          <UFormGroup label="Satuan Beli (pack)" required>
            <UInput v-model="form.pack_unit" placeholder="cth: kg, pcs, lembar" />
          </UFormGroup>

          <UFormGroup label="Isi per Pack" required help="Berapa satuan pakai dalam 1 pack">
            <UInput v-model.number="form.qty_per_pack" type="number" step="0.001" placeholder="cth: 1000 (untuk 1 kg dalam gram)" />
          </UFormGroup>

          <UFormGroup label="Satuan Pakai" required>
            <UInput v-model="form.use_unit" placeholder="cth: gram, butir, pcs" />
          </UFormGroup>

          <UFormGroup label="Harga per Pack (Rp)" required>
            <UInput v-model.number="form.price_per_pack" type="number" min="0" placeholder="cth: 15000" />
          </UFormGroup>

          <UFormGroup label="Stok Minimum" help="Alert restock jika stok di bawah nilai ini">
            <UInput v-model.number="form.min_stock" type="number" step="0.001" placeholder="cth: 500" />
          </UFormGroup>

          <UFormGroup label="Stok Awal Saat Ini" help="Opsional, default 0">
            <UInput v-model.number="form.current_stock" type="number" step="0.001" placeholder="0" />
          </UFormGroup>
        </div>

        <div v-if="computedPricePerUse > 0" class="p-4 bg-blue-50 rounded-lg">
          <p class="text-sm text-gray-700">
            <span class="font-semibold">Harga per {{ form.use_unit }}:</span>
            Rp {{ computedPricePerUse.toFixed(4).replace(/\.?0+$/, '') }}
          </p>
          <p class="text-xs text-gray-500 mt-1">Otomatis dihitung: harga_pack / qty_per_pack</p>
        </div>

        <div class="flex justify-end gap-3 pt-4 border-t">
          <UButton type="button" color="gray" variant="outline" @click="navigateTo('/inventory')">
            Batal
          </UButton>
          <UButton type="submit" :loading="loading" :disabled="!isValid">
            Simpan Bahan
          </UButton>
        </div>
      </form>
    </UCard>
  </div>
</template>

<script setup lang="ts">
useHead({ title: 'Tambah Bahan' })

const toast = useToast()
const { createIngredient } = useIngredients()

const loading = ref(false)
const form = ref({
  name: '',
  pack_unit: '',
  qty_per_pack: 0,
  use_unit: '',
  price_per_pack: 0,
  min_stock: 0,
  current_stock: 0
})

const computedPricePerUse = computed(() => {
  if (form.value.qty_per_pack > 0) {
    return form.value.price_per_pack / form.value.qty_per_pack
  }
  return 0
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

const handleSubmit = async () => {
  if (!isValid.value) {
    toast.add({ title: 'Error', description: 'Lengkapi semua field wajib', color: 'red' })
    return
  }

  loading.value = true
  const payload = {
    name: form.value.name,
    pack_unit: form.value.pack_unit,
    qty_per_pack: Number(form.value.qty_per_pack),
    use_unit: form.value.use_unit,
    price_per_pack: Number(form.value.price_per_pack),
    min_stock: Number(form.value.min_stock) || 0,
    current_stock: Number(form.value.current_stock) || 0
  }

  const { error } = await createIngredient(payload)
  loading.value = false

  if (error) {
    toast.add({
      title: 'Gagal',
      description: error.error || 'Gagal menyimpan bahan',
      color: 'red'
    })
    return
  }

  toast.add({ title: 'Berhasil', description: 'Bahan berhasil ditambahkan', color: 'green' })
  navigateTo('/inventory')
}
</script>
