export const useApi = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase

  const apiFetch = async (endpoint: string, options: any = {}) => {
    const token = useCookie('auth_token')
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...options.headers
    }

    if (token.value) {
      headers['Authorization'] = `Bearer ${token.value}`
    }

    try {
      const response = await $fetch(`${baseURL}${endpoint}`, {
        ...options,
        headers
      })
      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  return {
    get: (endpoint: string, options = {}) => apiFetch(endpoint, { method: 'GET', ...options }),
    post: (endpoint: string, body: any, options = {}) => apiFetch(endpoint, { method: 'POST', body, ...options }),
    put: (endpoint: string, body: any, options = {}) => apiFetch(endpoint, { method: 'PUT', body, ...options }),
    patch: (endpoint: string, body: any, options = {}) => apiFetch(endpoint, { method: 'PATCH', body, ...options }),
    delete: (endpoint: string, options = {}) => apiFetch(endpoint, { method: 'DELETE', ...options })
  }
}
