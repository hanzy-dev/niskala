<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()

const isAuthenticated = computed(() => authStore.isAuthenticated)

async function handleGoogleLogin() {
  await authStore.signInWithGoogle()
}
</script>

<template>
  <section class="page">
    <div class="page-card">
      <h1 class="page-title">Masuk</h1>
      <p class="page-subtitle">
        Masuk dengan Google melalui Supabase untuk mengakses keranjang, checkout, dan pesanan.
      </p>

      <div v-if="isAuthenticated">
        <p class="page-subtitle">Kamu sudah masuk sebagai {{ authStore.userEmail }}</p>
      </div>

      <div v-else>
        <button class="nav-button" @click="handleGoogleLogin">Masuk dengan Google</button>
      </div>
    </div>
  </section>
</template>