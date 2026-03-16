<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../lib/api'

const router = useRouter()
const loading = ref(false)
const error = ref('')
const fallbackMessage = ref('')

async function submitCheckout() {
  loading.value = true
  error.value = ''
  fallbackMessage.value = ''

  try {
    const response = await api.post(
      '/api/checkout',
      {},
      {
        headers: {
          'Idempotency-Key': crypto.randomUUID(),
        },
      },
    )

    const order = response.data.order ?? response.data

    if (order.pricing_fallback_used) {
      fallbackMessage.value = 'Discount service is unavailable. Checkout was completed with normal pricing.'
    }

    router.push('/checkout/success')
  } catch (err: any) {
    error.value = err?.response?.data?.error?.message || 'Checkout failed.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="page">
    <div class="page-card">
      <h1 class="page-title">Checkout</h1>
      <p class="page-subtitle">Submit checkout through the API.</p>

      <p v-if="fallbackMessage">{{ fallbackMessage }}</p>
      <p v-if="error">{{ error }}</p>

      <button :disabled="loading" @click="submitCheckout">
        {{ loading ? 'Processing...' : 'Place order' }}
      </button>
    </div>
  </section>
</template>