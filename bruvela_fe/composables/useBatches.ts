export const useBatches = () => {
  const config = useRuntimeConfig()
  const token = useCookie('auth_token')

  const fetchBatches = async () => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/batches`, {
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const fetchActiveBatch = async () => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/batches/active`, {
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const fetchBatchById = async (id: string) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/batches/${id}`, {
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const createBatch = async (batch: any) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/batches`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token.value}`
        },
        body: batch
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const updateBatch = async (id: string, batch: any) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/batches/${id}`, {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${token.value}`
        },
        body: batch
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const activateBatch = async (id: string) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/batches/${id}/activate`, {
        method: 'PATCH',
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const closeBatch = async (id: string) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/batches/${id}/close`, {
        method: 'PATCH',
        headers: {
          Authorization: `Bearer ${token.value}`
        }
      })

      return { data: response, error: null }
    } catch (error: any) {
      return { data: null, error: error.data || error.message }
    }
  }

  const deleteBatch = async (id: string) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/batches/${id}`, {
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

  const fetchBatchSummary = async (id: string) => {
    try {
      const response: any = await $fetch(`${config.public.apiBase}/batches/${id}/summary`, {
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
    fetchBatches,
    fetchActiveBatch,
    fetchBatchById,
    fetchBatchSummary,
    createBatch,
    updateBatch,
    activateBatch,
    closeBatch,
    deleteBatch
  }
}
