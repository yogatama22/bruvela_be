<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Manajemen Batch</h2>
        <p class="mt-1 text-sm text-gray-500">Kelola periode batch produksi</p>
      </div>
      <UButton @click="showCreateModal = true" icon="i-heroicons-plus" size="lg">
        Batch Baru
      </UButton>
    </div>

    <!-- Active Batch Card -->
    <UCard v-if="activeBatch" class="bg-primary-50 border-2 border-primary-500">
      <div class="flex items-center justify-between">
        <div>
          <div class="flex items-center gap-2">
            <UBadge color="green" variant="solid">AKTIF</UBadge>
            <h3 class="text-lg font-semibold text-gray-900">Batch #{{ activeBatch.batch_number }}</h3>
          </div>
          <p class="text-sm text-gray-600 mt-1">{{ activeBatch.name || 'Tanpa Nama' }}</p>
          <p class="text-xs text-gray-500 mt-1">
            {{ formatDate(activeBatch.start_date) }} - {{ activeBatch.end_date ? formatDate(activeBatch.end_date) : 'Berlangsung' }}
          </p>
        </div>
        <div class="text-right">
          <p class="text-sm text-gray-600">Total Modal</p>
          <p class="text-2xl font-bold text-gray-900">Rp {{ (activeBatch.total_modal || 0).toLocaleString('id-ID') }}</p>
        </div>
      </div>
    </UCard>

    <!-- Batches List -->
    <UCard>
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead>
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Batch #</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Nama</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Periode</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Modal</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Revenue</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Profit</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" class="hover:bg-gray-50">
              <td colspan="8" class="px-6 py-12 text-center">
                <div class="flex flex-col items-center justify-center">
                  <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mb-3"></div>
                  <p class="text-sm text-gray-500">Memuat data...</p>
                </div>
              </td>
            </tr>
            <tr v-else-if="!loading && batches.length === 0" class="hover:bg-gray-50">
              <td colspan="8" class="px-6 py-12 text-center">
                <div class="flex flex-col items-center justify-center">
                  <UIcon name="i-heroicons-inbox" class="w-12 h-12 text-gray-400 mb-3" />
                  <p class="text-sm font-medium text-gray-900 mb-1">Tidak ada data batch</p>
                  <p class="text-xs text-gray-500">Buat batch baru untuk memulai</p>
                </div>
              </td>
            </tr>
            <tr v-else v-for="batch in batches" :key="batch.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">#{{ batch.batch_number }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ batch.name || '-' }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(batch.start_date) }} - {{ batch.end_date ? formatDate(batch.end_date) : '-' }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <UBadge :color="batch.status === 'open' ? 'green' : 'gray'" variant="subtle" size="xs">
                  {{ batch.status === 'open' ? 'Aktif' : 'Tutup' }}
                </UBadge>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">Rp {{ (batch.total_modal || 0).toLocaleString('id-ID') }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">Rp {{ (batch.total_revenue || 0).toLocaleString('id-ID') }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium" :class="batch.gross_profit >= 0 ? 'text-green-600' : 'text-red-600'">
                Rp {{ (batch.gross_profit || 0).toLocaleString('id-ID') }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <UDropdown :items="getBatchActions(batch)">
                  <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal" />
                </UDropdown>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </UCard>

    <!-- Create/Edit Modal -->
    <UModal v-model="showCreateModal">
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold">{{ editingBatch ? 'Edit Batch' : 'Batch Baru' }}</h3>
        </template>

        <div class="space-y-4">
          <UFormGroup label="Batch Number" required>
            <UInput v-model.number="form.batch_number" type="number" placeholder="1" />
          </UFormGroup>
          <UFormGroup label="Nama Batch">
            <UInput v-model="form.name" placeholder="Batch Januari 2024" />
          </UFormGroup>
          <div class="grid grid-cols-2 gap-4">
            <UFormGroup label="Tanggal Mulai" required>
              <UInput v-model="form.start_date" type="date" />
            </UFormGroup>
            <UFormGroup label="Tanggal Selesai">
              <UInput v-model="form.end_date" type="date" />
            </UFormGroup>
          </div>
          <UFormGroup label="Status">
            <USelectMenu
              v-model="form.status"
              :options="[{ label: 'Aktif', value: 'open' }, { label: 'Tutup', value: 'closed' }]"
              option-attribute="label"
              value-attribute="value"
            />
          </UFormGroup>
        </div>

        <template #footer>
          <div class="flex justify-end gap-3">
            <UButton color="gray" variant="ghost" @click="closeModal">Batal</UButton>
            <UButton @click="saveBatch" :loading="saving">Simpan</UButton>
          </div>
        </template>
      </UCard>
    </UModal>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Manajemen Batch'
})

