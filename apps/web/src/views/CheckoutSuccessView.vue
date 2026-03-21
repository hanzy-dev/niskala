<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

const route = useRoute()

function formatPrice(value: number) {
  return new Intl.NumberFormat('id-ID').format(value)
}

const orderId = computed(() => String(route.query.order_id ?? ''))
const total = computed(() => Number(route.query.total ?? 0))
const fallbackUsed = computed(() => String(route.query.fallback ?? '0') === '1')
</script>

<template>
  <section class="page">
    <div class="page-card">
      <h1 class="page-title">Checkout Berhasil</h1>
      <p class="page-subtitle">Pesanan kamu berhasil dibuat.</p>

      <div style="margin-top: 1rem; display: grid; gap: 0.75rem;">
        <p><strong>ID Pesanan:</strong> {{ orderId || '-' }}</p>
        <p><strong>Total:</strong> Rp {{ formatPrice(total) }}</p>
        <p v-if="fallbackUsed">
          <strong>Catatan:</strong> Perhitungan checkout menggunakan fallback harga normal karena layanan pricing sedang tidak tersedia.
        </p>
      </div>

      <div style="display: flex; gap: 0.75rem; flex-wrap: wrap; margin-top: 1.25rem;">
        <RouterLink to="/orders">Lihat semua pesanan</RouterLink>
        <RouterLink to="/products">Kembali ke produk</RouterLink>
      </div>
    </div>
  </section>
</template>