<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
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

const route = useRoute()
const router = useRouter()

const product = ref<Product | null>(null)
const loading = ref(false)
const error = ref('')

async function loadProduct() {
  loading.value = true
  error.value = ''

  try {
    const response = await api.get(`/api/products/${route.params.id}`)
    product.value = response.data
  } catch {
    error.value = 'Failed to load product.'
  } finally {
    loading.value = false
  }
}

async function addToCart() {
  if (!product.value) return

  try {
    await api.post('/api/cart/items', {
      product_id: product.value.id,
      qty: 1,
    })
    router.push('/cart')
  } catch {
    error.value = 'Failed to add item to cart.'
  }
}

onMounted(loadProduct)
</script>

<template>
  <section class="page">
    <div v-if="loading" class="page-card">Loading product...</div>
    <div v-else-if="error" class="page-card">{{ error }}</div>

    <div v-else-if="product" class="page-card">
      <h1 class="page-title">{{ product.name }}</h1>
      <p class="page-subtitle">{{ product.description }}</p>
      <p>Price: {{ product.price_cents }}</p>
      <p>Stock: {{ product.stock }}</p>
      <button @click="addToCart">Add to cart</button>
    </div>
  </section>
</template>