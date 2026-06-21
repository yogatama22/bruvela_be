export const useRecipes = () => {
  const api = useApi()

  const fetchRecipes = async () => {
    return await api.get('/recipes')
  }

  const fetchRecipesByProductId = async (productId: string) => {
    return await api.get(`/recipes/product/${productId}`)
  }

  const createRecipe = async (recipeData: any) => {
    return await api.post('/recipes', recipeData)
  }

  const updateRecipe = async (id: string, recipeData: any) => {
    return await api.put(`/recipes/${id}`, recipeData)
  }

  const deleteRecipe = async (id: string) => {
    return await api.delete(`/recipes/${id}`)
  }

  const deleteRecipesByProductId = async (productId: string) => {
    return await api.delete(`/recipes/product/${productId}`)
  }

  const calculateProduction = async (items: any[]) => {
    return await api.post('/recipes/calculator', { items })
  }

  return {
    fetchRecipes,
    fetchRecipesByProductId,
    createRecipe,
    updateRecipe,
    deleteRecipe,
    deleteRecipesByProductId,
    calculateProduction
  }
}
