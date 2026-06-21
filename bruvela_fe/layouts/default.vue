<template>
  <div class="min-h-screen bg-gray-50">
    <div class="flex h-screen overflow-hidden">
      <aside 
        class="hidden lg:flex lg:flex-shrink-0"
        :class="sidebarOpen ? 'lg:w-64' : 'lg:w-20'"
      >
        <div class="flex flex-col w-full border-r border-gray-200 bg-white">
          <div class="flex items-center justify-between h-16 px-4 border-b border-gray-200">
            <div v-if="sidebarOpen" class="flex items-center space-x-3">
              <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
                <span class="text-white font-bold text-lg">B</span>
              </div>
              <span class="text-lg font-semibold text-gray-900">Bruvela</span>
            </div>
            <div v-else class="flex items-center justify-center w-full">
              <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
                <span class="text-white font-bold text-lg">B</span>
              </div>
            </div>
          </div>

          <nav class="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
            <NuxtLink
              v-for="item in navigation"
              :key="item.name"
              :to="item.href"
              class="flex items-center px-3 py-2 text-sm font-medium rounded-lg transition-colors"
              :class="[
                isActive(item.href)
                  ? 'bg-primary-50 text-primary-700'
                  : 'text-gray-700 hover:bg-gray-100'
              ]"
            >
              <UIcon :name="item.icon" class="flex-shrink-0 w-5 h-5" />
              <span v-if="sidebarOpen" class="ml-3">{{ item.name }}</span>
              <UBadge
                v-if="sidebarOpen && item.badge"
                :color="item.badgeColor || 'red'"
                variant="solid"
                size="xs"
                class="ml-auto"
              >
                {{ item.badge }}
              </UBadge>
            </NuxtLink>
          </nav>

          <div class="flex-shrink-0 p-4 border-t border-gray-200">
            <button
              @click="sidebarOpen = !sidebarOpen"
              class="flex items-center w-full px-3 py-2 text-sm font-medium text-gray-700 rounded-lg hover:bg-gray-100 transition-colors"
            >
              <UIcon 
                :name="sidebarOpen ? 'i-heroicons-chevron-left' : 'i-heroicons-chevron-right'" 
                class="w-5 h-5"
              />
              <span v-if="sidebarOpen" class="ml-3">Collapse</span>
            </button>
          </div>
        </div>
      </aside>

      <div class="flex flex-col flex-1 overflow-hidden">
        <header class="bg-white border-b border-gray-200">
          <div class="flex items-center justify-between h-16 px-4 sm:px-6 lg:px-8">
            <div class="flex items-center">
              <button
                @click="mobileMenuOpen = !mobileMenuOpen"
                class="lg:hidden p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100"
              >
                <UIcon name="i-heroicons-bars-3" class="w-6 h-6" />
              </button>
              <h1 class="ml-4 lg:ml-0 text-xl font-semibold text-gray-900">
                {{ currentPageTitle }}
              </h1>
            </div>

            <div class="flex items-center space-x-4">
              <UButton
                icon="i-heroicons-bell"
                color="gray"
                variant="ghost"
                size="lg"
              />
              <UDropdown :items="userMenuItems">
                <div class="flex items-center gap-3 cursor-pointer">
                  <UAvatar
                    :src="`https://ui-avatars.com/api/?name=${user?.name || user?.email || 'User'}&background=4f46e5&color=fff`"
                    size="sm"
                  />
                  <div class="hidden sm:block">
                    <p class="text-sm font-medium text-gray-900">{{ user?.name || 'User' }}</p>
                    <p class="text-xs text-gray-500">{{ user?.email }}</p>
                  </div>
                </div>
              </UDropdown>
            </div>
          </div>
        </header>

        <main class="flex-1 overflow-y-auto bg-gray-50">
          <div class="px-4 py-6 sm:px-6 lg:px-8">
            <slot />
          </div>
        </main>
      </div>
    </div>

    <USlideover v-model="mobileMenuOpen" side="left">
      <div class="p-4">
        <div class="flex items-center space-x-3 mb-6">
          <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
            <span class="text-white font-bold text-lg">B</span>
          </div>
          <span class="text-lg font-semibold text-gray-900">Bruvela</span>
        </div>

        <nav class="space-y-1">
          <NuxtLink
            v-for="item in navigation"
            :key="item.name"
            :to="item.href"
            @click="mobileMenuOpen = false"
            class="flex items-center px-3 py-2 text-sm font-medium rounded-lg transition-colors"
            :class="[
              isActive(item.href)
                ? 'bg-primary-50 text-primary-700'
                : 'text-gray-700 hover:bg-gray-100'
            ]"
          >
            <UIcon :name="item.icon" class="flex-shrink-0 w-5 h-5" />
            <span class="ml-3">{{ item.name }}</span>
            <UBadge
              v-if="item.badge"
              :color="item.badgeColor || 'red'"
              variant="solid"
              size="xs"
              class="ml-auto"
            >
              {{ item.badge }}
            </UBadge>
          </NuxtLink>
        </nav>
      </div>
    </USlideover>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const { user, logout } = useAuth()
const sidebarOpen = ref(true)
const mobileMenuOpen = ref(false)

const navigation = computed(() => [
  {
    name: 'Dashboard',
    href: '/',
    icon: 'i-heroicons-home',
  },
  {
    name: 'Pemesanan',
    href: '/orders',
    icon: 'i-heroicons-shopping-bag',
  },
  {
    name: 'Resep & Menu',
    href: '/recipes',
    icon: 'i-heroicons-book-open',
  },
  {
    name: 'Inventory',
    href: '/inventory',
    icon: 'i-heroicons-cube',
    // badge: inventoryStore.criticalStockCount,
    badgeColor: 'red'
  },
  {
    name: 'Stock Movement',
    href: '/inventory/stock-logs',
    icon: 'i-heroicons-arrows-right-left',
  },
  {
    name: 'Keuangan',
    href: '/finance',
    icon: 'i-heroicons-banknotes',
  },
  {
    name: 'Batch',
    href: '/batches',
    icon: 'i-heroicons-calendar-days',
  },
  {
    name: 'Laporan',
    href: '/reports',
    icon: 'i-heroicons-document-chart-bar',
  },
  {
    name: 'Pengaturan',
    href: '/settings',
    icon: 'i-heroicons-cog-6-tooth',
  },
])

const userMenuItems = [
  [{
    label: 'Profile',
    icon: 'i-heroicons-user-circle',
    click: () => {}
  }],
  [{
    label: 'Settings',
    icon: 'i-heroicons-cog-6-tooth',
    click: () => {}
  }],
  [{
    label: 'Logout',
    icon: 'i-heroicons-arrow-right-on-rectangle',
    click: () => logout()
  }]
]

const inventoryStore = useInventoryStore()

const currentPageTitle = computed(() => {
  const item = navigation.value.find(item => isActive(item.href))
  return item?.name || 'Dashboard'
})

const isActive = (href: string) => {
  if (href === '/') {
    return route.path === '/'
  }
  return route.path.startsWith(href)
}
</script>
