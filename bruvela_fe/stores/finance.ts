import { defineStore } from 'pinia'

export const useFinanceStore = defineStore('finance', {
  state: () => ({
    summary: {
      totalRevenue: 0,
      totalExpense: 0,
      grossProfit: 0
    }
  }),

  actions: {
    async fetchSummary() {
    }
  }
})
