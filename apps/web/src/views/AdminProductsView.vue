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
const saving = ref('')
const error = ref('')
const success = ref('')

function formatPrice(value: number) {
  return new Intl.NumberFormat('id-ID').format(value)
}

async function loadProducts() {
  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const response = await api.get('/api/products')
    products.value = response.data.items
  } catch (err: any) {
    error.value = err?.response?.data?.error?.message || 'Gagal memuat daftar produk admin.'
  } finally {
    loading.value = false
  }
}

async function updateStock(productId: string, stock: number) {
  saving.value = productId
  error.value = ''
  success.value = ''

  try {
    await api.patch(`/api/admin/products/${productId}/stock`, { stock })
    success.value = `Stok produk ${productId} berhasil diperbarui.`
    await loadProducts()
  } catch (err: any) {
    error.value = err?.response?.data?.error?.message || 'Gagal memperbarui stok produk.'
  } finally {
    saving.value = ''
  }
}

function increaseStock(product: Product) {
  updateStock(product.id, product.stock + 1)
}

function decreaseStock(product: Product) {
  if (product.stock <= 0) return
  updateStock(product.id, product.stock - 1)
}

onMounted(loadProducts)
</script>

<template>
  <section class="page">
    <div class="page-card">
      <h1 class="page-title">Admin Produk</h1>
      <p class="page-subtitle">
        Kelola produk dan stok melalui backend admin yang terproteksi.
      </p>

      <div style="margin-top: 1rem;">
        <RouterLink class="nav-button" to="/admin/products/new">Tambah produk baru</RouterLink>
      </div>
    </div>

    <div v-if="loading" class="page-card">Memuat daftar produk admin...</div>
    <div v-else-if="error" class="page-card">{{ error }}</div>
    <div v-else class="page">
      <div v-if="success" class="page-card">{{ success }}</div>

      <article v-for="product in products" :key="product.id" class="page-card">
        <p class="page-subtitle" style="margin-bottom: 0.25rem;">{{ product.category }}</p>
        <h2 style="margin-top: 0">{{ product.name }}</h2>
        <p class="page-subtitle">{{ product.description }}</p>

        <p><strong>ID:</strong> {{ product.id }}</p>
        <p><strong>Harga:</strong> Rp {{ formatPrice(product.price_cents) }}</p>
        <p><strong>Stok:</strong> {{ product.stock }}</p>
        <p><strong>Status:</strong> {{ product.is_active ? 'Aktif' : 'Nonaktif' }}</p>

        <div style="display: flex; gap: 0.75rem; flex-wrap: wrap; margin-top: 1rem;">
          <button class="nav-button" :disabled="saving === product.id || product.stock <= 0" @click="decreaseStock(product)">
            {{ saving === product.id ? 'Menyimpan...' : 'Kurangi stok' }}
          </button>

          <button class="nav-button" :disabled="saving === product.id" @click="increaseStock(product)">
            {{ saving === product.id ? 'Menyimpan...' : 'Tambah stok' }}
          </button>

          <RouterLink class="nav-button" :to="`/products/${product.id}`">Lihat sebagai pengguna</RouterLink>
        </div>
      </article>
    </div>
  </section>
</template>