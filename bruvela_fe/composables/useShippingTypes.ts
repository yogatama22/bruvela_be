export const useShippingTypes = () => {
  const config = useRuntimeConfig()
  const token = useCookie('auth_token')

  const fetchShippingTypes = async () => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/shipping-types`, {
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const createShippingType = async (shippingType: any) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/shipping-types`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token.value}`
        },
        body: shippingType
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const updateShippingType = async (id: string, shippingType: any) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/shipping-types/${id}`, {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${token.value}`
        },
        body: shippingType
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const deleteShippingType = async (id: string) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/shipping-types/${id}`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  return {
    fetchShippingTypes,
    createShippingType,
    updateShippingType,
    deleteShippingType
  }
}
