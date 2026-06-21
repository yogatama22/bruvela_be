export const useIngredients = () => {
  const api = useApi()

  const fetchIngredients = async () => {
    return await api.get('/ingredients')
  }

  const fetchIngredientsWithEstimation = async (batchId?: string) => {
    const query = batchId ? `?batch_id=${batchId}` : ''
    return await api.get(`/ingredients/estimation${query}`)
  }

  const fetchIngredientById = async (id: string) => {
    return await api.get(`/ingredients/${id}`)
  }

  const fetchCriticalStock = async () => {
    return await api.get('/ingredients/alerts')
  }

  const createIngredient = async (ingredientData: any) => {
    return await api.post('/ingredients', ingredientData)
  }

  const updateIngredient = async (id: string, ingredientData: any) => {
    return await api.put(`/ingredients/${id}`, ingredientData)
  }

  const deleteIngredient = async (id: string) => {
    return await api.delete(`/ingredients/${id}`)
  }

  const fetchPurchases = async () => {
    return await api.get('/ingredient-purchases')
  }

  const createPurchase = async (purchaseData: any) => {
    return await api.post('/ingredient-purchases', purchaseData)
  }

  return {
    fetchIngredients,
    fetchIngredientsWithEstimation,
    fetchIngredientById,
    fetchCriticalStock,
    createIngredient,
    updateIngredient,
    deleteIngredient,
    fetchPurchases,
    createPurchase
  }
}
