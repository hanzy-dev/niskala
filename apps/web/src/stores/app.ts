import { defineStore } from 'pinia'

type AppStatus = 'idle' | 'loading' | 'ready'

export const useAppStore = defineStore('app', {
  state: () => ({
    status: 'idle' as AppStatus,
    pageTitle: 'Niskala',
    debugUserId: localStorage.getItem('debug_user_id') || 'user_123',
    debugUserRole: localStorage.getItem('debug_user_role') || 'user',
  }),
  actions: {
    setStatus(status: AppStatus) {
      this.status = status
    },
    setPageTitle(title: string) {
      this.pageTitle = title
    },
    setDebugUser(userId: string, role: string) {
      this.debugUserId = userId
      this.debugUserRole = role

      localStorage.setItem('debug_user_id', userId)
      localStorage.setItem('debug_user_role', role)
    },
  },
})