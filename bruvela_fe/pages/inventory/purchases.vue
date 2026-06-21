<template>
  <div class="space-y-6">
    <div class="flex items-center space-x-4">
      <UButton to="/inventory" color="gray" variant="ghost" icon="i-heroicons-arrow-left" size="lg" />
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Pembelian Bahan</h2>
        <p class="mt-1 text-sm text-gray-500">
          Catat pembelian bahan baku
          <span v-if="activeBatch" class="ml-1 text-gray-400">— Batch #{{ activeBatch.batch_number }}</span>
        </p>
      </div>
      <div class="ml-auto">
        <UButton icon="i-heroicons-plus" size="lg" @click="openCreate" :disabled="!activeBatch">
          Catat Pembelian
        </UButton>
      </div>
    </div>

    <div v-if="!loading && !activeBatch" class="text-center py-16 bg-white rounded-lg border border-dashed border-gray-300">
      <UIcon name="i-heroicons-shopping-cart" class="w-12 h-12 text-gray-300 mx-auto mb-3" />
      <p class="text-sm font-medium text-gray-900 mb-1">Belum ada batch aktif</p>
      <p class="text-xs text-gray-500 mb-4">Buat batch terlebih dahulu untuk mulai mencatat pembelian</p>
      <UButton to="/batches" icon="i-heroicons-plus">Buat Batch</UButton>
    </div>

    <template v-else>
      <!-- Summary -->
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-3">
        <UCard>
          <p class="text-sm font-medium text-gray-600">Total Pembelian</p>
          <p class="mt-2 text-2xl font-bold text-gray-900">{{ purchases.length }}</p>
          <p class="mt-1 text-xs text-gray-500">transaksi</p>
        </UCard>
        <UCard>
          <p class="text-sm font-medium text-gray-600">Total Nilai</p>
          <p class="mt-2 text-2xl font-bold text-red-600">{{ formatRupiah(totalAmount) }}</p>
          <p class="mt-1 text-xs text-gray-500">pengeluaran bahan</p>
        </UCard>
        <UCard>
          <p class="text-sm font-medium text-gray-600">Supplier Unik</p>
          <p class="mt-2 text-2xl font-bold text-gray-900">{{ uniqueSuppliers }}</p>
          <p class="mt-1 text-xs text-gray-500">berbeda</p>
        </UCard>
      </div>

      <!-- Table -->
      <UCard>
        <div class="flex flex-col sm:flex-row gap-4 mb-4">
          <div class="flex-1">
            <UInput v-model="search" icon="i-heroicons-magnifying-glass" placeholder="Cari pembelian..." size="lg" />
          </div>
        </div>

        <div v-if="loading" class="flex justify-center py-12">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        </div>

        <div v-else-if="filteredPurchases.length === 0" class="text-center py-12">
          <UIcon name="i-heroicons-inbox" class="w-12 h-12 text-gray-300 mx-auto mb-2" />
          <p class="text-sm text-gray-500">Belum ada pembelian untuk batch ini</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead>
              <tr>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tanggal</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Bahan</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Supplier</th>
                <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Qty (pack)</th>
                <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Harga/Pack</th>
                <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Aksi</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="p in filteredPurchases" :key="p.id" class="hover:bg-gray-50">
                <td class="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{{ formatDate(p.purchase_date) }}</td>
                <td class="px-4 py-3">
                  <div class="text-sm font-medium text-gray-900">{{ p.Ingredient?.name || p.ingredient?.name || '—' }}</div>
                </td>
                <td class="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{{ p.supplier || '—' }}</td>
                <td class="px-4 py-3 whitespace-nowrap text-sm text-right">{{ Number(p.qty_pack).toFixed(2) }}</td>
                <td class="px-4 py-3 whitespace-nowrap text-sm text-right">Rp {{ Number(p.price_per_pack).toLocaleString('id-ID') }}</td>
                <td class="px-4 py-3 whitespace-nowrap text-sm text-right font-semibold text-red-600">
                  Rp {{ Number(p.total_price).toLocaleString('id-ID') }}
                </td>
                <td class="px-4 py-3 whitespace-nowrap text-sm">
                  <UButton color="red" variant="ghost" size="xs" icon="i-heroicons-trash" @click="handleDelete(p)" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </UCard>
    </template>

    <!-- Modal Create -->
    <UModal v-model="showCreate">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">Catat Pembelian Bahan</h3>
            <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark" @click="showCreate = false" />
          </div>
        </template>

        <form @submit.prevent="handleCreate" class="space-y-4">
          <UFormGroup label="Tanggal Pembelian" required>
            <UInput v-model="form.purchase_date" type="date" required />
          </UFormGroup>

          <UFormGroup label="Bahan" required>
            <USelectMenu
              v-model="form.ingredient_id"
              :options="ingredientOptions"
              value-attribute="value"
              option-attribute="label"
              placeholder="Pilih bahan"
            />
          </UFormGroup>

          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="Qty (pack)" required>
              <UInput v-model.number="form.qty_pack" type="number" step="0.001" min="0.001" placeholder="cth: 2" />
            </UFormGroup>
            <UFormGroup label="Harga per Pack (Rp)" required>
              <UInput v-model.number="form.price_per_pack" type="number" min="0" placeholder="cth: 15000" />
            </UFormGroup>
          </div>

          <div v-if="formTotal > 0" class="p-3 bg-blue-50 rounded-lg">
            <p class="text-sm text-gray-700">
              <span class="font-semibold">Total:</span> Rp {{ formTotal.toLocaleString('id-ID') }}
            </p>
            <p class="text-xs text-gray-500 mt-1">Akan otomatis menambah stok & jurnal expense</p>
          </div>

          <UFormGroup label="Supplier (opsional)">
            <UInput v-model="form.supplier" placeholder="cth: Toko Sumber Rezeki" />
          </UFormGroup>

          <UFormGroup label="Catatan (opsional)">
            <UTextarea v-model="form.note" placeholder="Catatan tambahan..." :rows="2" />
          </UFormGroup>

          <div class="flex justify-end gap-3 pt-4 border-t">
            <UButton type="button" color="gray" variant="outline" @click="showCreate = false">Batal</UButton>
            <UButton type="submit" :loading="submitting" :disabled="!isFormValid">Simpan</UButton>
          </div>
        </form>
      </UCard>
    </UModal>
  </div>
