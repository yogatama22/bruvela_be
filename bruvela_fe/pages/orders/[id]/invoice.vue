<template>
  <div class="min-h-screen bg-gray-100 py-8 px-4 print:bg-transparent print:p-0">
    <div class="max-w-4xl mx-auto print:max-w-none print:mx-0">
      
      <!-- Action Buttons (Hidden on Print) -->
      <div class="mb-4 flex justify-end gap-3 print:hidden">
        <UButton to="/orders" color="gray" variant="ghost" icon="i-heroicons-arrow-left">
          Kembali
        </UButton>
        <UButton @click="handlePrint" color="primary" icon="i-heroicons-printer">
          Cetak / Save PDF
        </UButton>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-24">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      </div>

      <!-- Invoice Document -->
      <div v-else-if="order" 
           class="invoice-page shadow-lg print:shadow-none mx-auto relative bg-white text-[#4a3e35] font-sans" 
           style="max-width: 210mm; min-height: 297mm;">
        
        <!-- LUNAS Watermark -->
        <div v-if="order.pay_status === 'lunas'" class="lunas-watermark">
          LUNAS
        </div>

        <!-- Invoice Padding -->
        <div class="p-[12mm] print:p-0">
          
          <!-- Header -->
          <div class="flex justify-between items-start mb-10">
            <div>
              <h1 class="font-['Georgia','Times_New_Roman',serif] text-[46px] text-[#4a3e35] tracking-[2px] m-0 leading-[1.1] lowercase">
                {{ company?.brand_name || 'bruvela' }}
              </h1>
              <div class="font-sans text-[13px] tracking-[5px] uppercase text-[#6d5d52] mb-2 mt-1">
                — BAKEHOUSE —
              </div>
              <div class="font-sans text-[10px] tracking-[1.5px] uppercase text-[#9c8e82]">
                MADE BY TWO, BAKED WITH LOVE
              </div>
            </div>
            <div class="text-right">
              <h2 class="text-[36px] font-bold text-[#4a3e35] tracking-[2px] m-0 mb-3">INVOICE</h2>
              <div class="grid grid-cols-[auto_1fr] gap-x-2 text-[12px] text-[#6d5d52] justify-end inline-grid text-right">
                <span class="font-bold">No. Invoice:</span> 
                <span>INV-{{ order.id?.substring(0, 8).toUpperCase() }}</span>
                
                <span class="font-bold">Tanggal:</span> 
                <span>{{ formatDate(order.order_date) }}</span>
                
                <!-- <span class="font-bold">Jatuh Tempo:</span> 
                <span>{{ formatDate(order.due_date) || '-' }}</span> -->
              </div>
            </div>
          </div>

          <!-- Customer Info -->
          <div class="grid grid-cols-2 gap-10 mb-10">
            <!-- Company Info -->
            <div class="pr-5">
              <h3 class="text-[13px] font-bold uppercase tracking-[1px] border-b border-[#dcd1c8] pb-1 mb-2 text-[#6d5d52]">
                Informasi Usaha
              </h3>
              <div class="text-[13px] text-[#4a3e35] leading-[1.6]">
                <strong>{{ company?.company_name || 'Bruvela Bakehouse' }}</strong><br>
                <template v-if="company?.address">
                  {{ company.address }}<br>
                </template>
                <template v-else>
                  JL. Masjid Al Barkah Syakira Residence 2 Blok E6<br>
                  Depok, Jawa Barat 16519<br>
                </template>
                IG: {{ company?.instagram || '@bruvelabakehouse' }}<br>
                WA: {{ company?.whatsapp || '0812-3456-7890' }}
              </div>
            </div>
            
            <!-- Bill To -->
            <div class="pl-5">
              <h3 class="text-[13px] font-bold uppercase tracking-[1px] border-b border-[#dcd1c8] pb-1 mb-2 text-[#6d5d52]">
                Ditagihkan Kepada
              </h3>
              <div class="text-[13px] text-[#4a3e35] leading-[1.6]">
                <strong>{{ order.customer_name || '-' }}</strong><br>
                <template v-if="order.shipping_dest">
                  {{ order.shipping_dest }}<br>
                </template>
                <template v-if="order.customer_phone">
                  {{ order.customer_phone }}<br>
                </template>
                <template v-if="order.shipping_type">
                  <span class="text-[#888] mt-1 inline-block">Pengiriman: {{ order.shipping_type }}</span>
                </template>
              </div>
            </div>
          </div>

          <!-- Items Table -->
          <table class="w-full mb-10 border-collapse">
            <thead>
              <tr>
                <th class="bg-[#eee7e0] border-b-2 border-[#dcd1c8] text-[#4a3e35] py-3 px-3 text-left text-[12px] uppercase tracking-[1px] w-[5%]">No</th>
                <th class="bg-[#eee7e0] border-b-2 border-[#dcd1c8] text-[#4a3e35] py-3 px-3 text-left text-[12px] uppercase tracking-[1px] w-[45%]">Deskripsi Pesanan</th>
                <th class="bg-[#eee7e0] border-b-2 border-[#dcd1c8] text-[#4a3e35] py-3 px-3 text-center text-[12px] uppercase tracking-[1px] w-[15%]">Jumlah</th>
                <th class="bg-[#eee7e0] border-b-2 border-[#dcd1c8] text-[#4a3e35] py-3 px-3 text-right text-[12px] uppercase tracking-[1px] w-[15%]">Harga Satuan</th>
                <th class="bg-[#eee7e0] border-b-2 border-[#dcd1c8] text-[#4a3e35] py-3 px-3 text-right text-[12px] uppercase tracking-[1px] w-[20%]">Total</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, i) in orderItems" :key="i">
                <td class="py-3 px-3 border-b border-[#e8e1da] text-[13px] align-top">{{ i + 1 }}</td>
                <td class="py-3 px-3 border-b border-[#e8e1da] text-[13px] align-top">
                  <strong>{{ item.product_name }}</strong>
                  <div v-if="item.variant || item.notes" class="text-[#888] text-[11px] mt-1 leading-tight">
                    {{ item.variant || item.notes }}
                  </div>
                </td>
                <td class="py-3 px-3 border-b border-[#e8e1da] text-[13px] align-top text-center">{{ item.qty_box || item.qty || 1 }}</td>
                <td class="py-3 px-3 border-b border-[#e8e1da] text-[13px] align-top text-right">Rp {{ (item.price_per_box || item.price || 0).toLocaleString('id-ID') }}</td>
                <td class="py-3 px-3 border-b border-[#e8e1da] text-[13px] align-top text-right">Rp {{ (item.subtotal || 0).toLocaleString('id-ID') }}</td>
              </tr>
            </tbody>
          </table>

          <!-- Totals & Payment -->
          <div class="grid grid-cols-[40%_60%] gap-5">
            <!-- Payment Method -->
            <div class="pr-5">
              <!-- <h3 class="text-[13px] font-bold uppercase tracking-[1px] border-b border-[#dcd1c8] pb-1 mb-3 text-[#6d5d52]">
                Metode Pembayaran
              </h3>
              <div class="bg-[#f5f0ec] p-4 rounded text-[13px] leading-[1.6]">
                <template v-if="company?.payment_info">
                  <div v-html="company.payment_info.replace(/\n/g, '<br>')"></div>
                </template>
                <template v-else>
                  Silakan transfer pembayaran ke rekening berikut:<br><br>
                  <strong>Bank BCA</strong><br>
                  No. Rekening: 1234 567 890<br>
                  Atas Nama: Bruvela Bakehouse<br><br>
                  <em>Mohon sertakan nomor invoice pada berita transfer.</em>
                </template>
              </div> -->

              <!-- Extra Note -->
              <div v-if="order.note" class="mt-4 text-[12px] text-[#6d5d52] italic">
                <strong>Catatan Pesanan:</strong> {{ order.note }}
              </div>
            </div>

            <!-- Summary Table -->
            <div>
              <table class="w-full text-[13px] border-collapse">
                <colgroup>
                  <col style="width: auto;">
                  <col style="width: auto;">
                  <col style="width: 5%;">
                </colgroup>
                <tbody>
                  <tr>
                    <td class="py-2 px-3">Subtotal</td>
                    <td class="py-2 px-3 text-right whitespace-nowrap">Rp {{ itemsSubtotal.toLocaleString('id-ID') }}</td>
                    <td></td>
                  </tr>
                  <tr v-if="order.discount">
                    <td class="py-2 px-3">Diskon</td>
                    <td class="py-2 px-3 text-right text-red-600 whitespace-nowrap">- Rp {{ (order.discount || 0).toLocaleString('id-ID') }}</td>
                    <td></td>
                  </tr>
                  <tr>
                    <td class="py-2 px-3">Ongkos Kirim <span v-if="order.shipping_courier" class="text-[11px] text-[#888]">({{ order.shipping_courier }})</span></td>
                    <td class="py-2 px-3 text-right whitespace-nowrap">Rp {{ (order.shipping_cost || 0).toLocaleString('id-ID') }}</td>
                    <td></td>
                  </tr>
                  <tr class="bg-[#eee7e0] border-y-2 border-[#dcd1c8] font-bold text-[15px]">
                    <td class="py-3 px-3">TOTAL</td>
                    <td class="py-3 px-3 text-right whitespace-nowrap">Rp {{ (order.total_bill || 0).toLocaleString('id-ID') }}</td>
                    <td></td>
                  </tr>
                </tbody>
              </table>

              <!-- Payment Status Indicator -->
              <div class="mt-3 text-right text-[13px] whitespace-nowrap pr-[4%]">
                Status: 
                <span class="font-bold uppercase" 
                      :class="order.pay_status === 'lunas' ? 'text-green-600' : 'text-red-600'">
                  {{ formatPayStatus(order.pay_status) }}
                </span>
              </div>
            </div>
          </div>

          <!-- Footer -->
          <div class="absolute bottom-0 left-0 w-full text-center border-t border-[#dcd1c8] pt-4 pb-8 print:pb-0 text-[11px] text-[#9c8e82] tracking-[1px] uppercase">
            Terima kasih atas pesanan Anda! | MADE BY TWO, BAKED WITH LOVE
          </div>

        </div>
      </div>

      <!-- Not Found State -->
      <div v-else class="text-center py-24">
        <p class="text-gray-500">Order tidak ditemukan</p>
        <UButton to="/orders" class="mt-4">Kembali ke Daftar Order</UButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: 'Invoice - Bruvela Bakehouse'
})

