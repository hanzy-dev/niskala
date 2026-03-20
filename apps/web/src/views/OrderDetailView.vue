<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '../lib/api'

const route = useRoute()
const order = ref<any>(null)
const loading = ref(false)
const error = ref('')

async function loadOrder() {
  loading.value = true
  error.value = ''

  try {
    const response = await api.get(`/api/orders/${route.params.id}`)
    order.value = response.data
  } catch {
    error.value = 'Gagal memuat detail pesanan.'
  } finally {
    loading.value = false
  }
}

onMounted(loadOrder)
</script>

<template>
  <section class="page">
    <div v-if="loading" class="page-card">Memuat detail pesanan...</div>
    <div v-else-if="error" class="page-card">{{ error }}</div>

    <div v-else-if="order" class="page-card">
      <h1 class="page-title">Detail Pesanan</h1>
      <p class="page-subtitle">Status: {{ order.status }}</p>
      <p>Total: {{ order.total_cents }}</p>

      <ul class="page-list">
        <li v-for="item in order.items" :key="item.product_id">
          {{ item.product_name_snapshot }} × {{ item.qty }}
        </li>
      </ul>
    </div>
  </section>
</template>