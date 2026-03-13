import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'

const APP_NAME = 'Outlet Ready'

const routes = [
  // Auth
  { path: '/login', name: 'Login', component: () => import('../pages/Login.vue'), meta: { guest: true, title: 'Login' } },

  // Admin Dashboard
  {
    path: '/',
    component: () => import('../layouts/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      // Dashboard — all admin roles
      { path: '', name: 'Dashboard', component: () => import('../pages/Dashboard.vue'), meta: { title: 'Dashboard' } },

      // Outlets — master, admin
      { path: 'outlets', name: 'Outlets', component: () => import('../pages/outlets/OutletList.vue'), meta: { roles: ['master', 'admin'], title: 'Kelola Outlet' } },
      { path: 'outlets/create', name: 'OutletCreate', component: () => import('../pages/outlets/OutletForm.vue'), meta: { roles: ['master', 'admin'], title: 'Tambah Outlet' } },
      { path: 'outlets/:id', name: 'OutletEdit', component: () => import('../pages/outlets/OutletForm.vue'), meta: { roles: ['master', 'admin'], title: 'Edit Outlet' } },
      { path: 'outlet-categories', name: 'OutletCategories', component: () => import('../pages/outlets/OutletCategoryList.vue'), meta: { roles: ['master', 'admin'], title: 'Kategori Outlet' } },

      // Partnerships — master, admin, finance (view)
      { path: 'partnerships', name: 'Partnerships', component: () => import('../pages/partnerships/PartnershipList.vue'), meta: { roles: ['master', 'admin', 'finance'], title: 'Partnership' } },
      { path: 'partnerships/:id', name: 'PartnershipDetail', component: () => import('../pages/partnerships/PartnershipDetail.vue'), meta: { roles: ['master', 'admin', 'finance'], title: 'Detail Partnership' } },

      // Applications — master, admin
      { path: 'applications', name: 'Applications', component: () => import('../pages/applications/ApplicationList.vue'), meta: { roles: ['master', 'admin'], title: 'Aplikasi Masuk' } },

      // Payments — master, finance
      { path: 'payments', name: 'Payments', component: () => import('../pages/payments/PaymentList.vue'), meta: { roles: ['master', 'finance'], title: 'Pembayaran' } },

      // Meetings — master, admin
      { path: 'meetings', name: 'Meetings', component: () => import('../pages/meetings/MeetingList.vue'), meta: { roles: ['master', 'admin'], title: 'Meeting' } },

      // Mitra — master, admin
      { path: 'mitra', name: 'Mitra', component: () => import('../pages/mitra/MitraList.vue'), meta: { roles: ['master', 'admin'], title: 'Kelola Mitra' } },

      // Leader — master, admin
      { path: 'leader', name: 'Leader', component: () => import('../pages/leader/LeaderList.vue'), meta: { roles: ['master', 'admin'], title: 'Kelola Affiliator' } },

      // Users — master only
      { path: 'users', name: 'Users', component: () => import('../pages/users/UserList.vue'), meta: { roles: ['master'], title: 'Kelola User' } },

      // Settings — master, admin
      { path: 'settings', name: 'Settings', component: () => import('../pages/settings/SettingList.vue'), meta: { roles: ['master', 'admin'], title: 'Pengaturan' } },

      // Location Submissions — master, admin
      { path: 'locations', name: 'Locations', component: () => import('../pages/locations/LocationList.vue'), meta: { roles: ['master', 'admin'], title: 'Lokasi' } },
      { path: 'locations/:id', name: 'LocationDetail', component: () => import('../pages/locations/LocationDetail.vue'), meta: { roles: ['master', 'admin'], title: 'Detail Lokasi' } },

      // Ebooks — master, admin
      { path: 'ebooks', name: 'Ebooks', component: () => import('../pages/ebooks/EbookList.vue'), meta: { roles: ['master', 'admin'], title: 'E-Book' } },
      { path: 'ebooks/create', name: 'EbookCreate', component: () => import('../pages/ebooks/EbookForm.vue'), meta: { roles: ['master', 'admin'], title: 'Tambah E-Book' } },
      { path: 'ebooks/:id', name: 'EbookEdit', component: () => import('../pages/ebooks/EbookForm.vue'), meta: { roles: ['master', 'admin'], title: 'Edit E-Book' } },
      { path: 'ebook-orders', name: 'EbookOrders', component: () => import('../pages/ebooks/EbookOrderList.vue'), meta: { roles: ['master', 'admin'], title: 'Pesanan E-Book' } },
      { path: 'ebook-categories', name: 'EbookCategories', component: () => import('../pages/ebooks/EbookCategories.vue'), meta: { roles: ['master', 'admin'], title: 'Kategori E-Book' } },
      { path: 'download-requests', name: 'DownloadRequests', component: () => import('../pages/ebooks/DownloadRequests.vue'), meta: { roles: ['master', 'admin'], title: 'Permintaan Download' } },

      // AI Konsultan — master, admin
      { path: 'ai-settings', name: 'AISettings', component: () => import('../pages/ai/AISettings.vue'), meta: { roles: ['master', 'admin'], title: 'AI Konsultan' } },

      // Affiliator Management — master, admin
      { path: 'affiliator-management', name: 'AffiliatorManagement', component: () => import('../pages/affiliator/AffiliatorManagement.vue'), meta: { roles: ['master', 'admin'], title: 'Manajemen Affiliator' } },

      // Affiliator — affiliator self-service
      { path: 'affiliator-dashboard', name: 'AffiliatorDashboard', component: () => import('../pages/affiliator/AffiliatorDashboard.vue'), meta: { roles: ['affiliator'], title: 'Dashboard Affiliator' } },
    ],
  },

  { path: '/:pathMatch(.*)*', name: 'NotFound', component: () => import('../pages/NotFound.vue'), meta: { title: 'Halaman Tidak Ditemukan' } },
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

  // When authenticated guest tries to access login, redirect to appropriate dashboard
  if (to.meta.guest && auth.isAuthenticated) {
    if (auth.userRole === 'affiliator') {
      return next({ name: 'AffiliatorDashboard' })
    }
    return next({ name: 'Dashboard' })
  }

  // Affiliator visiting root Dashboard should go to their own dashboard
  if (to.name === 'Dashboard' && auth.isAuthenticated && auth.userRole === 'affiliator') {
    return next({ name: 'AffiliatorDashboard' })
  }

  // Role check
  const allowedRoles = to.meta.roles
  if (allowedRoles && allowedRoles.length > 0) {
    if (!auth.hasRole(...allowedRoles)) {
      if (auth.userRole === 'affiliator') {
        return next({ name: 'AffiliatorDashboard' })
      }
      return next({ name: 'Dashboard' })
    }
  }

  next()
})

// ── Dynamic page title ──
router.afterEach((to) => {
  const pageTitle = to.meta.title
  document.title = pageTitle ? `${pageTitle} — ${APP_NAME}` : APP_NAME
})

export default router
