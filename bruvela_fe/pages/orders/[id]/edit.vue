<template>
  <div class="space-y-6">
    <div class="flex items-center space-x-4">
      <UButton
        :to="`/orders/${route.params.id}`"
        color="gray"
        variant="ghost"
        icon="i-heroicons-arrow-left"
        size="lg"
      />
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Edit Order</h2>
        <p class="mt-1 text-sm text-gray-500">Order #{{ route.params.id?.toString().substring(0, 8) }}</p>
      </div>
    </div>

    <div v-if="loadingData" class="flex justify-center py-12">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        <p class="mt-2 text-sm text-gray-500">Memuat data...</p>
      </div>
    </div>

    <div v-else-if="form" class="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <div class="lg:col-span-2 space-y-6">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900">Informasi Customer</h3>
          </template>
          <div class="space-y-4">
            <UFormGroup label="Nama Customer" required>
              <UInput v-model="form.customerName" placeholder="Masukkan nama customer" size="lg" />
            </UFormGroup>
            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="Tanggal Order" required>
                <UInput v-model="form.orderDate" type="date" size="lg" />
              </UFormGroup>
              <UFormGroup label="Channel" required>
                <USelectMenu
                  v-model="form.channel"
                  :options="['whatsapp', 'instagram', 'offline', 'titip_teman']"
                  placeholder="Pilih channel"
                  size="lg"
                />
              </UFormGroup>
            </div>
          </div>
        </UCard>

        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-900">Produk</h3>
              <UButton @click="addProduct" icon="i-heroicons-plus" size="sm">
                Tambah Produk
              </UButton>
            </div>
          </template>
          <div class="space-y-4">
            <div v-for="(item, index) in form.items" :key="index" class="p-4 border border-gray-200 rounded-lg">
              <div class="flex items-start justify-between mb-4">
                <span class="text-sm font-medium text-gray-700">Produk {{ index + 1 }}</span>
                <UButton
                  v-if="form.items.length > 1"
                  @click="removeProduct(index)"
                  color="red"
                  variant="ghost"
                  icon="i-heroicons-trash"
                  size="xs"
                />
              </div>
              <div class="grid grid-cols-2 gap-4">
                <UFormGroup label="Varian">
                  <USelectMenu
                    v-model="item.productId"
                    :options="productOptions"
                    option-attribute="label"
                    value-attribute="value"
                    placeholder="Pilih varian"
                    size="lg"
                  />
                </UFormGroup>
                <UFormGroup label="Jumlah Box">
                  <UInput v-model.number="item.qty" type="number" min="1" placeholder="0" size="lg" />
                </UFormGroup>
              </div>
            </div>
          </div>
        </UCard>

        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900">Pengiriman</h3>
          </template>
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <UFormGroup label="Tipe Pengiriman">
                <USelectMenu
                  v-model="form.shippingType"
                  :options="shippingTypes.map((st: any) => ({ label: st.shipping_name, value: st.shipping_code }))"
                  option-attribute="label"
                  value-attribute="value"
                  placeholder="Pilih tipe"
                  size="lg"
                />
              </UFormGroup>
              <UFormGroup label="Ongkir">
                <UInput v-model.number="form.shippingCost" type="number" placeholder="0" size="lg">
                  <template #leading>
                    <span class="text-gray-500">Rp</span>
                  </template>
                </UInput>
              </UFormGroup>
            </div>
            <UFormGroup label="Alamat Tujuan">
              <UTextarea v-model="form.shippingDest" placeholder="Masukkan alamat lengkap" rows="3" />
            </UFormGroup>
          </div>
        </UCard>
      </div>

      <div class="space-y-6">
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900">Ringkasan</h3>
          </template>
          <div class="space-y-3">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Subtotal Produk</span>
              <span class="font-medium text-gray-900">Rp {{ (subtotal || 0).toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Ongkir</span>
              <span class="font-medium text-gray-900">Rp {{ (form.shippingCost || 0).toLocaleString('id-ID') }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Diskon</span>
              <div class="flex items-center gap-2">
                <span class="text-gray-500 text-xs">Rp</span>
                <UInput v-model.number="form.discount" type="number" placeholder="0" size="sm" class="w-32" />
              </div>
            </div>
            <div class="pt-3 border-t border-gray-200">
              <div class="flex justify-between">
                <span class="text-base font-semibold text-gray-900">Total</span>
                <span class="text-lg font-bold text-primary-600">Rp {{ (totalBill || 0).toLocaleString('id-ID') }}</span>
              </div>
            </div>
          </div>
        </UCard>

        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-gray-900">Status</h3>
          </template>
          <div class="space-y-4">
            <UFormGroup label="Status Order">
              <USelectMenu
                v-model="form.prodStatus"
                :options="orderStatuses"
                placeholder="Pilih status"
                size="lg"
              />
            </UFormGroup>
            <UFormGroup label="Status Pembayaran">
              <USelectMenu
                v-model="form.paymentStatus"
                :options="paymentStatuses.map((ps: any) => ({ label: ps.status_name, value: ps.status_code }))"
                option-attribute="label"
                value-attribute="value"
                placeholder="Pilih status"
                size="lg"
              />
            </UFormGroup>
            <UFormGroup label="Catatan">
              <UTextarea v-model="form.note" placeholder="Catatan tambahan (opsional)" rows="3" />
            </UFormGroup>
          </div>
        </UCard>

        <div class="flex flex-col gap-3">
          <UButton @click="saveOrder" size="lg" block :loading="loading" :disabled="loading">
            Simpan Perubahan
          </UButton>
          <UButton :to="`/orders/${route.params.id}`" color="gray" variant="outline" size="lg" block>
            Batal
          </UButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Edit Order'
})

const route = useRoute()
const toast = useToast()
const { updateOrder, fetchOrderById } = useOrders()
const { fetchProducts } = useProducts()
const { fetchShippingTypes } = useShippingTypes()
const { fetchPaymentStatuses } = usePaymentStatuses()

const loading = ref(false)
const loadingData = ref(true)
const products = ref<any[]>([])
const productMap = ref<Record<string, any>>({})
const shippingTypes = ref<any[]>([])
const paymentStatuses = ref<any[]>([])
const originalOrder = ref<any>(null)

const orderStatuses = [
  { label: 'Baru', value: 'baru' },
  { label: 'Proses', value: 'proses' },
  { label: 'Siap Kirim', value: 'siap_kirim' },
  { label: 'Selesai', value: 'selesai' },
  { label: 'Batal', value: 'batal' }
]

const form = ref<any>({
  customerName: '',
  orderDate: '',
  channel: 'whatsapp',
  items: [],
  shippingType: '',
  shippingCost: 0,
  shippingDest: '',
  discount: 0,
  prodStatus: 'baru',
  paymentStatus: 'belum_bayar',
  note: ''
})

const loadProducts = async () => {
  const { data, error } = await fetchProducts()

  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal memuat data produk',
      color: 'red'
    })
    return
  }

  products.value = (data as any[]) || []
  productMap.value = ((data as any[]) || []).reduce((acc: any, p: any) => {
    acc[p.id] = p
    return acc
  }, {})
}

