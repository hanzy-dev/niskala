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
      <p>Status: {{ order.status }}</p>
      <p>Total: {{ order.total_cents }}</p>

      <ul class="page-list">
        <li v-for="item in order.items" :key="item.product_id">
          {{ item.product_name_snapshot }} × {{ item.qty }}
        </li>
      </ul>
    </div>
  </section>
</template>