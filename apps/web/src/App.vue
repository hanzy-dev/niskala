<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, RouterView, useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const authLabel = computed(() => {
  if (authStore.isAuthenticated) {
    return authStore.userEmail || 'Sudah masuk'
  }

  return 'Masuk'
})

async function handleSignOut() {
  await authStore.signOut()
  router.push('/login')
}
</script>

<template>
  <div class="app-shell">
    <header class="app-header">
      <div class="app-header__inner">
        <RouterLink class="brand" to="/">Niskala</RouterLink>

        <nav class="nav">
          <RouterLink to="/products">Produk</RouterLink>
          <RouterLink to="/cart">Keranjang</RouterLink>
          <RouterLink to="/orders">Pesanan</RouterLink>
          <RouterLink to="/admin/products">Admin</RouterLink>
          <RouterLink to="/login">{{ authLabel }}</RouterLink>
          <button v-if="authStore.isAuthenticated" class="nav-button" @click="handleSignOut">
            Keluar
          </button>
        </nav>
      </div>
    </header>

    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>