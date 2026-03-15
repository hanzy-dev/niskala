import { createRouter, createWebHistory } from 'vue-router'

import AdminProductFormView from '../views/AdminProductFormView.vue'
import AdminProductsView from '../views/AdminProductsView.vue'
import CartView from '../views/CartView.vue'
import CheckoutSuccessView from '../views/CheckoutSuccessView.vue'
import CheckoutView from '../views/CheckoutView.vue'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import OrderDetailView from '../views/OrderDetailView.vue'
import OrdersView from '../views/OrdersView.vue'
import ProductDetailView from '../views/ProductDetailView.vue'
import ProductsView from '../views/ProductsView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    { path: '/login', name: 'login', component: LoginView },
    { path: '/products', name: 'products', component: ProductsView },
    { path: '/products/:id', name: 'product-detail', component: ProductDetailView },
    { path: '/cart', name: 'cart', component: CartView },
    { path: '/checkout', name: 'checkout', component: CheckoutView },
    { path: '/checkout/success', name: 'checkout-success', component: CheckoutSuccessView },
    { path: '/orders', name: 'orders', component: OrdersView },
    { path: '/orders/:id', name: 'order-detail', component: OrderDetailView },
    { path: '/admin/products', name: 'admin-products', component: AdminProductsView },
    { path: '/admin/products/new', name: 'admin-product-create', component: AdminProductFormView },
    { path: '/admin/products/:id/edit', name: 'admin-product-edit', component: AdminProductFormView },
  ],
})

export default router