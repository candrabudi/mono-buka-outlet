<template>
  <!-- Desktop Full Sidebar -->
  <aside class="sidebar-full" :class="{ active: expanded }">
    <div class="sidebar-full-header">
      <router-link to="/" class="sidebar-logo-link">
        <div class="sidebar-logo-icon">B</div>
        <div>
          <div class="sidebar-logo-text">BukaOutlet</div>
          <div class="sidebar-logo-sub">Admin Panel</div>
        </div>
      </router-link>
      <button type="button" class="sidebar-toggle-btn" @click="$emit('toggle')" title="Toggle sidebar">
        <svg width="16" height="40" viewBox="0 0 16 40" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M0 10C0 4.47715 4.47715 0 10 0H16V40H10C4.47715 40 0 35.5228 0 30V10Z" fill="var(--primary)"/>
          <path d="M10 15L6 20.0049L10 25.0098" stroke="#fff" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </button>
    </div>

    <div class="sidebar-body">
      <nav class="sidebar-nav-wrapper">
        <!-- Menu Section -->
        <div class="sidebar-section">
          <h4 class="sidebar-section-label">Menu</h4>
          <ul>
            <li>
              <router-link to="/" class="sidebar-link" :class="{ active: $route.name === 'Dashboard' }">
                <span class="sidebar-link-icon">
                  <svg width="18" height="21" viewBox="0 0 18 21" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path class="path-1" d="M0 8.84719C0 7.99027 0.366443 7.17426 1.00691 6.60496L6.34255 1.86217C7.85809 0.515019 10.1419 0.515019 11.6575 1.86217L16.9931 6.60496C17.6336 7.17426 18 7.99027 18 8.84719V17C18 19.2091 16.2091 21 14 21H4C1.79086 21 0 19.2091 0 17V8.84719Z" fill="currentColor"/>
                    <path class="path-2" d="M5 17C5 14.7909 6.79086 13 9 13C11.2091 13 13 14.7909 13 17V21H5V17Z" fill="var(--primary)"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Dashboard</span>
              </router-link>
            </li>
          </ul>
        </div>

        <!-- Kemitraan -->
        <div class="sidebar-section" v-if="auth.hasRole('master', 'admin')">
          <h4 class="sidebar-section-label">Kemitraan</h4>
          <ul>

            <li>
              <router-link to="/outlets" class="sidebar-link" :class="{ active: $route.path.startsWith('/outlets') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M1 3C1 1.89543 1.89543 1 3 1H17C18.1046 1 19 1.89543 19 3V5L10 9L1 5V3Z" fill="var(--primary)" class="path-2"/>
                    <path d="M1 7L10 11L19 7V17C19 18.1046 18.1046 19 17 19H3C1.89543 19 1 18.1046 1 17V7Z" fill="currentColor" class="path-1"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Outlet</span>
              </router-link>
            </li>
            <li>
              <router-link to="/outlet-categories" class="sidebar-link" :class="{ active: $route.path.startsWith('/outlet-categories') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <rect x="1" y="1" width="7" height="7" rx="2" fill="var(--primary)" class="path-2"/>
                    <rect x="12" y="1" width="7" height="7" rx="2" fill="currentColor" class="path-1"/>
                    <rect x="1" y="12" width="7" height="7" rx="2" fill="currentColor" class="path-1"/>
                    <rect x="12" y="12" width="7" height="7" rx="2" fill="currentColor" class="path-1"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Kategori</span>
              </router-link>
            </li>
            <li>
              <router-link to="/applications" class="sidebar-link" :class="{ active: $route.path.startsWith('/applications') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M9 12h6M9 16h6M17 21H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" stroke="currentColor" stroke-width="2" fill="none" class="path-1"/>
                    <path d="M13 3v5a1 1 0 001 1h5" stroke="var(--primary)" stroke-width="2" fill="none" class="path-2"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Pengajuan</span>
              </router-link>
            </li>
          </ul>
        </div>

        <!-- Partnership -->
        <div class="sidebar-section" v-if="auth.hasRole('master', 'admin', 'finance')">
          <h4 class="sidebar-section-label">Partnership</h4>
          <ul>
            <li>
              <router-link to="/partnerships" class="sidebar-link" :class="{ active: $route.path.startsWith('/partnerships') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 20 20" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M18 11C18 15.9706 13.9706 20 9 20C4.02944 20 0 15.9706 0 11C0 6.02944 4.02944 2 9 2C13.9706 2 18 6.02944 18 11Z" fill="currentColor" class="path-1"/>
                    <path d="M19.8025 8.01277C19.0104 4.08419 15.9158 0.989557 11.9872 0.197453C10.9045 -0.0208635 10 0.89543 10 2V8C10 9.10457 10.8954 10 12 10H18C19.1046 10 20.0209 9.09555 19.8025 8.01277Z" fill="var(--primary)" class="path-2"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Partnerships</span>
              </router-link>
            </li>
            <li v-if="auth.hasRole('master', 'admin')">
              <router-link to="/meetings" class="sidebar-link" :class="{ active: $route.path.startsWith('/meetings') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <rect x="3" y="4" width="18" height="18" rx="2" fill="currentColor" class="path-1"/>
                    <line x1="16" y1="2" x2="16" y2="6" stroke="var(--primary)" stroke-width="2" stroke-linecap="round" class="path-2"/>
                    <line x1="8" y1="2" x2="8" y2="6" stroke="var(--primary)" stroke-width="2" stroke-linecap="round" class="path-2"/>
                    <line x1="3" y1="10" x2="21" y2="10" stroke="var(--primary)" stroke-width="1.5" class="path-2"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Meetings</span>
              </router-link>
            </li>
            <li>
              <router-link to="/locations" class="sidebar-link" :class="{ active: $route.path.startsWith('/locations') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" fill="currentColor" class="path-1"/>
                    <circle cx="12" cy="10" r="3" stroke="var(--primary)" stroke-width="2" fill="none" class="path-2"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Lokasi</span>
              </router-link>
            </li>
          </ul>
        </div>

        <!-- Keuangan -->
        <div class="sidebar-section" v-if="auth.hasRole('master', 'finance')">
          <h4 class="sidebar-section-label">Keuangan</h4>
          <ul>
            <li>
              <router-link to="/payments" class="sidebar-link" :class="{ active: $route.path.startsWith('/payments') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="18" viewBox="0 0 20 18" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M20 4C20 1.79086 18.2091 0 16 0H4C1.79086 0 0 1.79086 0 4V14C0 16.2091 1.79086 18 4 18H16C18.2091 18 20 16.2091 20 14V4Z" fill="currentColor" class="path-1"/>
                    <path d="M6 9C6 7.34315 4.65685 6 3 6H0V12H3C4.65685 12 6 10.6569 6 9Z" fill="var(--primary)" class="path-2"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Pembayaran</span>
              </router-link>
            </li>
          </ul>
        </div>

        <!-- Management -->
        <div class="sidebar-section" v-if="auth.hasRole('master', 'admin')">
          <h4 class="sidebar-section-label">Management</h4>
          <ul>
            <li>
              <router-link to="/mitra" class="sidebar-link" :class="{ active: $route.path.startsWith('/mitra') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" stroke="currentColor" stroke-width="2" fill="none" class="path-1"/>
                    <circle cx="9" cy="7" r="4" fill="var(--primary)" class="path-2"/>
                    <path d="M23 21v-2a4 4 0 0 0-3-3.87" stroke="currentColor" stroke-width="2" fill="none" class="path-1"/>
                    <circle cx="16" cy="3.13" r="1" fill="var(--primary)" class="path-2"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Mitra</span>
              </router-link>
            </li>
            <li>
              <router-link to="/leader" class="sidebar-link" :class="{ active: $route.path.startsWith('/leader') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z" fill="currentColor" class="path-1"/>
                    <circle cx="12" cy="12" r="4" fill="var(--primary)" class="path-2"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Leader</span>
              </router-link>
            </li>
            <li v-if="auth.hasRole('master')">
              <router-link to="/users" class="sidebar-link" :class="{ active: $route.path.startsWith('/users') }">
                <span class="sidebar-link-icon">
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <ellipse cx="11.7778" cy="17.5555" rx="7.77778" ry="4.44444" fill="currentColor" class="path-1"/>
                    <circle cx="11.7778" cy="6.44444" r="4.44444" fill="var(--primary)" class="path-2"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Users</span>
              </router-link>
            </li>
            <li>
              <router-link to="/settings" class="sidebar-link" :class="{ active: $route.path.startsWith('/settings') }">
                <span class="sidebar-link-icon">
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <circle cx="12" cy="12" r="3" fill="var(--primary)" class="path-2"/>
                    <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z" stroke="currentColor" stroke-width="1.5" fill="none" class="path-1"/>
                  </svg>
                </span>
                <span class="sidebar-link-text">Pengaturan</span>
              </router-link>
            </li>
          </ul>
        </div>
      </nav>

      <!-- User Footer -->
      <div class="sidebar-user-area">
        <div class="sidebar-user-card">
          <div class="sidebar-user-avatar">{{ auth.userInitial }}</div>
          <div class="sidebar-user-info">
            <div class="sidebar-user-name">{{ auth.userName }}</div>
            <div class="sidebar-user-role">{{ auth.roleLabel(auth.userRole) }}</div>
          </div>
          <button @click="$emit('logout')" class="sidebar-logout-btn" title="Logout">
            <svg width="21" height="18" viewBox="0 0 21 18" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path fill-rule="evenodd" clip-rule="evenodd" d="M17.1464 10.4394C16.8536 10.7323 16.8536 11.2072 17.1464 11.5001C17.4393 11.7929 17.9142 11.7929 18.2071 11.5001L19.5 10.2072C20.1834 9.52375 20.1834 8.41571 19.5 7.73229L18.2071 6.4394C17.9142 6.1465 17.4393 6.1465 17.1464 6.4394C16.8536 6.73229 16.8536 7.20716 17.1464 7.50006L17.8661 8.21973H11.75C11.3358 8.21973 11 8.55551 11 8.96973C11 9.38394 11.3358 9.71973 11.75 9.71973H17.8661L17.1464 10.4394Z" fill="var(--primary)" class="path-2"/>
              <path fill-rule="evenodd" clip-rule="evenodd" d="M4.75 17.75H12C14.6234 17.75 16.75 15.6234 16.75 13C16.75 12.5858 16.4142 12.25 16 12.25C15.5858 12.25 15.25 12.5858 15.25 13C15.25 14.7949 13.7949 16.25 12 16.25H8.21412C7.34758 17.1733 6.11614 17.75 4.75 17.75ZM8.21412 1.75H12C13.7949 1.75 15.25 3.20507 15.25 5C15.25 5.41421 15.5858 5.75 16 5.75C16.4142 5.75 16.75 5.41421 16.75 5C16.75 2.37665 14.6234 0.25 12 0.25H4.75C6.11614 0.25 7.34758 0.82673 8.21412 1.75Z" fill="currentColor" class="path-1"/>
              <path fill-rule="evenodd" clip-rule="evenodd" d="M0 5C0 2.37665 2.12665 0.25 4.75 0.25C7.37335 0.25 9.5 2.37665 9.5 5V13C9.5 15.6234 7.37335 17.75 4.75 17.75C2.12665 17.75 0 15.6234 0 13V5Z" fill="currentColor" class="path-1"/>
            </svg>
          </button>
        </div>
      </div>
    </div>
  </aside>

  <!-- Mobile overlay -->
  <div v-if="expanded" class="sidebar-overlay" @click="$emit('toggle')"></div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()

defineProps({
  expanded: { type: Boolean, default: true }
})

defineEmits(['toggle', 'logout'])
</script>
