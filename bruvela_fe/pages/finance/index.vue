<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Keuangan</h2>
        <p class="mt-1 text-sm text-gray-500">
          Laporan keuangan per batch
          <span v-if="activeBatch" class="ml-1 text-gray-400">
            — Batch #{{ activeBatch.batch_number }}
          </span>
        </p>
      </div>
      <UButton icon="i-heroicons-plus" size="lg" @click="openCreateModal" :disabled="!activeBatch">
        Tambah Jurnal
      </UButton>
    </div>

    <!-- Empty state if no active batch -->
    <div v-if="!loading && !activeBatch" class="text-center py-16 bg-white rounded-lg border border-dashed border-gray-300">
      <UIcon name="i-heroicons-banknotes" class="w-12 h-12 text-gray-300 mx-auto mb-3" />
      <p class="text-sm font-medium text-gray-900 mb-1">Belum ada batch aktif</p>
      <p class="text-xs text-gray-500 mb-4">Buat batch terlebih dahulu untuk mulai mencatat jurnal</p>
      <UButton to="/batches" icon="i-heroicons-plus">Buat Batch</UButton>
    </div>

    <template v-else>
      <!-- Summary Cards -->
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <UCard>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Total Pemasukan</p>
              <p class="mt-2 text-2xl font-bold text-green-600">{{ formatRupiah(summary.total_income) }}</p>
              <p class="mt-1 text-xs text-gray-500">Modal + Pendapatan</p>
            </div>
            <div class="p-3 bg-green-100 rounded-lg">
              <UIcon name="i-heroicons-arrow-trending-up" class="w-6 h-6 text-green-600" />
            </div>
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Total Pengeluaran</p>
              <p class="mt-2 text-2xl font-bold text-red-600">{{ formatRupiah(summary.total_expense) }}</p>
              <p class="mt-1 text-xs text-gray-500">Bahan + Operasional</p>
            </div>
            <div class="p-3 bg-red-100 rounded-lg">
              <UIcon name="i-heroicons-arrow-trending-down" class="w-6 h-6 text-red-600" />
            </div>
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Saldo</p>
              <p class="mt-2 text-2xl font-bold" :class="(summary.balance || 0) >= 0 ? 'text-gray-900' : 'text-red-600'">
                {{ formatRupiah(summary.balance) }}
              </p>
              <p class="mt-1 text-xs text-gray-500">Pemasukan − Pengeluaran</p>
            </div>
            <div class="p-3 bg-blue-100 rounded-lg">
              <UIcon name="i-heroicons-scale" class="w-6 h-6 text-blue-600" />
            </div>
          </div>
        </UCard>

        <UCard>
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Jumlah Transaksi</p>
              <p class="mt-2 text-2xl font-bold text-gray-900">{{ journalEntries.length }}</p>
              <p class="mt-1 text-xs text-gray-500">entri jurnal</p>
            </div>
            <div class="p-3 bg-purple-100 rounded-lg">
              <UIcon name="i-heroicons-document-text" class="w-6 h-6 text-purple-600" />
            </div>
          </div>
        </UCard>
      </div>

      <!-- Journal List -->
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900">Jurnal Transaksi</h3>
            <UButton color="gray" variant="ghost" size="xs" icon="i-heroicons-arrow-path" @click="loadData" :loading="loading" />
          </div>
        </template>

        <div v-if="loading && journalEntries.length === 0" class="flex justify-center py-12">
          <div class="text-center">
            <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
            <p class="mt-2 text-sm text-gray-500">Memuat jurnal...</p>
          </div>
        </div>

        <div v-else-if="journalEntries.length === 0" class="text-center py-12">
          <UIcon name="i-heroicons-document-plus" class="w-12 h-12 text-gray-300 mx-auto mb-2" />
          <p class="text-sm text-gray-500">Belum ada jurnal. Klik "Tambah Jurnal" untuk mulai.</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead>
              <tr>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tanggal</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Keterangan</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipe</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Partner</th>
                <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Jumlah</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="entry in journalEntries" :key="entry.id" class="hover:bg-gray-50">
                <td class="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{{ formatDate(entry.entry_date) }}</td>
                <td class="px-4 py-3 text-sm text-gray-900">{{ entry.description }}</td>
                <td class="px-4 py-3 whitespace-nowrap">
                  <UBadge :color="getTypeColor(entry.type)" variant="subtle" size="xs">
                    {{ getTypeLabel(entry.type) }}
                  </UBadge>
                </td>
                <td class="px-4 py-3 whitespace-nowrap text-sm text-gray-500">{{ entry.partner || '—' }}</td>
                <td class="px-4 py-3 whitespace-nowrap text-sm font-semibold text-right" :class="isIncome(entry.type) ? 'text-green-600' : 'text-red-600'">
                  {{ isIncome(entry.type) ? '+' : '−' }} Rp {{ Math.abs(entry.amount).toLocaleString('id-ID') }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </UCard>
    </template>

    <!-- Modal Tambah Jurnal -->
    <UModal v-model="showCreateModal">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">Tambah Jurnal Keuangan</h3>
            <UButton color="gray" variant="ghost" icon="i-heroicons-x-mark" @click="showCreateModal = false" />
          </div>
        </template>

        <form @submit.prevent="handleCreate" class="space-y-4">
          <UFormGroup label="Tanggal" required>
            <UInput v-model="form.entry_date" type="date" required />
          </UFormGroup>

          <UFormGroup label="Tipe Transaksi" required>
            <USelectMenu
              v-model="form.type"
              :options="typeOptions"
              value-attribute="value"
              option-attribute="label"
            />
          </UFormGroup>

          <UFormGroup label="Keterangan" required>
            <UInput v-model="form.description" placeholder="cth: Pembelian bahan Batch 3" />
          </UFormGroup>

          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="Jumlah (Rp)" required>
              <UInput v-model.number="form.amount" type="number" min="1" placeholder="0" />
            </UFormGroup>
            <UFormGroup label="Partner (opsional)">
              <UInput v-model="form.partner" placeholder="cth: Aul / Dhavinna" />
            </UFormGroup>
          </div>

          <div class="flex justify-end gap-3 pt-4 border-t">
            <UButton type="button" color="gray" variant="outline" @click="showCreateModal = false">
              Batal
            </UButton>
            <UButton type="submit" :loading="submitting" :disabled="!isFormValid">
              Simpan
            </UButton>
          </div>
        </form>
      </UCard>
    </UModal>
  </div>
</template>

<script setup lang="ts">
useHead({ title: 'Keuangan' })

const toast = useToast()
const { fetchActiveBatch } = useBatches()
const { getJournal, createJournalEntry, getFinanceSummary } = useFinance()

const loading = ref(false)
const submitting = ref(false)
const activeBatch = ref<any>(null)
const summary = ref<any>({ total_income: 0, total_expense: 0, balance: 0 })
const journalEntries = ref<any[]>([])
const showCreateModal = ref(false)

const typeOptions = [
  { label: 'Pemasukan (Pendapatan)', value: 'income' },
  { label: 'Pengeluaran (Bahan)', value: 'expense' },
  { label: 'Modal Masuk', value: 'modal' },
  { label: 'Transfer', value: 'transfer' }
]

const form = ref({
  entry_date: new Date().toISOString().split('T')[0],
  type: 'income',
  description: '',
  amount: 0,
  partner: ''
})

const isFormValid = computed(() => {
  return form.value.description && form.value.amount > 0 && form.value.entry_date
})

const loadData = async () => {
  loading.value = true
  try {
    const batchRes = await fetchActiveBatch()
    if (batchRes.error || !batchRes.data) {
      activeBatch.value = null
      return
    }
    activeBatch.value = batchRes.data
    const batchID = (batchRes.data as any).id

    const [journalRes, summaryRes] = await Promise.all([
      getJournal(batchID),
      getFinanceSummary(batchID)
    ])

    if (journalRes.error) {
      toast.add({ title: 'Error', description: 'Gagal memuat jurnal', color: 'red' })
    } else {
      journalEntries.value = (journalRes.data as any[]) || []
    }

    if (summaryRes.error) {
      console.warn('Summary error:', summaryRes.error)
      summary.value = { total_income: 0, total_expense: 0, balance: 0 }
    } else {
      summary.value = (summaryRes.data as any) || { total_income: 0, total_expense: 0, balance: 0 }
    }
  } catch (err) {
    console.error('Finance load error:', err)
    toast.add({ title: 'Error', description: 'Terjadi kesalahan', color: 'red' })
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  if (!activeBatch.value) {
    toast.add({ title: 'Info', description: 'Belum ada batch aktif', color: 'yellow' })
    return
  }
  // Reset form
  form.value = {
    entry_date: new Date().toISOString().split('T')[0],
    type: 'income',
    description: '',
    amount: 0,
    partner: ''
  }
  showCreateModal.value = true
}

const handleCreate = async () => {
  if (!isFormValid.value) {
    toast.add({ title: 'Error', description: 'Lengkapi semua field wajib', color: 'red' })
    return
  }
  submitting.value = true
  try {
    const payload = {
      batch_id: activeBatch.value.id,
      entry_date: form.value.entry_date,
      type: form.value.type,
      description: form.value.description,
      amount: form.value.amount,
      partner: form.value.partner || null
    }
    const { error } = await createJournalEntry(payload)
    submitting.value = false

    if (error) {
      toast.add({ title: 'Gagal', description: 'Gagal menyimpan jurnal', color: 'red' })
      return
    }

    toast.add({ title: 'Berhasil', description: 'Jurnal berhasil ditambahkan', color: 'green' })
    showCreateModal.value = false
    loadData()
  } catch (err) {
    submitting.value = false
    console.error(err)
    toast.add({ title: 'Error', description: 'Terjadi kesalahan', color: 'red' })
  }
}

// === Helpers ===
const formatRupiah = (value: number) => {
  if (!value && value !== 0) return 'Rp 0'
  return `Rp ${Math.abs(value).toLocaleString('id-ID')}`
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleDateString('id-ID', {
      day: '2-digit', month: 'short', year: 'numeric'
    })
  } catch {
    return dateStr
  }
}

const isIncome = (type: string) => ['income', 'modal'].includes(type)

const getTypeLabel = (type: string) => {
  const map: Record<string, string> = {
    income: 'Pemasukan', expense: 'Pengeluaran',
    modal: 'Modal', transfer: 'Transfer'
  }
  return map[type] || type
}

const getTypeColor = (type: string) => {
  const map: Record<string, string> = {
    income: 'green', expense: 'red',
    modal: 'blue', transfer: 'yellow'
  }
  return map[type] || 'gray'
}

onMounted(() => {
  loadData()
})
</script>
