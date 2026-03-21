<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { api } from '../lib/api'
import { useAuthStore } from '../stores/auth'

type CartItem = {
  product_id: string
  qty: number
}

type Cart = {
  user_id: string
  items: CartItem[]
}

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

const router = useRouter()
const authStore = useAuthStore()

const cart = ref<Cart | null>(null)
const products = ref<Record<string, Product>>({})
const loading = ref(false)
const error = ref('')

const isAuthenticated = computed(() => authStore.isAuthenticated)

function formatPrice(value: number) {
  return new Intl.NumberFormat('id-ID').format(value)
}

async function loadProductsMap() {
  const response = await api.get('/api/products')
  const items: Product[] = response.data.items ?? []
  products.value = Object.fromEntries(items.map((item) => [item.id, item]))
}

async function loadCart() {
  if (!isAuthenticated.value) {
    cart.value = null
    return
  }

  loading.value = true
  error.value = ''

  try {
    const [cartResponse] = await Promise.all([
      api.get('/api/cart'),
      loadProductsMap(),
    ])
    cart.value = cartResponse.data
  } catch (err: any) {
    error.value = err?.response?.data?.error?.message || 'Gagal memuat keranjang.'
  } finally {
    loading.value = false
  }
}

async function removeItem(productId: string) {
  try {
    const response = await api.delete(`/api/cart/items/${productId}`)
    cart.value = response.data
  } catch (err: any) {
    error.value = err?.response?.data?.error?.message || 'Gagal menghapus item dari keranjang.'
  }
}

function goCheckout() {
  router.push('/checkout')
}

function goLogin() {
  router.push('/login')
}

function getProduct(productId: string) {
  return products.value[productId]
}

const estimatedTotal = computed(() => {
  if (!cart.value) return 0

  return cart.value.items.reduce((total, item) => {
    const product = getProduct(item.product_id)
    if (!product) return total
    return total + product.price_cents * item.qty
  }, 0)
})

onMounted(loadCart)
</script>

<template>
  <section class="page">
    <div class="page-card">
      <h1 class="page-title">Keranjang</h1>
      <p class="page-subtitle">Keranjang dimuat dari API.</p>
    </div>

    <div v-if="!isAuthenticated" class="page-card">
      <p>Kamu harus masuk terlebih dahulu untuk melihat keranjang.</p>
      <button class="nav-button" @click="goLogin">Ke halaman masuk</button>
    </div>

    <div v-else-if="loading" class="page-card">Memuat keranjang...</div>
    <div v-else-if="error" class="page-card">{{ error }}</div>

    <div v-else-if="cart" class="page-card">
      <div v-if="cart.items.length === 0">Keranjang kamu masih kosong.</div>

      <div v-else style="display: grid; gap: 1rem;">
        <article
          v-for="item in cart.items"
          :key="item.product_id"
          style="border: 1px solid #e5e7eb; border-radius: 12px; padding: 1rem;"
        >
          <h2 style="margin-top: 0;">
            {{ getProduct(item.product_id)?.name ?? item.product_id }}
          </h2>
          <p class="page-subtitle">
            {{ getProduct(item.product_id)?.description ?? 'Produk belum termuat.' }}
          </p>
          <p><strong>Jumlah:</strong> {{ item.qty }}</p>
          <p>
            <strong>Harga satuan:</strong>
            Rp {{ formatPrice(getProduct(item.product_id)?.price_cents ?? 0) }}
          </p>
          <p>
            <strong>Subtotal:</strong>
            Rp {{ formatPrice((getProduct(item.product_id)?.price_cents ?? 0) * item.qty) }}
          </p>

          <div style="display: flex; gap: 0.75rem; flex-wrap: wrap; margin-top: 0.75rem;">
            <RouterLink :to="`/products/${item.product_id}`">Lihat detail</RouterLink>
            <button class="nav-button" @click="removeItem(item.product_id)">Hapus</button>
          </div>
        </article>

        <div class="page-card" style="padding: 1rem;">
          <p><strong>Perkiraan total:</strong> Rp {{ formatPrice(estimatedTotal) }}</p>
          <button style="margin-top: 1rem;" @click="goCheckout">Lanjut ke checkout</button>
        </div>
      </div>
    </div>
  </section>
</template>