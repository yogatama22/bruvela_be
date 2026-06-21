export const useCompanySettings = () => {
  const api = useApi()

  const fetchSettings = async () => {
    return await api.get('/company-settings')
  }

  const updateSettings = async (settings: any) => {
    return await api.put('/company-settings', settings)
  }

  return {
    fetchSettings,
    updateSettings
  }
}
