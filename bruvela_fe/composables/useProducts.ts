export const useProducts = () => {
  const api = useApi()

  const fetchProducts = async () => {
    return await api.get('/products')
  }

  const fetchProductById = async (id: string) => {
    return await api.get(`/products/${id}`)
  }

  const fetchProductRecipe = async (id: string) => {
    return await api.get(`/products/${id}/recipe`)
  }

  const createProduct = async (productData: any) => {
    return await api.post('/products', productData)
  }

  const updateProduct = async (id: string, productData: any) => {
    return await api.put(`/products/${id}`, productData)
  }

  const deleteProduct = async (id: string) => {
    return await api.delete(`/products/${id}`)
  }

  return {
    fetchProducts,
    fetchProductById,
    fetchProductRecipe,
    createProduct,
    updateProduct,
    deleteProduct
  }
}
