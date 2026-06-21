export const useDashboard = () => {
  const api = useApi()

  /**
   * Fetch dashboard summary for a specific batch.
   * @param _unused - kept for backward compatibility
   * @param batchId - UUID of the active batch (required by backend)
   */
  const fetchSummary = async (_unused?: any, batchId?: string) => {
    if (!batchId) {
      return { data: null, error: 'batch_id required' }
    }
    return await api.get(`/dashboard/summary?batch_id=${batchId}`)
  }

  /**
   * Sales chart data per product variant for a batch.
   * Backend endpoint not yet implemented; client-side aggregation in dashboard page is fallback.
   */
  const fetchSalesChart = async (batchId?: string) => {
    if (!batchId) return { data: null, error: 'batch_id required' }
    return await api.get(`/dashboard/charts/sales?batch_id=${batchId}`)
  }

  /**
   * Finance chart data (income/expense/profit trend).
   * Backend endpoint not yet implemented.
   */
  const fetchFinanceChart = async (batchId?: string) => {
    if (!batchId) return { data: null, error: 'batch_id required' }
    return await api.get(`/dashboard/charts/finance?batch_id=${batchId}`)
  }

  return {
    fetchSummary,
    fetchSalesChart,
    fetchFinanceChart
  }
}
