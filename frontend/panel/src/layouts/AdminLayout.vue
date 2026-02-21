<template>
  <div class="layout-wrapper" :class="{ active: sidebarExpanded }">
    <div class="layout-flex">
      <!-- Sidebar Component -->
      <AppSidebar
        :expanded="sidebarExpanded"
        @toggle="sidebarExpanded = !sidebarExpanded"
        @logout="handleLogout"
      />

      <!-- Main Body -->
      <div class="body-wrapper">
        <!-- Header Component -->
        <AppHeader
          :title="pageTitle"
          :subtitle="pageSubtitle"
          :sidebar-expanded="sidebarExpanded"
          :user-name="auth.userName"
          :user-initial="auth.userInitial"
          :user-role="auth.roleLabel(auth.userRole)"
          @toggle-sidebar="sidebarExpanded = !sidebarExpanded"
          @logout="handleLogout"
        />

        <!-- Page Content -->
        <main class="main-content-area">
          <router-view />
        </main>
      </div>
    </div>

    <!-- Mobile toggle -->
    <button
      @click="sidebarExpanded = !sidebarExpanded"
      class="mobile-toggle-btn"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
      </svg>
    </button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import { authApi } from '../services/api'
import AppSidebar from '../components/AppSidebar.vue'
import AppHeader from '../components/AppHeader.vue'

const auth = useAuthStore()
const toast = useToastStore()
const router = useRouter()
const route = useRoute()
const sidebarExpanded = ref(true)

// ── Session Health Check ──
const SESSION_CHECK_INTERVAL = 5 * 60 * 1000 // 5 minutes
let sessionCheckTimer = null

async function checkSession() {
  if (!auth.isAuthenticated) return
  try {
    await authApi.profile()
  } catch (err) {
    // 401 is handled by the axios interceptor, which already redirects to login
    // This catch is just to prevent unhandled promise rejection
  }
}

function handleVisibilityChange() {
  if (document.visibilityState === 'visible') {
    checkSession()
  }
}

onMounted(() => {
  // Start polling session health
  sessionCheckTimer = setInterval(checkSession, SESSION_CHECK_INTERVAL)

  // Check session when user returns to this tab
  document.addEventListener('visibilitychange', handleVisibilityChange)

  // Check immediately on mount
  checkSession()
})

onUnmounted(() => {
  if (sessionCheckTimer) clearInterval(sessionCheckTimer)
  document.removeEventListener('visibilitychange', handleVisibilityChange)
})

const pageTitles = {
  Dashboard: 'Dashboard',

  Outlets: 'Outlet Management',
  OutletCreate: 'Tambah Outlet',
  OutletEdit: 'Edit Outlet',
  OutletCategories: 'Kategori Outlet',
  Partnerships: 'Partnership Management',
  PartnershipDetail: 'Detail Partnership',
  Payments: 'Pembayaran',
  Meetings: 'Meeting Management',
  Mitra: 'Management Mitra',
  Leader: 'Management Leader',
  Users: 'User Management',
  Settings: 'Pengaturan Sistem',
  Locations: 'Pengajuan Lokasi',
  LocationDetail: 'Detail Lokasi',
}
const pageSubtitles = {
  Dashboard: "Let's check your update today",

  Outlets: 'Kelola semua data outlet',
  OutletCreate: 'Isi informasi outlet baru',
  OutletEdit: 'Perbarui informasi outlet',
  OutletCategories: 'Kelola kategori pengelompokan outlet',
  Partnerships: 'Kelola semua kemitraan aktif',
  Payments: 'Verifikasi & kelola pembayaran',
  Meetings: 'Kelola agenda, notulensi, dan action plan',
  Mitra: 'Kelola data mitra dan akses portal',
  Leader: 'Kelola data leader yang menangani mitra',
  Users: 'Kelola pengguna sistem',
  Settings: 'Konfigurasi aplikasi & integrasi pembayaran',
  Locations: 'Kelola pengajuan & survei lokasi mitra',
  LocationDetail: 'Detail informasi, scoring, survei & approval',
}

const pageTitle = computed(() => pageTitles[route.name] || '')
const pageSubtitle = computed(() => pageSubtitles[route.name] || '')

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>
