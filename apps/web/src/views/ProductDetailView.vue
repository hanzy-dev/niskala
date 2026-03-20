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
    error.value = 'Gagal memuat detail produk.'
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
    error.value = 'Gagal menambahkan produk ke keranjang.'
  }
}

onMounted(loadProduct)
</script>

<template>
  <section class="page">
    <div v-if="loading" class="page-card">Memuat detail produk...</div>
    <div v-else-if="error" class="page-card">{{ error }}</div>

    <div v-else-if="product" class="page-card">
      <h1 class="page-title">Detail Produk</h1>
      <h2>{{ product.name }}</h2>
      <p class="page-subtitle">{{ product.description }}</p>
      <p>Harga: {{ product.price_cents }}</p>
      <p>Stok: {{ product.stock }}</p>
      <button @click="addToCart">Tambah ke keranjang</button>
    </div>
  </section>
</template>