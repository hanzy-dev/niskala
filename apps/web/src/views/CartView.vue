<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
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

const router = useRouter()
const authStore = useAuthStore()

const cart = ref<Cart | null>(null)
const loading = ref(false)
const error = ref('')

const isAuthenticated = computed(() => authStore.isAuthenticated)

async function loadCart() {
  if (!isAuthenticated.value) {
    cart.value = null
    return
  }

  loading.value = true
  error.value = ''

  try {
    const response = await api.get('/api/cart')
    cart.value = response.data
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

      <ul v-else class="page-list">
        <li v-for="item in cart.items" :key="item.product_id">
          {{ item.product_id }} × {{ item.qty }}
          <button style="margin-left: 0.75rem;" @click="removeItem(item.product_id)">Hapus</button>
        </li>
      </ul>

      <button v-if="cart.items.length > 0" style="margin-top: 1rem;" @click="goCheckout">
        Lanjut ke checkout
      </button>
    </div>
  </section>
</template>