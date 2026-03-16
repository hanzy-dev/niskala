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
    error.value = 'Failed to load products.'
  } finally {
    loading.value = false
  }
}

onMounted(loadProducts)
</script>

<template>
  <section class="page">
    <div class="page-card">
      <h1 class="page-title">Products</h1>
      <p class="page-subtitle">Catalog loaded from the API.</p>
    </div>

    <div v-if="loading" class="page-card">Loading products...</div>
    <div v-else-if="error" class="page-card">{{ error }}</div>

    <div v-else class="page" style="grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));">
      <article v-for="product in products" :key="product.id" class="page-card">
        <h2 style="margin-top: 0;">{{ product.name }}</h2>
        <p class="page-subtitle">{{ product.description }}</p>
        <p>Price: {{ product.price_cents }}</p>
        <p>Stock: {{ product.stock }}</p>
        <RouterLink :to="`/products/${product.id}`">View detail</RouterLink>
      </article>
    </div>
  </section>
</template>