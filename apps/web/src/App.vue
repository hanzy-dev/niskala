<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router';
</script>

<template>
  <div class="app-shell">
    <header class="app-header">
      <div class="app-header__inner">
        <RouterLink class="brand" to="/">Niskala</RouterLink>

        <nav class="nav">
          <RouterLink to="/products">Products</RouterLink>
          <RouterLink to="/cart">Cart</RouterLink>
          <RouterLink to="/orders">Orders</RouterLink>
          <RouterLink to="/admin/products">Admin</RouterLink>
          <RouterLink to="/login">Login</RouterLink>
        </nav>
      </div>
    </header>

    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>