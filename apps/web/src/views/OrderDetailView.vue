<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../lib/api'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const order = ref<any>(null)
const loading = ref(false)
const error = ref('')

const isAuthenticated = computed(() => authStore.isAuthenticated)

function formatPrice(value: number) {
  return new Intl.NumberFormat('id-ID').format(value)
}

async function loadOrder() {
  if (!isAuthenticated.value) {
    order.value = null
    return
  }

  loading.value = true
  error.value = ''

  try {
    const response = await api.get(`/api/orders/${route.params.id}`)
    order.value = response.data
  } catch (err: any) {
    error.value = err?.response?.data?.error?.message || 'Gagal memuat detail pesanan.'
  } finally {
    loading.value = false
  }
}

function goLogin() {
  router.push('/login')
}

onMounted(loadOrder)
</script>

<template>
  <section class="page">
    <div v-if="!isAuthenticated" class="page-card">
      <p>Kamu harus masuk terlebih dahulu untuk melihat detail pesanan.</p>
      <button class="nav-button" @click="goLogin">Ke halaman masuk</button>
    </div>

    <div v-else-if="loading" class="page-card">Memuat detail pesanan...</div>
    <div v-else-if="error" class="page-card">{{ error }}</div>

    <div v-else-if="order" class="page-card">
      <h1 class="page-title">Detail Pesanan</h1>
      <p class="page-subtitle">ID: {{ order.id }}</p>
      <p><strong>Status:</strong> {{ order.status }}</p>
      <p><strong>Subtotal:</strong> Rp {{ formatPrice(order.subtotal_cents ?? 0) }}</p>
      <p><strong>Diskon:</strong> Rp {{ formatPrice(order.discount_cents ?? 0) }}</p>
      <p><strong>Total:</strong> Rp {{ formatPrice(order.total_cents ?? 0) }}</p>

      <div style="margin-top: 1rem; display: grid; gap: 0.75rem;">
        <article
          v-for="item in order.items"
          :key="item.product_id"
          style="border: 1px solid #e5e7eb; border-radius: 12px; padding: 1rem;"
        >
          <h2 style="margin-top: 0;">{{ item.product_name_snapshot }}</h2>
          <p><strong>Jumlah:</strong> {{ item.qty }}</p>
          <p><strong>Harga satuan:</strong> Rp {{ formatPrice(item.price_cents) }}</p>
          <p><strong>Subtotal item:</strong> Rp {{ formatPrice(item.price_cents * item.qty) }}</p>
        </article>
      </div>
    </div>
  </section>
</template>