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
        <svg width="16" height="40" viewBox="0 0 16 40" fill="none">
          <path d="M0 10C0 4.47715 4.47715 0 10 0H16V40H10C4.47715 40 0 35.5228 0 30V10Z" fill="var(--primary)"/>
          <path d="M10 15L6 20.0049L10 25.0098" stroke="#fff" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
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
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z" stroke="currentColor" stroke-width="2" class="path-1"/>
                    <polyline points="9 22 9 12 15 12 15 22" stroke="currentColor" stroke-width="2" class="path-2"/>
                  </svg>
                </span>
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
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path d="M3 3h18v18H3V3z" stroke="currentColor" stroke-width="2" class="path-1"/>
                    <path d="M3 9h18M9 21V9" stroke="currentColor" stroke-width="2" class="path-2"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Outlet</span>
              </router-link>
            </li>
            <li>
              <router-link to="/applications" class="sidebar-link" :class="{ active: $route.path.startsWith('/applications') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path d="M9 5H7a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V7a2 2 0 0 0-2-2h-2" stroke="currentColor" stroke-width="2" class="path-1"/>
                    <rect x="9" y="3" width="6" height="4" rx="1" stroke="currentColor" stroke-width="2" class="path-2"/>
                  </svg>
                </span>
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
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="2" class="path-1"/>
                    <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2" class="path-2"/>
                    <line x1="16" y1="13" x2="8" y2="13" stroke="currentColor" stroke-width="2" class="path-2"/>
                    <line x1="16" y1="17" x2="8" y2="17" stroke="currentColor" stroke-width="2" class="path-2"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Invoice</span>
              </router-link>
            </li>
            <li>
              <router-link to="/agreements" class="sidebar-link" :class="{ active: $route.path.startsWith('/agreements') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" stroke-width="2" class="path-1"/>
                    <polyline points="14 2 14 8 20 8" stroke="currentColor" stroke-width="2" class="path-2"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Agreement</span>
              </router-link>
            </li>
            <li>
              <router-link to="/locations" class="sidebar-link" :class="{ active: $route.path.startsWith('/locations') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                    <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" stroke="currentColor" stroke-width="2" class="path-1"/>
                    <circle cx="12" cy="10" r="3" stroke="currentColor" stroke-width="2" class="path-2"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Lokasi</span>
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
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
            <path d="M15 10L13.7071 11.2929C13.3166 11.6834 13.3166 12.3166 13.7071 12.7071L15 14M14 12L22 12M6 20C3.79086 20 2 18.2091 2 16V8C2 5.79086 3.79086 4 6 4M6 20C8.20914 20 10 18.2091 10 16V8C10 5.79086 8.20914 4 6 4M6 20H14C16.2091 20 18 18.2091 18 16M6 4H14C16.2091 4 18 5.79086 18 8" class="path-1"/>
          </svg>
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
