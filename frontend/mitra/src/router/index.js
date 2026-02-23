import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'

const routes = [
  { path: '/login', name: 'Login', component: () => import('../pages/Login.vue'), meta: { guest: true } },
  { path: '/register', name: 'Register', component: () => import('../pages/Register.vue'), meta: { guest: true } },

  {
    path: '/',
    component: () => import('../layouts/MitraLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', name: 'Dashboard', component: () => import('../pages/Dashboard.vue') },
      { path: 'outlets', name: 'Outlets', component: () => import('../pages/Outlets.vue') },
      { path: 'outlets/:id', name: 'OutletDetail', component: () => import('../pages/OutletDetail.vue') },
      { path: 'outlets/:id/apply', name: 'ApplyForm', component: () => import('../pages/ApplyForm.vue') },
      { path: 'applications', name: 'Applications', component: () => import('../pages/Applications.vue') },
      { path: 'applications/:id', name: 'ApplicationDetail', component: () => import('../pages/ApplicationDetail.vue') },
      { path: 'invoices', name: 'Invoices', component: () => import('../pages/Invoices.vue') },
      { path: 'invoices/:id', name: 'InvoiceDetail', component: () => import('../pages/InvoiceDetail.vue') },
      { path: 'agreements', name: 'Agreements', component: () => import('../pages/Agreements.vue') },
      { path: 'locations', name: 'Locations', component: () => import('../pages/Locations.vue') },
      { path: 'ebooks', name: 'Ebooks', component: () => import('../pages/Ebooks.vue') },
      { path: 'ebooks/:id', name: 'EbookDetail', component: () => import('../pages/EbookDetail.vue') },
      { path: 'ebooks/:id/read', name: 'EbookReader', component: () => import('../pages/EbookReader.vue') },
      { path: 'ebook-orders', name: 'EbookOrders', component: () => import('../pages/EbookOrders.vue') },
      { path: 'settings', name: 'Settings', component: () => import('../pages/Settings.vue') },
    ],
  },

  { path: '/:pathMatch(.*)*', name: 'NotFound', component: () => import('../pages/NotFound.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

function isTokenExpired() {
  const token = localStorage.getItem('mitra_token')
  if (!token) return true
  try {
    const parts = token.split('.')
    if (parts.length !== 3) return true
    const payload = JSON.parse(atob(parts[1]))
    if (!payload.exp) return true
    return Date.now() >= payload.exp * 1000
  } catch {
    return true
  }
}

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()

  if (to.meta.requiresAuth && auth.isAuthenticated && isTokenExpired()) {
    auth.logout()
    const toast = useToastStore()
    toast.error('Sesi telah berakhir, silakan login kembali')
    return next({ name: 'Login', query: { redirect: to.fullPath } })
  }

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return next({ name: 'Login', query: { redirect: to.fullPath } })
  }
  if (to.meta.guest && auth.isAuthenticated) {
    return next({ name: 'Dashboard' })
  }

  next()
})

export default router
