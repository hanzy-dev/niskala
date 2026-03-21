<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { api } from '../lib/api'
import { useAuthStore } from '../stores/auth'

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
const authStore = useAuthStore()

const product = ref<Product | null>(null)
const loading = ref(false)
const error = ref('')
const success = ref('')

const isAuthenticated = computed(() => authStore.isAuthenticated)

function formatPrice(value: number) {
  return new Intl.NumberFormat('id-ID').format(value)
}

async function loadProduct() {
  loading.value = true
  error.value = ''
  success.value = ''

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

  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }

  try {
    await api.post('/api/cart/items', {
      product_id: product.value.id,
      qty: 1,
    })
    success.value = 'Produk berhasil ditambahkan ke keranjang.'
  } catch {
    error.value = 'Gagal menambahkan produk ke keranjang.'
  }
}

function goToCart() {
  router.push('/cart')
}

onMounted(loadProduct)
</script>

<template>
  <section class="page">
    <div v-if="loading" class="page-card">Memuat detail produk...</div>
    <div v-else-if="error" class="page-card">{{ error }}</div>

    <div v-else-if="product" class="page-card">
      <p class="page-subtitle">{{ product.category }}</p>
      <h1 class="page-title">{{ product.name }}</h1>
      <p class="page-subtitle">{{ product.description }}</p>

      <div style="display: grid; gap: 0.75rem; margin-top: 1rem;">
        <p><strong>Harga:</strong> Rp {{ formatPrice(product.price_cents) }}</p>
        <p><strong>Stok tersedia:</strong> {{ product.stock }}</p>
        <p><strong>Status:</strong> {{ product.is_active ? 'Aktif' : 'Nonaktif' }}</p>
      </div>

      <div style="display: flex; gap: 0.75rem; flex-wrap: wrap; margin-top: 1rem;">
        <button class="nav-button" @click="addToCart">Tambah ke keranjang</button>
        <button class="nav-button" @click="goToCart">Lihat keranjang</button>
      </div>

      <p v-if="success" style="margin-top: 1rem;">{{ success }}</p>
    </div>
  </section>
</template>