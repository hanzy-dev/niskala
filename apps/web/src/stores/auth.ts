import type { Session, User } from '@supabase/supabase-js'
import { defineStore } from 'pinia'
import { supabase } from '../lib/supabase'

type AuthStatus = 'idle' | 'loading' | 'authenticated' | 'unauthenticated'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    status: 'idle' as AuthStatus,
    session: null as Session | null,
    user: null as User | null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.session,
    userEmail: (state) => state.user?.email ?? '',
  },

  actions: {
    async bootstrap() {
      this.status = 'loading'

      const { data } = await supabase.auth.getSession()
      this.session = data.session
      this.user = data.session?.user ?? null
      this.status = data.session ? 'authenticated' : 'unauthenticated'

      supabase.auth.onAuthStateChange((_event, session) => {
        this.session = session
        this.user = session?.user ?? null
        this.status = session ? 'authenticated' : 'unauthenticated'
      })
    },

    async signInWithGoogle() {
      await supabase.auth.signInWithOAuth({
        provider: 'google',
        options: {
          redirectTo: window.location.origin,
        },
      })
    },

    async signOut() {
      await supabase.auth.signOut()
      this.session = null
      this.user = null
      this.status = 'unauthenticated'
    },
  },
})