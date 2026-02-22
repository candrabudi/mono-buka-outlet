<template>
  <!-- Desktop Full Sidebar -->
  <aside class="sidebar-full" :class="{ active: expanded }">
    <div class="sidebar-full-header">
      <router-link to="/" class="sidebar-logo-link">
        <div class="sidebar-logo-icon">B</div>
        <div>
          <div class="sidebar-logo-text">BukaOutlet</div>
          <div class="sidebar-logo-sub">Portal Mitra</div>
        </div>
      </router-link>
      <button class="sidebar-toggle-btn" @click="$emit('toggle')">
        <i class="ri-arrow-left-s-line"></i>
      </button>
    </div>

    <div class="sidebar-body">
      <nav class="sidebar-nav-wrapper">
        <!-- Menu -->
        <div class="sidebar-section">
          <h4 class="sidebar-section-label">Menu</h4>
          <ul>
            <li>
              <router-link to="/" class="sidebar-link" :class="{ active: $route.path === '/' }">
                <span class="sidebar-link-icon"><i class="ri-home-5-line"></i></span>
                <span class="sidebar-link-text">Dashboard</span>
              </router-link>
            </li>
          </ul>
        </div>

        <!-- Jelajahi -->
        <div class="sidebar-section">
          <h4 class="sidebar-section-label">Jelajahi</h4>
          <ul>
            <li>
              <router-link to="/outlets" class="sidebar-link" :class="{ active: $route.path.startsWith('/outlets') }">
                <span class="sidebar-link-icon"><i class="ri-store-2-line"></i></span>
                <span class="sidebar-link-text">Outlet</span>
              </router-link>
            </li>
            <li>
              <router-link to="/applications" class="sidebar-link" :class="{ active: $route.path.startsWith('/applications') }">
                <span class="sidebar-link-icon"><i class="ri-file-list-3-line"></i></span>
                <span class="sidebar-link-text">Pengajuan</span>
              </router-link>
            </li>
          </ul>
        </div>

        <!-- Kemitraan -->
        <div class="sidebar-section">
          <h4 class="sidebar-section-label">Kemitraan</h4>
          <ul>
            <li>
              <router-link to="/invoices" class="sidebar-link" :class="{ active: $route.path.startsWith('/invoices') }">
                <span class="sidebar-link-icon"><i class="ri-receipt-line"></i></span>
                <span class="sidebar-link-text">Invoice</span>
              </router-link>
            </li>
            <li>
              <router-link to="/agreements" class="sidebar-link" :class="{ active: $route.path.startsWith('/agreements') }">
                <span class="sidebar-link-icon"><i class="ri-file-shield-line"></i></span>
                <span class="sidebar-link-text">Agreement</span>
              </router-link>
            </li>
            <li>
              <router-link to="/locations" class="sidebar-link" :class="{ active: $route.path.startsWith('/locations') }">
                <span class="sidebar-link-icon"><i class="ri-map-pin-line"></i></span>
                <span class="sidebar-link-text">Lokasi</span>
              </router-link>
            </li>
          </ul>
        </div>

        <!-- Akun -->
        <div class="sidebar-section">
          <h4 class="sidebar-section-label">Akun</h4>
          <ul>
            <li>
              <router-link to="/settings" class="sidebar-link" :class="{ active: $route.path.startsWith('/settings') }">
                <span class="sidebar-link-icon"><i class="ri-settings-3-line"></i></span>
                <span class="sidebar-link-text">Pengaturan</span>
              </router-link>
            </li>
          </ul>
        </div>
      </nav>
    </div>

    <!-- User area -->
    <div class="sidebar-user-area">
      <div class="sidebar-user-card">
        <div class="sidebar-user-avatar">{{ auth.userInitial }}</div>
        <div class="sidebar-user-info">
          <div class="sidebar-user-name">{{ auth.userName }}</div>
          <div class="sidebar-user-role">Mitra</div>
        </div>
        <button class="sidebar-logout-btn" @click="$emit('logout')" title="Logout">
          <i class="ri-logout-box-r-line"></i>
        </button>
      </div>
    </div>
  </aside>

  <!-- Overlay -->
  <div v-if="expanded" class="sidebar-overlay" @click="$emit('toggle')"></div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'
const auth = useAuthStore()

defineProps({ expanded: { type: Boolean, default: false } })
defineEmits(['toggle', 'logout'])
</script>
