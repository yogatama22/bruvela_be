export const useStockLogs = () => {
  const api = useApi()

  const fetchStockLogs = async (filters?: any) => {
    const queryParams = new URLSearchParams()
    if (filters?.batchId) queryParams.append('batch_id', filters.batchId)
    if (filters?.ingredientId) queryParams.append('ingredient_id', filters.ingredientId)
    if (filters?.logType) queryParams.append('log_type', filters.logType)
    if (filters?.referenceType) queryParams.append('reference_type', filters.referenceType)

    const query = queryParams.toString() ? `?${queryParams.toString()}` : ''
    return await api.get(`/stock-logs${query}`)
  }

  return {
    fetchStockLogs
  }
}
