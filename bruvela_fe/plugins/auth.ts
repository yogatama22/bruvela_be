export default defineNuxtPlugin(() => {
  const { isAuthenticated, getMe } = useAuth()
  
  addRouteMiddleware('auth', (to) => {
    // Skip middleware for login page
    if (to.path === '/login') {
      return
    }

    // Check if user is authenticated
    if (!isAuthenticated.value) {
      return navigateTo('/login')
    }
  }, { global: true })

  // Load user data on app start if authenticated
  if (isAuthenticated.value && process.client) {
    getMe()
  }
})
