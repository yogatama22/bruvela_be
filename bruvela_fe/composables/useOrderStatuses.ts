export const useOrderStatuses = () => {
  const config = useRuntimeConfig()
  const token = useCookie('auth_token')

  const fetchOrderStatuses = async () => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/order-statuses`, {
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const createOrderStatus = async (orderStatus: any) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/order-statuses`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token.value}`
        },
        body: orderStatus
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const updateOrderStatus = async (id: string, orderStatus: any) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/order-statuses/${id}`, {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${token.value}`
        },
        body: orderStatus
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const deleteOrderStatus = async (id: string) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/order-statuses/${id}`, {
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
    fetchOrderStatuses,
    createOrderStatus,
    updateOrderStatus,
    deleteOrderStatus
  }
}
