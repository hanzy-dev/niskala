<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { api } from '../lib/api'
import { useAuthStore } from '../stores/auth'

type Order = {
  id: string
  status: string
  total_cents: number
}

const authStore = useAuthStore()
const router = useRouter()

const orders = ref<Order[]>([])
const loading = ref(false)
const error = ref('')

const isAuthenticated = computed(() => authStore.isAuthenticated)

function formatPrice(value: number) {
  return new Intl.NumberFormat('id-ID').format(value)
}

async function loadOrders() {
  if (!isAuthenticated.value) {
    orders.value = []
    return
  }

  loading.value = true
  error.value = ''

  try {
    const response = await api.get('/api/orders')
    orders.value = response.data.items
  } catch (err: any) {
    error.value = err?.response?.data?.error?.message || 'Gagal memuat daftar pesanan.'
  } finally {
    loading.value = false
  }
}

function goLogin() {
  router.push('/login')
}

onMounted(loadOrders)
</script>

<template>
  <section class="page">
    <div class="page-card">
      <h1 class="page-title">Pesanan</h1>
      <p class="page-subtitle">Riwayat pesanan dimuat dari API.</p>
    </div>

    <div v-if="!isAuthenticated" class="page-card">
      <p>Kamu harus masuk terlebih dahulu untuk melihat pesanan.</p>
      <button class="nav-button" @click="goLogin">Ke halaman masuk</button>
    </div>

    <div v-else-if="loading" class="page-card">Memuat pesanan...</div>
    <div v-else-if="error" class="page-card">{{ error }}</div>

    <div v-else class="page">
      <div v-if="orders.length === 0" class="page-card">Belum ada pesanan.</div>

      <article v-for="order in orders" :key="order.id" class="page-card">
        <h2 style="margin-top: 0;">Pesanan {{ order.id }}</h2>
        <p class="page-subtitle">Status: {{ order.status }}</p>
        <p><strong>Total:</strong> Rp {{ formatPrice(order.total_cents) }}</p>
        <RouterLink :to="`/orders/${order.id}`">Lihat detail pesanan</RouterLink>
      </article>
    </div>
  </section>
</template>