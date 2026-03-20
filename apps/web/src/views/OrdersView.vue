<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { api } from '../lib/api'

type Order = {
  id: string
  status: string
  total_cents: number
}

const orders = ref<Order[]>([])
const loading = ref(false)
const error = ref('')

async function loadOrders() {
  loading.value = true
  error.value = ''

  try {
    const response = await api.get('/api/orders')
    orders.value = response.data.items
  } catch {
    error.value = 'Gagal memuat daftar pesanan.'
  } finally {
    loading.value = false
  }
}

onMounted(loadOrders)
</script>

<template>
  <section class="page">
    <div class="page-card">
      <h1 class="page-title">Pesanan</h1>
      <p class="page-subtitle">Riwayat pesanan dimuat dari API.</p>
    </div>

    <div v-if="loading" class="page-card">Memuat pesanan...</div>
    <div v-else-if="error" class="page-card">{{ error }}</div>

    <div v-else class="page-card">
      <div v-if="orders.length === 0">Belum ada pesanan.</div>

      <ul v-else class="page-list">
        <li v-for="order in orders" :key="order.id">
          <RouterLink :to="`/orders/${order.id}`">
            {{ order.id }} — {{ order.status }} — {{ order.total_cents }}
          </RouterLink>
        </li>
      </ul>
    </div>
  </section>
</template>