const toast = useToast()
const { fetchBatches, fetchActiveBatch, createBatch, updateBatch, activateBatch, closeBatch, deleteBatch } = useBatches()

const loading = ref(false)
const saving = ref(false)
const batches = ref<any[]>([])
const activeBatch = ref<any>(null)
const showCreateModal = ref(false)
const editingBatch = ref<any>(null)

const form = ref({
  batch_number: 0,
  name: '',
  start_date: new Date().toISOString().split('T')[0],
  end_date: '',
  status: 'open'
})

const loadBatches = async () => {
  loading.value = true
  const { data, error } = await fetchBatches()
  loading.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal memuat data batch',
      color: 'red'
    })
    return
  }

  batches.value = data || []
}

const loadActiveBatch = async () => {
  const { data } = await fetchActiveBatch()
  activeBatch.value = data
}

const saveBatch = async () => {
  if (!form.value.batch_number) {
    toast.add({
      title: 'Error',
      description: 'Batch number harus diisi',
      color: 'red'
    })
    return
  }

  saving.value = true

  const batchData = {
    batch_number: Number(form.value.batch_number),
    name: form.value.name,
    start_date: new Date(form.value.start_date).toISOString(),
    end_date: form.value.end_date ? new Date(form.value.end_date).toISOString() : null,
    status: form.value.status
  }

  const { error } = editingBatch.value
    ? await updateBatch(editingBatch.value.id, batchData)
    : await createBatch(batchData)

  saving.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: error,
      color: 'red'
    })
    return
  }

  toast.add({
    title: 'Berhasil',
    description: `Batch berhasil ${editingBatch.value ? 'diupdate' : 'dibuat'}`,
    color: 'green'
  })

  closeModal()
  loadBatches()
  loadActiveBatch()
}

const handleActivate = async (batch: any) => {
  if (!confirm(`Aktifkan Batch #${batch.batch_number}? Batch lain akan ditutup.`)) return

  const { error } = await activateBatch(batch.id)

  if (error) {
    toast.add({
      title: 'Error',
      description: error,
      color: 'red'
    })
    return
  }

  toast.add({
    title: 'Berhasil',
    description: 'Batch berhasil diaktifkan',
    color: 'green'
  })

  loadBatches()
  loadActiveBatch()
}

const handleClose = async (batch: any) => {
  if (!confirm(`Tutup Batch #${batch.batch_number}?`)) return

  const { error } = await closeBatch(batch.id)

  if (error) {
    toast.add({
      title: 'Error',
      description: error,
      color: 'red'
    })
    return
  }

  toast.add({
    title: 'Berhasil',
    description: 'Batch berhasil ditutup',
    color: 'green'
  })

  loadBatches()
  loadActiveBatch()
}

const handleEdit = (batch: any) => {
  editingBatch.value = batch
  form.value = {
    batch_number: batch.batch_number,
    name: batch.name || '',
    start_date: new Date(batch.start_date).toISOString().split('T')[0],
    end_date: batch.end_date ? new Date(batch.end_date).toISOString().split('T')[0] : '',
    status: batch.status
  }
  showCreateModal.value = true
}

const handleDelete = async (batch: any) => {
  if (!confirm(`Hapus Batch #${batch.batch_number}?`)) return

  const { error } = await deleteBatch(batch.id)

  if (error) {
    toast.add({
      title: 'Error',
      description: error,
      color: 'red'
    })
    return
  }

  toast.add({
    title: 'Berhasil',
    description: 'Batch berhasil dihapus',
    color: 'green'
  })

  loadBatches()
  loadActiveBatch()
}

const closeModal = () => {
  showCreateModal.value = false
  editingBatch.value = null
  form.value = {
    batch_number: 0,
    name: '',
    start_date: new Date().toISOString().split('T')[0],
    end_date: '',
    status: 'open'
  }
}

const getBatchActions = (batch: any) => {
  const actions = []

  if (batch.status !== 'open') {
    actions.push([{
      label: 'Aktifkan',
      icon: 'i-heroicons-check-circle',
      click: () => handleActivate(batch)
    }])
  }

  if (batch.status === 'open') {
    actions.push([{
      label: 'Tutup',
      icon: 'i-heroicons-x-circle',
      click: () => handleClose(batch)
    }])
  }

  actions.push([
    {
      label: 'Detail',
      icon: 'i-heroicons-eye',
      click: () => navigateTo(`/batches/${batch.id}`)
    },
    {
      label: 'Edit',
      icon: 'i-heroicons-pencil-square',
      click: () => handleEdit(batch)
    },
    {
      label: 'Hapus',
      icon: 'i-heroicons-trash',
      click: () => handleDelete(batch)
    }
  ])

  return actions
}

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
}

onMounted(() => {
  loadBatches()
  loadActiveBatch()
})
</script>
