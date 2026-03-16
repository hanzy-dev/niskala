<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '../stores/app'

const appStore = useAppStore()
const router = useRouter()

const userId = ref(appStore.debugUserId)
const role = ref(appStore.debugUserRole)

function saveSession() {
  appStore.setDebugUser(userId.value, role.value)
  router.push('/products')
}
</script>

<template>
  <section class="page">
    <div class="page-card">
      <h1 class="page-title">Login</h1>
      <p class="page-subtitle">
        Debug auth placeholder for local development.
      </p>

      <div style="display: grid; gap: 0.75rem; max-width: 360px; margin-top: 1rem;">
        <input v-model="userId" placeholder="User ID" />
        <select v-model="role">
          <option value="user">user</option>
          <option value="admin">admin</option>
        </select>
        <button @click="saveSession">Save debug session</button>
      </div>
    </div>
  </section>
</template>