const productOptions = computed(() => {
  return products.value.map(p => ({
    label: `${p.name} - Rp ${p.price?.toLocaleString('id-ID')}`,
    value: p.id
  }))
})

const loadOrder = async () => {
  loadingData.value = true
  const { data, error } = await fetchOrderById(route.params.id as string)
  
  if (error) {
    toast.add({
      title: 'Error',
      description: 'Gagal memuat detail order',
      color: 'red'
    })
    loadingData.value = false
    return
  }

  const orderData = data as any
  originalOrder.value = orderData

  form.value = {
    customerName: orderData.customer_name || '',
    orderDate: orderData.order_date ? new Date(orderData.order_date).toISOString().split('T')[0] : '',
    channel: orderData.channel || 'whatsapp',
    items: (orderData.items || []).map((item: any) => ({
      productId: item.product_id,
      qty: item.qty_box
    })),
    shippingType: orderData.shipping_type || '',
    shippingCost: Number(orderData.shipping_cost) || 0,
    shippingDest: orderData.shipping_dest || '',
    discount: Number(orderData.discount) || 0,
    prodStatus: orderData.prod_status || 'baru',
    paymentStatus: orderData.pay_status || 'belum_bayar',
    note: orderData.note || ''
  }

  if (form.value.items.length === 0) {
    form.value.items = [{ productId: '', qty: 1 }]
  }

  loadingData.value = false
}

const addProduct = () => {
  form.value.items.push({ productId: '', qty: 1 })
}

const removeProduct = (index: number) => {
  form.value.items.splice(index, 1)
}

const subtotal = computed(() => {
  if (!form.value) return 0
  return form.value.items.reduce((sum: number, item: any) => {
    const product = productMap.value[item.productId]
    const price = product?.price || 0
    return sum + (price * Number(item.qty || 0))
  }, 0)
})

const totalBill = computed(() => {
  if (!form.value) return 0
  return subtotal.value + Number(form.value.shippingCost || 0) - Number(form.value.discount || 0)
})

const saveOrder = async () => {
  if (!form.value.customerName) {
    toast.add({
      title: 'Error',
      description: 'Nama customer harus diisi',
      color: 'red'
    })
    return
  }

  if (form.value.items.length === 0 || !form.value.items[0].productId) {
    toast.add({
      title: 'Error',
      description: 'Minimal 1 produk harus dipilih',
      color: 'red'
    })
    return
  }

  loading.value = true

  const orderData = {
    batch_id: originalOrder.value.batch_id,
    customer_name: form.value.customerName,
    order_date: new Date(form.value.orderDate).toISOString(),
    channel: form.value.channel,
    shipping_type: form.value.shippingType,
    shipping_dest: form.value.shippingDest,
    shipping_cost: Number(form.value.shippingCost),
    discount: Number(form.value.discount),
    prod_status: form.value.prodStatus,
    pay_status: form.value.paymentStatus,
    note: form.value.note,
    items: form.value.items.map((item: any) => {
      const product = productMap.value[item.productId]
      return {
        product_id: item.productId,
        product_code: product?.code || '',
        product_name: product?.name || '',
        qty_box: Number(item.qty),
        price_per_box: product?.price || 0
      }
    })
  }

  const { data, error } = await updateOrder(route.params.id as string, orderData)
  loading.value = false

  if (error) {
    toast.add({
      title: 'Error',
      description: error.message || 'Gagal menyimpan perubahan',
      color: 'red'
    })
    return
  }

  toast.add({
    title: 'Berhasil',
    description: 'Order berhasil diperbarui',
    color: 'green'
  })

  navigateTo(`/orders/${route.params.id}`)
}

onMounted(async () => {
  await loadProducts()
  await loadShippingTypes()
  await loadPaymentStatuses()
  await loadOrder()
})

const loadPaymentStatuses = async () => {
  const { data, error } = await fetchPaymentStatuses()
  if (error) {
    console.error('Failed to load payment statuses:', error)
    return
  }
  paymentStatuses.value = data || []
}

const loadShippingTypes = async () => {
  const { data, error } = await fetchShippingTypes()
  if (error) {
    console.error('Failed to load shipping types:', error)
    return
  }
  shippingTypes.value = data || []
}
</script>
