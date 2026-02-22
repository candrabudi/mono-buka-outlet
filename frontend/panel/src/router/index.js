import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'

const routes = [
  // Auth
  { path: '/login', name: 'Login', component: () => import('../pages/Login.vue'), meta: { guest: true } },

  // Admin Dashboard
  {
    path: '/',
    component: () => import('../layouts/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      // Dashboard — all admin roles
      { path: '', name: 'Dashboard', component: () => import('../pages/Dashboard.vue') },





      // Outlets — master, admin
      { path: 'outlets', name: 'Outlets', component: () => import('../pages/outlets/OutletList.vue'), meta: { roles: ['master', 'admin'] } },
      { path: 'outlets/create', name: 'OutletCreate', component: () => import('../pages/outlets/OutletForm.vue'), meta: { roles: ['master', 'admin'] } },
      { path: 'outlets/:id', name: 'OutletEdit', component: () => import('../pages/outlets/OutletForm.vue'), meta: { roles: ['master', 'admin'] } },
      { path: 'outlet-categories', name: 'OutletCategories', component: () => import('../pages/outlets/OutletCategoryList.vue'), meta: { roles: ['master', 'admin'] } },

      // Partnerships — master, admin, finance (view)
      { path: 'partnerships', name: 'Partnerships', component: () => import('../pages/partnerships/PartnershipList.vue'), meta: { roles: ['master', 'admin', 'finance'] } },
      { path: 'partnerships/:id', name: 'PartnershipDetail', component: () => import('../pages/partnerships/PartnershipDetail.vue'), meta: { roles: ['master', 'admin', 'finance'] } },

      // Applications — master, admin
      { path: 'applications', name: 'Applications', component: () => import('../pages/applications/ApplicationList.vue'), meta: { roles: ['master', 'admin'] } },

      // Payments — master, finance
      { path: 'payments', name: 'Payments', component: () => import('../pages/payments/PaymentList.vue'), meta: { roles: ['master', 'finance'] } },

      // Meetings — master, admin
      { path: 'meetings', name: 'Meetings', component: () => import('../pages/meetings/MeetingList.vue'), meta: { roles: ['master', 'admin'] } },

      // Mitra — master, admin
      { path: 'mitra', name: 'Mitra', component: () => import('../pages/mitra/MitraList.vue'), meta: { roles: ['master', 'admin'] } },

      // Leader — master, admin
      { path: 'leader', name: 'Leader', component: () => import('../pages/leader/LeaderList.vue'), meta: { roles: ['master', 'admin'] } },

      // Users — master only
      { path: 'users', name: 'Users', component: () => import('../pages/users/UserList.vue'), meta: { roles: ['master'] } },

      // Settings — master, admin
      { path: 'settings', name: 'Settings', component: () => import('../pages/settings/SettingList.vue'), meta: { roles: ['master', 'admin'] } },

      // Location Submissions — master, admin
      { path: 'locations', name: 'Locations', component: () => import('../pages/locations/LocationList.vue'), meta: { roles: ['master', 'admin'] } },
      { path: 'locations/:id', name: 'LocationDetail', component: () => import('../pages/locations/LocationDetail.vue'), meta: { roles: ['master', 'admin'] } },
    ],
  },

  { path: '/:pathMatch(.*)*', name: 'NotFound', component: () => import('../pages/NotFound.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

/**
 * Check if the JWT token stored in localStorage has expired
 * by decoding the payload and comparing the `exp` claim.
 */
function isTokenExpired() {
  const token = localStorage.getItem('admin_token')
  if (!token) return true

  try {
    const parts = token.split('.')
    if (parts.length !== 3) return true

    const payload = JSON.parse(atob(parts[1]))
    if (!payload.exp) return true

    // exp is in seconds, Date.now() is in ms
    return Date.now() >= payload.exp * 1000
  } catch {
    return true
  }
}

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()

  // ── Proactive token expiry check ──
  if (to.meta.requiresAuth && auth.isAuthenticated && isTokenExpired()) {
    auth.logout()
    const toast = useToastStore()
    toast.error('Sesi telah berakhir, silakan login kembali')
    return next({ name: 'Login', query: { redirect: to.fullPath } })
  }

  // Auth check
  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return next({ name: 'Login', query: { redirect: to.fullPath } })
  }
  if (to.meta.guest && auth.isAuthenticated) {
    return next({ name: 'Dashboard' })
  }

  // Role check
  const allowedRoles = to.meta.roles
  if (allowedRoles && allowedRoles.length > 0) {
    if (!auth.hasRole(...allowedRoles)) {
      return next({ name: 'Dashboard' })
    }
  }

  next()
})

export default router
