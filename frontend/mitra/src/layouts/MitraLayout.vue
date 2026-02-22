<template>
  <div class="layout-wrapper" :class="{ active: sidebarExpanded }">
    <div class="layout-flex">
      <AppSidebar
        :expanded="sidebarExpanded"
        @toggle="sidebarExpanded = !sidebarExpanded"
        @logout="handleLogout"
      />

      <div class="body-wrapper">
        <AppHeader
          :title="pageTitle"
          :subtitle="pageSubtitle"
          :sidebar-expanded="sidebarExpanded"
          :user-name="auth.userName"
          :user-initial="auth.userInitial"
          user-role="Mitra"
          @toggle-sidebar="sidebarExpanded = !sidebarExpanded"
          @logout="handleLogout"
        />

        <main class="main-content-area">
          <router-view />
        </main>
      </div>
    </div>

    <button @click="sidebarExpanded = !sidebarExpanded" class="mobile-toggle-btn">
      <i class="ri-menu-line" style="font-size:24px"></i>
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

const SESSION_CHECK_INTERVAL = 5 * 60 * 1000
let sessionCheckTimer = null

async function checkSession() {
  if (!auth.isAuthenticated) return
  try {
    await authApi.profile()
  } catch { /* 401 handled by interceptor */ }
}

function handleVisibilityChange() {
  if (document.visibilityState === 'visible') checkSession()
}

onMounted(() => {
  sessionCheckTimer = setInterval(checkSession, SESSION_CHECK_INTERVAL)
  document.addEventListener('visibilitychange', handleVisibilityChange)
  checkSession()
})

onUnmounted(() => {
  if (sessionCheckTimer) clearInterval(sessionCheckTimer)
  document.removeEventListener('visibilitychange', handleVisibilityChange)
})

const pageTitles = {
  Dashboard: 'Dashboard',
  Outlets: 'Jelajahi Outlet',
  OutletDetail: 'Detail Outlet',
  ApplyForm: 'Form Pengajuan',
  Applications: 'Pengajuan Kemitraan',
  ApplicationDetail: 'Detail Pengajuan',
  Invoices: 'Invoice',
  InvoiceDetail: 'Detail Invoice',
  Agreements: 'Agreement',
  Locations: 'Pengajuan Lokasi',
  Settings: 'Pengaturan Profil',
}
const pageSubtitles = {
  Dashboard: 'Pantau perkembangan kemitraan Anda',
  Outlets: 'Temukan peluang kemitraan yang sesuai',
  OutletDetail: 'Informasi lengkap & paket kemitraan',
  ApplyForm: 'Lengkapi data pengajuan kemitraan',
  Applications: 'Status pengajuan kemitraan Anda',
  ApplicationDetail: 'Informasi lengkap pengajuan kemitraan',
  Invoices: 'Riwayat tagihan & pembayaran',
  InvoiceDetail: 'Informasi lengkap tagihan Anda',
  Agreements: 'Dokumen perjanjian kemitraan',
  Locations: 'Pengajuan & status lokasi outlet',
  Settings: 'Kelola informasi akun & keamanan',
}

const pageTitle = computed(() => pageTitles[route.name] || '')
const pageSubtitle = computed(() => pageSubtitles[route.name] || '')

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>
