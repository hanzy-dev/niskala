<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../lib/api'

type CartItem = {
  product_id: string
  qty: number
}

type Cart = {
  user_id: string
  items: CartItem[]
}

const router = useRouter()
const cart = ref<Cart | null>(null)
const loading = ref(false)
const error = ref('')

async function loadCart() {
  loading.value = true
  error.value = ''

  try {
    const response = await api.get('/api/cart')
    cart.value = response.data
  } catch {
    error.value = 'Failed to load cart.'
  } finally {
    loading.value = false
  }
}

async function removeItem(productId: string) {
  try {
    const response = await api.delete(`/api/cart/items/${productId}`)
    cart.value = response.data
  } catch {
    error.value = 'Failed to remove cart item.'
  }
}

function goCheckout() {
  router.push('/checkout')
}

onMounted(loadCart)
</script>

<template>
  <section class="page">
    <div class="page-card">
      <h1 class="page-title">Cart</h1>
      <p class="page-subtitle">Cart loaded from the API.</p>
    </div>

    <div v-if="loading" class="page-card">Loading cart...</div>
    <div v-else-if="error" class="page-card">{{ error }}</div>

    <div v-else-if="cart" class="page-card">
      <div v-if="cart.items.length === 0">Your cart is empty.</div>

      <ul v-else class="page-list">
        <li v-for="item in cart.items" :key="item.product_id">
          {{ item.product_id }} × {{ item.qty }}
          <button style="margin-left: 0.75rem;" @click="removeItem(item.product_id)">Remove</button>
        </li>
      </ul>

      <button v-if="cart.items.length > 0" style="margin-top: 1rem;" @click="goCheckout">
        Continue to checkout
      </button>
    </div>
  </section>
</template>