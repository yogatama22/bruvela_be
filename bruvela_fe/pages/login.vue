<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary-50 to-primary-100 px-4">
    <div class="max-w-md w-full">
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-16 h-16 bg-primary-600 rounded-2xl mb-4">
          <UIcon name="i-heroicons-cake" class="w-10 h-10 text-white" />
        </div>
        <h1 class="text-3xl font-bold text-gray-900">Bruvela Bakehouse</h1>
        <p class="mt-2 text-gray-600">Sistem Manajemen Toko</p>
      </div>

      <UCard>
        <template #header>
          <h2 class="text-xl font-semibold text-gray-900">Login</h2>
          <p class="text-sm text-gray-500 mt-1">Masuk ke akun Anda</p>
        </template>

        <form @submit.prevent="handleLogin" class="space-y-4">
          <UFormGroup label="Email" required>
            <UInput
              v-model="form.email"
              type="email"
              placeholder="admin@bruvela.com"
              size="lg"
              icon="i-heroicons-envelope"
              :disabled="loading"
            />
          </UFormGroup>

          <UFormGroup label="Password" required>
            <UInput
              v-model="form.password"
              type="password"
              placeholder="••••••••"
              size="lg"
              icon="i-heroicons-lock-closed"
              :disabled="loading"
            />
          </UFormGroup>

          <div class="flex items-center justify-between">
            <UCheckbox v-model="rememberMe" label="Ingat saya" />
            <UButton variant="link" color="gray" size="sm" :disabled="loading">
              Lupa password?
            </UButton>
          </div>

          <UButton
            type="submit"
            size="lg"
            block
            :loading="loading"
            :disabled="loading"
          >
            Masuk
          </UButton>
        </form>

        <template #footer>
          <div class="text-center text-sm text-gray-500">
            <p>Test Account:</p>
            <p class="font-mono text-xs mt-1">admin@bruvela.com / admin123</p>
          </div>
        </template>
      </UCard>

      <div class="mt-6 text-center text-sm text-gray-600">
        <p>&copy; 2024 Bruvela Bakehouse. All rights reserved.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false
})

useHead({
  title: 'Login'
})

const toast = useToast()
const { login } = useAuth()

const form = ref({
  email: '',
  password: ''
})

const rememberMe = ref(false)
const loading = ref(false)

const handleLogin = async () => {
  if (!form.value.email || !form.value.password) {
    toast.add({
      title: 'Error',
      description: 'Email dan password harus diisi',
      color: 'red'
    })
    return
  }

  loading.value = true

  const { data, error } = await login(form.value.email, form.value.password)

  loading.value = false

  if (error) {
    toast.add({
      title: 'Login Gagal',
      description: error.error || 'Email atau password salah',
      color: 'red'
    })
    return
  }

  toast.add({
    title: 'Login Berhasil',
    description: `Selamat datang, ${data.user.name || data.user.email}!`,
    color: 'green'
  })

  window.location.href = '/'
}

// Auto-fill for testing
onMounted(() => {
  form.value.email = 'admin@bruvela.com'
  form.value.password = 'admin123'
})
</script>