const route = useRoute()
const { fetchOrderById } = useOrders() // Asumsi composable Anda
const { fetchSettings } = useCompanySettings() // Asumsi composable Anda

const loading = ref(true)
const order = ref<any>(null)
const orderItems = ref<any[]>([])
const company = ref<any>(null)

const loadOrder = async () => {
  loading.value = true
  const { data, error } = await fetchOrderById(route.params.id as string)
  loading.value = false

  if (error) return

  order.value = data
  orderItems.value = data?.items || []
}

const loadCompany = async () => {
  const { data } = await fetchSettings()
  company.value = data
}

const itemsSubtotal = computed(() => {
  return orderItems.value.reduce((sum: number, item: any) => sum + (item.subtotal || 0), 0)
})

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })
}

const formatPayStatus = (status: string) => {
  const map: Record<string, string> = {
    'belum_bayar': 'Belum Bayar',
    'dp': 'DP',
    'lunas': 'Lunas'
  }
  return map[status] || status
}

const handlePrint = () => {
  window.print()
}

onMounted(() => {
  loadOrder()
  loadCompany()
})
</script>

<style scoped>
/* Lunas Watermark Customization */
.lunas-watermark {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%) rotate(-30deg);
  font-size: 140px;
  font-weight: 900;
  color: rgba(34, 197, 94, 0.08); /* Sangat tipis agar tidak mengganggu bacaan */
  pointer-events: none;
  z-index: 0;
  letter-spacing: 15px;
  user-select: none;
}

/* Base Print Page Rules */
@media print {
  @page {
    size: A4;
    margin: 10mm;
  }

  body {
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
  }

  /* Hide layout chrome: sidebar, header, mobile menu */
  :global(aside),
  :global(header),
  :global(.fixed),
  :global([class*="slideover"]) {
    display: none !important;
  }

  /* Reset layout containers for print */
  :global(.flex.h-screen) {
    display: block !important;
    overflow: visible !important;
    height: auto !important;
  }

  :global(.flex.flex-col.flex-1) {
    overflow: visible !important;
    height: auto !important;
  }

  :global(main) {
    overflow: visible !important;
    height: auto !important;
  }

  :global(.px-4.py-6) {
    padding: 0 !important;
  }

  .invoice-page {
    background-color: #ffffff !important;
    min-height: auto !important;
    max-width: 100% !important;
    overflow: visible !important;
    margin: 0 !important;
  }
  
  /* Menghilangkan Absolute Position dari footer untuk cetak jika halaman berlebih */
  .absolute.bottom-0 {
    position: relative !important;
    margin-top: 40px !important;
    padding-bottom: 0 !important;
  }
}
</style>