</template>

<script setup lang="ts">
useHead({ title: 'Pembelian Bahan' })

const toast = useToast()
const { fetchActiveBatch } = useBatches()
const { fetchIngredients } = useIngredients()
const { fetchPurchases, createPurchase } = useIngredients()

const loading = ref(false)
const submitting = ref(false)
const activeBatch = ref<any>(null)
const ingredients = ref<any[]>([])
const purchases = ref<any[]>([])
const showCreate = ref(false)
const search = ref('')

const form = ref({
  purchase_date: new Date().toISOString().split('T')[0],
  ingredient_id: '' as string,
  qty_pack: 0,
  price_per_pack: 0,
  supplier: '',
  note: ''
})

const ingredientOptions = computed(() => {
  return ingredients.value.map(i => ({ label: i.name, value: i.id }))
})

const isFormValid = computed(() => {
  return form.value.ingredient_id && form.value.qty_pack > 0 && form.value.price_per_pack > 0
})

const formTotal = computed(() => {
  return Math.round((form.value.qty_pack || 0) * (form.value.price_per_pack || 0))
})

const filteredPurchases = computed(() => {
  let r = purchases.value
  if (search.value) {
    const q = search.value.toLowerCase()
    r = r.filter((p: any) =>
      (p.Ingredient?.name || p.ingredient?.name || '').toLowerCase().includes(q) ||
      (p.supplier || '').toLowerCase().includes(q)
    )
  }
  return r
})

const totalAmount = computed(() => {
  return purchases.value.reduce((sum, p) => sum + Number(p.total_price || 0), 0)
})

const uniqueSuppliers = computed(() => {
  const set = new Set(purchases.value.map(p => p.supplier).filter(Boolean))
  return set.size
})

const loadAll = async () => {
  loading.value = true
  try {
    const batchRes = await fetchActiveBatch()
    if (batchRes.error || !batchRes.data) {
      activeBatch.value = null
      return
    }
    activeBatch.value = batchRes.data
    const batchID = (batchRes.data as any).id

    const [ingRes, purRes] = await Promise.all([
      fetchIngredients(),
      fetchPurchases()
    ])

    if (!ingRes.error) ingredients.value = (ingRes.data as any[]) || []

    if (!purRes.error) {
      let data = (purRes.data as any[]) || []
      // Filter by current batch client-side if backend didn't filter
      data = data.filter((p: any) => p.batch_id === batchID || p.Batch?.id === batchID)
      purchases.value = data
    }
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  form.value = {
    purchase_date: new Date().toISOString().split('T')[0],
    ingredient_id: '',
    qty_pack: 0,
    price_per_pack: 0,
    supplier: '',
    note: ''
  }
  showCreate.value = true
}

const handleCreate = async () => {
  if (!isFormValid.value || !activeBatch.value) {
    toast.add({ title: 'Error', description: 'Lengkapi semua field', color: 'red' })
    return
  }

  submitting.value = true
  const payload = {
    batch_id: activeBatch.value.id,
    ingredient_id: form.value.ingredient_id,
    purchase_date: form.value.purchase_date,
    qty_pack: Number(form.value.qty_pack),
    price_per_pack: Number(form.value.price_per_pack),
    supplier: form.value.supplier,
    note: form.value.note
  }

  const { error } = await createPurchase(payload)
  submitting.value = false

  if (error) {
    toast.add({ title: 'Gagal', description: 'Gagal mencatat pembelian', color: 'red' })
    return
  }
  toast.add({
    title: 'Berhasil',
    description: `Pembelian dicatat. Stok bahan bertambah & jurnal expense otomatis terbuat.`,
    color: 'green'
  })
  showCreate.value = false
  await loadAll()
}

const handleDelete = async (p: any) => {
  if (!confirm(`Hapus pembelian ${p.Ingredient?.name || ''}? Stok & jurnal TIDAK akan di-rollback.`)) return
  // We need delete API; for now show warning
  toast.add({
    title: 'Info',
    description: 'Hapus pembelian belum tersedia di backend. Edit manual via DB jika urgent.',
    color: 'yellow'
  })
}

const formatRupiah = (val: number) => `Rp ${(val || 0).toLocaleString('id-ID')}`
const formatDate = (d: string) => {
  if (!d) return '-'
  try {
    return new Date(d).toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
  } catch { return d }
}

onMounted(() => loadAll())
</script>
