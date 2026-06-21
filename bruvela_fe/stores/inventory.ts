import { defineStore } from 'pinia'

export const useInventoryStore = defineStore('inventory', {
  state: () => ({
    criticalStockCount: 8
  }),

  actions: {
    async fetchCriticalStock() {
    }
  }
})
