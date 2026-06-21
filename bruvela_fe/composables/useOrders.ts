export const useOrders = () => {
  const api = useApi()

  const fetchOrders = async (filters?: any) => {
    const queryParams = new URLSearchParams()
    if (filters?.prodStatus) queryParams.append('prod_status', filters.prodStatus)
    if (filters?.payStatus) queryParams.append('pay_status', filters.payStatus)
    if (filters?.batchId) queryParams.append('batch_id', filters.batchId)
    
    const query = queryParams.toString() ? `?${queryParams.toString()}` : ''
    return await api.get(`/orders${query}`)
  }

  const fetchOrderById = async (id: string) => {
    return await api.get(`/orders/${id}`)
  }

  const createOrder = async (orderData: any) => {
    return await api.post('/orders', orderData)
  }

  const updateOrder = async (id: string, orderData: any) => {
    return await api.put(`/orders/${id}`, orderData)
  }

  const updateOrderStatus = async (id: string, status: string) => {
    return await api.patch(`/orders/${id}/status`, { status })
  }

  const updatePaymentStatus = async (id: string, payStatus: string) => {
    return await api.patch(`/orders/${id}/pay`, { pay_status: payStatus })
  }

  const deleteOrder = async (id: string) => {
    return await api.delete(`/orders/${id}`)
  }

  return {
    fetchOrders,
    fetchOrderById,
    createOrder,
    updateOrder,
    updateOrderStatus,
    updatePaymentStatus,
    deleteOrder
  }
}
