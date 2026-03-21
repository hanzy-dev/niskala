<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { api } from '../lib/api'

const router = useRouter()

const loading = ref(false)
const error = ref('')
const success = ref('')

const form = reactive({
  id: '',
  name: '',
  description: '',
  price_cents: 0,
  stock: 0,
  category: '',
  image_url: '',
  is_active: true,
})

async function submitProduct() {
  loading.value = true
  error.value = ''
  success.value = ''

  try {
    await api.post('/api/admin/products', {
      id: form.id,
      name: form.name,
      description: form.description,
      price_cents: Number(form.price_cents),
      stock: Number(form.stock),
      category: form.category,
      image_url: form.image_url,
      is_active: form.is_active,
    })

    success.value = 'Produk berhasil dibuat.'
    router.push('/admin/products')
  } catch (err: any) {
    error.value = err?.response?.data?.error?.message || 'Gagal membuat produk.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="page">
    <div class="page-card">
      <h1 class="page-title">Form Produk Admin</h1>
      <p class="page-subtitle">Tambahkan produk baru melalui backend admin.</p>
    </div>

    <form class="page-card" style="display: grid; gap: 1rem;" @submit.prevent="submitProduct">
      <label>
        <div>ID Produk</div>
        <input v-model="form.id" type="text" placeholder="contoh: prod_4" />
      </label>

      <label>
        <div>Nama Produk</div>
        <input v-model="form.name" type="text" placeholder="Masukkan nama produk" />
      </label>

      <label>
        <div>Deskripsi</div>
        <textarea v-model="form.description" rows="4" placeholder="Masukkan deskripsi produk" />
      </label>

      <label>
        <div>Harga (dalam cents)</div>
        <input v-model="form.price_cents" type="number" min="0" />
      </label>

      <label>
        <div>Stok</div>
        <input v-model="form.stock" type="number" min="0" />
      </label>

      <label>
        <div>Kategori</div>
        <input v-model="form.category" type="text" placeholder="contoh: stationery" />
      </label>

      <label>
        <div>URL Gambar</div>
        <input v-model="form.image_url" type="text" placeholder="opsional" />
      </label>

      <label style="display: flex; align-items: center; gap: 0.5rem;">
        <input v-model="form.is_active" type="checkbox" />
        Produk aktif
      </label>

      <p v-if="error">{{ error }}</p>
      <p v-if="success">{{ success }}</p>

      <div style="display: flex; gap: 0.75rem; flex-wrap: wrap;">
        <button class="nav-button" type="submit" :disabled="loading">
          {{ loading ? 'Menyimpan...' : 'Simpan produk' }}
        </button>

        <button class="nav-button" type="button" @click="router.push('/admin/products')">
          Kembali
        </button>
      </div>
    </form>
  </section>
</template>