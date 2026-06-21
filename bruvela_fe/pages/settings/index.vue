<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-2xl font-bold text-gray-900">Pengaturan Perusahaan</h2>
      <p class="mt-1 text-sm text-gray-500">Data perusahaan untuk invoice dan dokumen</p>
    </div>

    <UCard v-if="loading" class="py-12">
      <div class="flex justify-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      </div>
    </UCard>

    <UCard v-else>
      <template #header>
        <h3 class="text-lg font-semibold text-gray-900">Informasi Perusahaan</h3>
      </template>
      <div class="space-y-4">
        <UFormGroup label="Nama Perusahaan" required>
          <UInput v-model="form.company_name" placeholder="Bruvela Bakehouse" size="lg" />
        </UFormGroup>

        <UFormGroup label="Alamat">
          <UTextarea v-model="form.address" placeholder="Jl. Contoh No. 123, Kota" :rows="3" />
        </UFormGroup>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <UFormGroup label="Nomor Telepon">
            <UInput v-model="form.phone_number" placeholder="0812-3456-7890" />
          </UFormGroup>
          <UFormGroup label="Email">
            <UInput v-model="form.email" placeholder="info@bruvela.com" type="email" />
          </UFormGroup>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <UFormGroup label="WhatsApp">
            <UInput v-model="form.whatsapp" placeholder="0812-3456-7890" />
          </UFormGroup>
          <UFormGroup label="Instagram">
            <UInput v-model="form.instagram" placeholder="@bruvela.bakehouse" />
          </UFormGroup>
        </div>

        <UFormGroup label="TikTok">
          <UInput v-model="form.tiktok" placeholder="@bruvela" />
        </UFormGroup>

        <UFormGroup label="Catatan Invoice (opsional)">
          <UTextarea v-model="form.invoice_note" placeholder="Terima kasih telah berbelanja!" :rows="2" />
        </UFormGroup>
      </div>

      <template #footer>
        <div class="flex justify-end">
          <UButton @click="saveSettings" :loading="saving" size="lg">
            Simpan
          </UButton>
        </div>
      </template>
    </UCard>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Pengaturan Perusahaan'
})

const toast = useToast()
const { fetchSettings, updateSettings } = useCompanySettings()

const loading = ref(true)
const saving = ref(false)

const form = ref({
  company_name: '',
  address: '',
  phone_number: '',
  email: '',
  whatsapp: '',
  instagram: '',
  tiktok: '',
  invoice_note: ''
})

const loadSettings = async () => {
  loading.value = true
  const { data, error } = await fetchSettings()
  loading.value = false

  if (error) {
    toast.add({ title: 'Error', description: 'Gagal memuat pengaturan', color: 'red' })
    return
  }

  if (data) {
    form.value = {
      company_name: data.company_name || '',
      address: data.address || '',
      phone_number: data.phone_number || '',
      email: data.email || '',
      whatsapp: data.whatsapp || '',
      instagram: data.instagram || '',
      tiktok: data.tiktok || '',
      invoice_note: data.invoice_note || ''
    }
  }
}

const saveSettings = async () => {
  if (!form.value.company_name) {
    toast.add({ title: 'Error', description: 'Nama perusahaan harus diisi', color: 'red' })
    return
  }

  saving.value = true
  const { error } = await updateSettings(form.value)
  saving.value = false

  if (error) {
    toast.add({ title: 'Error', description: 'Gagal menyimpan pengaturan', color: 'red' })
    return
  }

  toast.add({ title: 'Berhasil', description: 'Pengaturan perusahaan disimpan', color: 'green' })
}

onMounted(() => {
  loadSettings()
})
</script>
