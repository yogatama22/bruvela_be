import { defineStore } from 'pinia'

export const useOrdersStore = defineStore('orders', {
  state: () => ({
    orders: [],
    loading: false
  }),

  actions: {
    async fetchOrders() {
    }
  }
})
