export const usePaymentStatuses = () => {
  const config = useRuntimeConfig()
  const token = useCookie('auth_token')

  const fetchPaymentStatuses = async () => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/payment-statuses`, {
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const createPaymentStatus = async (paymentStatus: any) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/payment-statuses`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token.value}`
        },
        body: paymentStatus
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const updatePaymentStatus = async (id: string, paymentStatus: any) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/payment-statuses/${id}`, {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${token.value}`
        },
        body: paymentStatus
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const deletePaymentStatus = async (id: string) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/payment-statuses/${id}`, {
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
    fetchPaymentStatuses,
    createPaymentStatus,
    updatePaymentStatus,
    deletePaymentStatus
  }
}
