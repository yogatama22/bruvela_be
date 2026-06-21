export const useFinance = () => {
  const api = useApi()

  /**
   * Fetch journal entries for a batch.
   * Backend route: GET /finance/journal?batch_id=...
   */
  const getJournal = async (batchId: string) => {
    if (!batchId) return { data: null, error: 'batch_id required' }
    return await api.get(`/finance/journal?batch_id=${batchId}`)
  }

  /**
   * Create a new journal entry for a batch.
   * Backend route: POST /finance/journal
   */
  const createJournalEntry = async (entryData: any) => {
    return await api.post('/finance/journal', entryData)
  }

  /**
   * Fetch finance summary (total income, expense, balance) for a batch.
   * Backend route: GET /finance/summary?batch_id=...
   */
  const getFinanceSummary = async (batchId: string) => {
    if (!batchId) return { data: null, error: 'batch_id required' }
    return await api.get(`/finance/summary?batch_id=${batchId}`)
  }

  /**
   * Profit & Loss report. Backend route not yet implemented.
   */
  const fetchProfitLoss = async (batchId?: string) => {
    if (!batchId) return { data: null, error: 'batch_id required' }
    return await api.get(`/finance/profit-loss?batch_id=${batchId}`)
  }

  /**
   * Batch helpers re-exposed here for convenience.
   */
  const fetchBatches = async () => {
    return await api.get('/batches')
  }
  const createBatch = async (batchData: any) => {
    return await api.post('/batches', batchData)
  }
  const closeBatch = async (id: string) => {
    return await api.patch(`/batches/${id}/close`, {})
  }

  return {
    getJournal,
    createJournalEntry,
    getFinanceSummary,
    fetchProfitLoss,
    fetchBatches,
    createBatch,
    closeBatch
  }
}
