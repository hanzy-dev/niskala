import { defineStore } from 'pinia'

type AppStatus = 'idle' | 'loading' | 'ready'

export const useAppStore = defineStore('app', {
  state: () => ({
    status: 'idle' as AppStatus,
    pageTitle: 'Niskala',
  }),
  actions: {
    setStatus(status: AppStatus) {
      this.status = status
    },
    setPageTitle(title: string) {
      this.pageTitle = title
    },
  },
})