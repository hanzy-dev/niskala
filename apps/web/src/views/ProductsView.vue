<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { api } from '../lib/api'

type Product = {
  id: string
  name: string
  description: string
  price_cents: number
  stock: number
  category: string
  image_url: string
  is_active: boolean
}

const products = ref<Product[]>([])
const loading = ref(false)
const error = ref('')

async function loadProducts() {
  loading.value = true
  error.value = ''

  try {
    const response = await api.get('/api/products')
    products.value = response.data.items
  } catch {
    error.value = 'Gagal memuat daftar produk.'
  } finally {
    loading.value = false
  }
}

onMounted(loadProducts)
</script>

<template>
  <section class="page">
    <div class="page-card">
      <h1 class="page-title">Produk</h1>
      <p class="page-subtitle">Katalog produk dimuat dari API.</p>
    </div>

    <div v-if="loading" class="page-card">Memuat produk...</div>
    <div v-else-if="error" class="page-card">{{ error }}</div>

    <div v-else class="page" style="grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));">
      <article v-for="product in products" :key="product.id" class="page-card">
        <h2 style="margin-top: 0;">{{ product.name }}</h2>
        <p class="page-subtitle">{{ product.description }}</p>
        <p>Harga: {{ product.price_cents }}</p>
        <p>Stok: {{ product.stock }}</p>
        <RouterLink :to="`/products/${product.id}`">Lihat detail</RouterLink>
      </article>
    </div>
  </section>
</template>