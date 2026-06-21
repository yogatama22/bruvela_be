export const useAuth = () => {
  const config = useRuntimeConfig()
  const router = useRouter()
  const token = useCookie('auth_token', {
    maxAge: 60 * 60 * 24 * 7 // 7 days
  })
  const user = useState<any>('user', () => null)

  const login = async (email: string, password: string) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/auth/login`, {
        method: 'POST',
        body: { email, password }
      })

      token.value = response.token
      user.value = response.user

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    router.push('/login')
  }

  const getMe = async () => {
    if (!token.value) return { data: null, error: 'No token' }

    try {
      const response = await $fetch(`${config.public.apiBase}/auth/me`, {
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      })

      user.value = response
      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const isAuthenticated = computed(() => !!token.value)

  return {
    login,
    logout,
    getMe,
    isAuthenticated,
    user,
    token
  }
}
